package app_db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/g45t345rt/g45w/app_db/order_column"
	"github.com/g45t345rt/g45w/app_db/schema_version"
)

type IPFSGateway struct {
	ID           int64
	Name         string
	Endpoint     string
	FetchHeaders map[string]string
	Active       bool
	OrderNumber  int
}

var gatewayOrderer = order_column.Orderer{
	TableName:  "ipfs_gateways",
	ColumnName: "order_number",
}

var TRUSTED_IPFS_GATEWAYS = []IPFSGateway{
	{Name: "deronfts.com", Endpoint: "https://ipfs.deronfts.com/ipfs/{cid}"},
	{Name: "ipfs.io", Endpoint: "https://ipfs.io/ipfs/{cid}"},
	{Name: "pinata.cloud", Endpoint: "https://gateway.pinata.cloud/ipfs/{cid}"},
	{Name: "dweb.link", Endpoint: "https://dweb.link/ipfs/{cid}"},
	{Name: "nftstorage.link", Endpoint: "https://{cid}.ipfs.nftstorage.link"},
}

func initTableIPFSGateways() error {
	version, err := schema_version.GetVersion(DB, "ipfs_gateways")
	if err != nil {
		return err
	}

	if version == 0 {
		_, err := DB.Exec(`
			CREATE TABLE IF NOT EXISTS ipfs_gateways (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				endpoint VARCHAR,
				name VARCHAR,
				fetch_headers VARCHAR,
				active BOOL
			);
		`)
		if err != nil {
			return err
		}

		_, err = DB.Exec(`
			ALTER TABLE ipfs_gateways ADD COLUMN order_number INT NOT NULL DEFAULT 0;
		`)
		if err != nil {
			return err
		}

		err = ReOrderIPFSGateways(DB)
		if err != nil {
			return err
		}

		version = 1
		err = schema_version.StoreVersion(DB, "ipfs_gateways", version)
		if err != nil {
			return err
		}
	}

	return err
}

func ResetIPFSGateways() error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM ipfs_gateways`)
	if err != nil {
		tx.Rollback()
		return err
	}

	for i, gateway := range TRUSTED_IPFS_GATEWAYS {
		_, err = tx.Exec(`
			INSERT INTO ipfs_gateways (endpoint, name, active, order_number)
			VALUES (?,?,?,?);
		`, gateway.Endpoint, gateway.Name, true, i)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

type GetIPFSGatewaysParams struct {
	Active sql.NullBool
}

func GetIPFSGateways(params GetIPFSGatewaysParams) ([]IPFSGateway, error) {
	query := sq.Select("*").From("ipfs_gateways").OrderBy("order_number ASC")

	if params.Active.Valid {
		query = query.Where(sq.Eq{"active": params.Active.Bool})
	}

	rows, err := query.RunWith(DB).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ipfsGateways []IPFSGateway
	for rows.Next() {
		var ipfsGateway IPFSGateway

		var fetchHeaders sql.NullString
		err = rows.Scan(
			&ipfsGateway.ID,
			&ipfsGateway.Endpoint,
			&ipfsGateway.Name,
			&fetchHeaders,
			&ipfsGateway.Active,
			&ipfsGateway.OrderNumber,
		)
		if err != nil {
			return nil, err
		}

		if fetchHeaders.Valid {
			err = json.Unmarshal([]byte(fetchHeaders.String), &ipfsGateway.FetchHeaders)
			if err != nil {
				return nil, err
			}
		}

		ipfsGateways = append(ipfsGateways, ipfsGateway)
	}

	return ipfsGateways, nil
}

func GetIPFSGateway(id int64) (gateway IPFSGateway, err error) {
	row := DB.QueryRow(`
	SELECT * FROM ipfs_gateways
	WHERE id = ?;
	`, id)
	if err != nil {
		return
	}

	err = row.Err()
	if err != nil {
		return
	}

	var fetchHeaders sql.NullString
	err = row.Scan(
		&gateway.ID,
		&gateway.Endpoint,
		&gateway.Name,
		&fetchHeaders,
		&gateway.Active,
		&gateway.OrderNumber,
	)
	if err != nil {
		return
	}

	if fetchHeaders.Valid {
		err = json.Unmarshal([]byte(fetchHeaders.String), &gateway.FetchHeaders)
		if err != nil {
			return
		}
	}

	return
}

func InsertIPFSGateway(gateway IPFSGateway) error {
	fetchHeadersBytes, err := json.Marshal(gateway.FetchHeaders)
	if err != nil {
		return err
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	gateway.OrderNumber, err = gatewayOrderer.GetNewOrderNumber(tx)
	if err != nil {
		return err
	}

	/*
		err = walletOrderer.Insert(tx, gateway.OrderNumber)
		if err != nil {
			tx.Rollback()
			return err
		}
	*/

	result, err := tx.Exec(`
		INSERT INTO ipfs_gateways (name,endpoint,fetch_headers,active,order_number)
		VALUES (?,?,?,?,?);
	`, gateway.Name, gateway.Endpoint, string(fetchHeadersBytes), gateway.Active, gateway.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	gateway.ID = id
	return tx.Commit()
}

func UpdateIPFSGateway(gateway IPFSGateway) error {
	fetchHeadersBytes, err := json.Marshal(gateway.FetchHeaders)
	if err != nil {
		return err
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	currentGateway, err := GetIPFSGateway(gateway.ID)
	if err != nil {
		return err
	}

	err = gatewayOrderer.Update(tx, currentGateway.OrderNumber, gateway.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE ipfs_gateways
		SET name = ?, endpoint = ?, fetch_headers = ?, active = ?, order_number = ?
		WHERE id = ?;
	`, gateway.Name, gateway.Endpoint, string(fetchHeadersBytes), gateway.Active, gateway.OrderNumber, gateway.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DelIPFSGateway(id int64) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	gateway, err := GetIPFSGateway(id)
	if err != nil {
		return err
	}

	err = gatewayOrderer.Delete(tx, gateway.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM ipfs_gateways
		WHERE id = ?;
	`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (i IPFSGateway) Fetch(cId string, timeout time.Duration) (*http.Response, error) {
	endpoint := strings.Replace(i.Endpoint, "{cid}", cId, -1)

	client := new(http.Client)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (i IPFSGateway) TestFetch() error {
	cId := "bafybeibozpulxtpv5nhfa2ue3dcjx23ndh3gwr5vwllk7ptoyfwnfjjr4q"
	res, err := i.Fetch(cId, 5*time.Second)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%d", res.StatusCode)
	}

	return nil
}

func ReOrderIPFSGateways(db *sql.DB) error {
	gateways, err := GetIPFSGateways(GetIPFSGatewaysParams{})
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for i, node := range gateways {
		_, err = tx.Exec(`
				UPDATE ipfs_gateways
				SET order_number = ?
				WHERE id = ?;
			`, i, node.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
