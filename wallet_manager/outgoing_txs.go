package wallet_manager

import (
	"database/sql"
	"encoding/hex"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/deroproject/derohe/block"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/transaction"
	"github.com/deroproject/derohe/walletapi"
)

type OutgoingTx struct {
	TxId        string
	HeightBuilt sql.NullInt64
	Timestamp   sql.NullInt64
	Status      sql.NullString
	TxType      sql.NullInt32
	HexData     sql.NullString
	BlockHeight sql.NullInt64
	Description sql.NullString
}

func initDatabaseOutgoingTxs(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS outgoing_txs (
			tx_id VARCHAR PRIMARY KEY,
			height_built BIGINT,
			timestamp BIGINT,
			status VARCHAR,
			tx_type VARCHAR,
			hex_data VARCHAR,
			block_height BIGINT,
			description VARCHAR
		);
	`)
	return err
}

func rowsScanOutgoingTxs(rows *sql.Rows) ([]OutgoingTx, error) {
	defer rows.Close()

	var outgoingTxs []OutgoingTx
	for rows.Next() {
		var outgoingTx OutgoingTx
		err := rows.Scan(
			&outgoingTx.TxId,
			&outgoingTx.HeightBuilt,
			&outgoingTx.Timestamp,
			&outgoingTx.Status,
			&outgoingTx.TxType,
			&outgoingTx.HexData,
			&outgoingTx.BlockHeight,
			&outgoingTx.Description,
		)
		if err != nil {
			return nil, err
		}

		outgoingTxs = append(outgoingTxs, outgoingTx)
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return outgoingTxs, nil
}

type GetOutgoingTxsParams struct {
	Descending bool
	OrderBy    string
	Limit      *uint64
	TxType     *transaction.TransactionType
}

func (w *Wallet) GetOutgoingTxs(params GetOutgoingTxsParams) ([]OutgoingTx, error) {
	query := sq.Select("*").From("outgoing_txs")

	if params.TxType != nil {
		query = query.Where(sq.Eq{"tx_type": params.TxType})
	}

	if len(params.OrderBy) > 0 {
		direction := "ASC"
		if params.Descending {
			direction = "DESC"
		}

		query = query.OrderBy(params.OrderBy, direction)
	}

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	rows, err := w.DB.Query(query.ToSql())
	if err != nil {
		return nil, err
	}

	return rowsScanOutgoingTxs(rows)
}

var pendingTries map[string]int

func (w *Wallet) CheckRegistrationTx(tx transaction.Transaction) (rpc.GetEncryptedBalance_Result, bool, error) {
	// registration does not give a valid block even if successful
	// use GetEncryptedBalance to get registration height/block

	var balanceResult rpc.GetEncryptedBalance_Result
	err := walletapi.RPC_Client.Call("DERO.GetEncryptedBalance", rpc.GetEncryptedBalance_Params{
		TopoHeight: -1,
		Address:    w.Info.Addr,
	}, &balanceResult)
	if err != nil {
		return balanceResult, false, err
	}

	var blockResult rpc.GetBlock_Result
	err = walletapi.RPC_Client.Call("DERO.GetBlock", rpc.GetBlock_Params{
		Height: uint64(balanceResult.Registration),
	}, &blockResult)
	if err != nil {
		return balanceResult, false, err
	}

	var block block.Block
	data, _ := hex.DecodeString(blockResult.Blob)
	block.Deserialize(data)

	for _, blockTx := range block.Tx_hashes {
		if blockTx.String() == tx.GetHash().String() {
			return balanceResult, true, nil
		}
	}

	return balanceResult, false, nil
}

func (w *Wallet) UpdatePendingOutgoingTxs() (int, error) {
	if !walletapi.Connected {
		return 0, nil
	}

	if pendingTries == nil {
		pendingTries = make(map[string]int)
	}

	rows, err := w.DB.Query(`
		SELECT *
		FROM outgoing_txs
		WHERE status = 'pending';
	`)
	if err != nil {
		return 0, err
	}

	outgoingTxs, err := rowsScanOutgoingTxs(rows)
	if err != nil {
		return 0, err
	}

	var txIds []string
	for _, outgoingTx := range outgoingTxs {
		txIds = append(txIds, outgoingTx.TxId)
	}

	if len(txIds) == 0 {
		return 0, nil
	}

	var txResult rpc.GetTransaction_Result
	err = walletapi.RPC_Client.Call("DERO.GetTransaction", rpc.GetTransaction_Params{
		Tx_Hashes: txIds,
	}, &txResult)
	if err != nil {
		return 0, err
	}

	updated := 0

	for i, info := range txResult.Txs {
		txHex := txResult.Txs_as_hex[i]

		var tx transaction.Transaction
		data, _ := hex.DecodeString(txHex)
		err := tx.Deserialize(data)
		if err != nil {
			return updated, err
		}

		txId := tx.GetHash().String()
		valid := false
		var blockHeight int64

		if tx.TransactionType == transaction.REGISTRATION {
			balance, regValid, err := w.CheckRegistrationTx(tx)
			if err != nil {
				return updated, err
			}

			if regValid {
				valid = true
				blockHeight = balance.Registration
			}
		} else if info.ValidBlock != "" {
			valid = true
			blockHeight = info.Block_Height
		}

		if valid {
			err = w.UpdateOugoingTx(txId, "valid", blockHeight)
			if err != nil {
				return updated, err
			}

			updated += 1
		} else {
			// if after 30 tries the transaction is still not in a valid block we set invalid status
			tries, ok := pendingTries[txId]
			if !ok {
				pendingTries[txId] = 1
			}

			if tries >= 30 {
				err = w.UpdateOugoingTx(txId, "invalid", 0)
				if err != nil {
					return updated, err
				}

				updated += 1
				delete(pendingTries, txId)
			} else {
				pendingTries[txId] = tries + 1
			}
		}
	}

	return updated, nil
}

func (w *Wallet) UpdateOugoingTx(txId string, status string, blockHeight int64) error {
	_, err := w.DB.Exec(`
		UPDATE outgoing_txs
		SET status = ?, block_height = ?
		WHERE tx_id = ?;
	`, status, blockHeight, txId)
	return err
}

func (w *Wallet) StoreOutgoingTx(tx *transaction.Transaction, description string) error {
	txId := tx.GetHash().String()
	height := tx.Height
	txType := tx.TransactionType
	hexData := hex.EncodeToString(tx.Serialize())

	_, err := w.DB.Exec(`
		INSERT INTO outgoing_txs (tx_id,height_built,tx_type,timestamp,status,hex_data,description)
		VALUES (?,?,?,?,?,?,?)
		ON CONFLICT DO NOTHING;
	`, txId, height, txType, time.Now().Unix(), "pending", hexData, description)
	return err
}

func (w *Wallet) DelOutgoingTx(txId string) error {
	_, err := w.DB.Exec(`
		DELETE FROM outgoing_txs
		WHERE tx_id = ?;
	`, txId)
	return err
}
