package app_db

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	sq "github.com/Masterminds/squirrel"
	"github.com/deroproject/derohe/globals"
	"github.com/g45t345rt/g45w/app_db/order_column"
	"github.com/g45t345rt/g45w/app_db/schema_version"
	"github.com/g45t345rt/g45w/settings"
)

type WalletInfo struct {
	Addr              string
	Name              string
	RegistrationTxHex string
	Timestamp         int64
	OrderNumber       int
}

var walletOrderer = order_column.Orderer{
	TableName:  "wallets",
	ColumnName: "order_number",
}

func initTableWallets() error {
	version, err := schema_version.GetVersion(DB, "wallets")
	if err != nil {
		return err
	}

	if version == 0 {
		_, err := DB.Exec(`
			CREATE TABLE IF NOT EXISTS wallets (
				addr VARCHAR PRIMARY KEY,
				name VARCHAR NOT NULL,
				registration_tx_hex VARCHAR NOT NULL,
				timestamp BIGINT NOT NULL,
				order_number INT NOT NULL
			);
		`)
		if err != nil {
			return err
		}

		err = migrateJsonWalletsInfo()
		if err != nil {
			return err
		}

		version = 1
		err = schema_version.StoreVersion(DB, "wallets", version)
		if err != nil {
			return err
		}
	}

	return nil
}

func delWalletInfoIfNoFolder() error {
	wallets, err := GetWallets()
	if err != nil {
		return err
	}

	walletsDir := settings.WalletsDir
	for _, info := range wallets {
		walletPath := filepath.Join(walletsDir, info.Addr)
		_, err := os.Stat(walletPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = DelWalletInfo(info.Addr)
				if err != nil {
					return err
				}
			}

			return err
		}
	}

	return nil
}

type JsonWalletInfo struct {
	Name              string `json:"name"`
	Addr              string `json:"addr"`
	RegistrationTxHex string `json:"registration_tx_hex"`
	Timestamp         int64  `json:"timestamp"`
}

// Before using sql lite I was using json file to store wallet info
// use this func to migrate json wallet files in sql db for previous app versions
func migrateJsonWalletsInfo() error {
	walletsDir := settings.WalletsDir

	return filepath.Walk(walletsDir, func(path string, info fs.FileInfo, fileErr error) error {
		if walletsDir == path {
			return nil
		}

		if info != nil && info.IsDir() {
			addr := info.Name()

			_, err := globals.ParseValidateAddress(addr)
			if err != nil {
				return nil
			}

			walletInfoPath := filepath.Join(walletsDir, addr, "info.json")
			data, err := os.ReadFile(walletInfoPath)
			if err != nil {
				return nil
			}

			var walletInfo JsonWalletInfo
			err = json.Unmarshal(data, &walletInfo)
			if err != nil {
				return err
			}

			err = InsertWalletInfo(WalletInfo{
				Addr:              walletInfo.Addr,
				Name:              walletInfo.Name,
				RegistrationTxHex: walletInfo.RegistrationTxHex,
				Timestamp:         walletInfo.Timestamp,
				// OrderNumber will be automatically set to last
			})
			if err != nil {
				return err
			}

			err = os.Remove(walletInfoPath)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func GetWallets() ([]WalletInfo, error) {
	query := sq.Select("*").From("wallets").OrderBy("order_number ASC")

	var wallets []WalletInfo
	rows, err := query.RunWith(DB).Query()
	if err != nil {
		return nil, nil
	}

	for rows.Next() {
		var wallet WalletInfo
		err = rows.Scan(
			&wallet.Addr,
			&wallet.Name,
			&wallet.RegistrationTxHex,
			&wallet.Timestamp,
			&wallet.OrderNumber,
		)
		if err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func GetWalletInfo(addr string) (WalletInfo, error) {
	query := sq.Select("*").From("wallets").Where(sq.Eq{"addr": addr})

	var walletInfo WalletInfo
	row := query.RunWith(DB).QueryRow()
	err := row.Scan(
		&walletInfo.Addr,
		&walletInfo.Name,
		&walletInfo.RegistrationTxHex,
		&walletInfo.Timestamp,
		&walletInfo.OrderNumber,
	)
	return walletInfo, err
}

func InsertWalletInfo(walletInfo WalletInfo) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	walletInfo.OrderNumber, err = walletOrderer.GetNewOrderNumber(tx)
	if err != nil {
		return err
	}

	err = walletOrderer.Insert(tx, walletInfo.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO wallets (addr,name,registration_tx_hex,timestamp,order_number)
		VALUES (?,?,?,?,?);
	`, walletInfo.Addr, walletInfo.Name, walletInfo.RegistrationTxHex, walletInfo.Timestamp, walletInfo.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func UpdateWalletInfo(walletInfo WalletInfo) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	currentWalletInfo, err := GetWalletInfo(walletInfo.Addr)
	if err != nil {
		return err
	}

	err = walletOrderer.Update(tx, currentWalletInfo.OrderNumber, walletInfo.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE wallets
		SET name = ?,
				registration_tx_hex = ?,
				order_number = ?
		WHERE addr = ?;
	`, walletInfo.Name, walletInfo.RegistrationTxHex, walletInfo.OrderNumber, walletInfo.Addr)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DelWalletInfo(addr string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	walletInfo, err := GetWalletInfo(addr)
	if err != nil {
		return err
	}

	err = walletOrderer.Delete(tx, walletInfo.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM wallets
		WHERE addr = ?;
	`, addr)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
