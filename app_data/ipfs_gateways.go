package app_data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type IPFSGateway struct {
	ID           int64
	Name         string
	Endpoint     string
	FetchHeaders map[string]string
	Active       bool
}

var TRUSTED_IPFS_GATEWAYS = []IPFSGateway{
	{Name: "deronfts.com", Endpoint: "https://ipfs.deronfts.com/ipfs/{cid}"},
	{Name: "ipfs.io", Endpoint: "https://ipfs.io/ipfs/{cid}"},
	{Name: "pinata.cloud", Endpoint: "https://gateway.pinata.cloud/ipfs/{cid}"},
	{Name: "dweb.link", Endpoint: "https://dweb.link/ipfs/{cid}"},
	{Name: "nftstorage.link", Endpoint: "https://{cid}.ipfs.nftstorage.link"},
}

func initDatabaseIPFSGateways() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS ipfs_gateways (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			endpoint VARCHAR,
			name VARCHAR,
			fetch_headers VARCHAR,
			active BOOL
		);
	`)
	return err
}

func StoreTrustedIPFSGateways() error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	for _, gateway := range TRUSTED_IPFS_GATEWAYS {
		_, err = tx.Exec(`
			INSERT INTO ipfs_gateways (endpoint, name, active)
			VALUES (?,?,?);
		`, gateway.Endpoint, gateway.Name, true)
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
	query := sq.Select("*").From("ipfs_gateways")

	if params.Active.Valid {
		query = query.Where(sq.Eq{"active": params.Active.Bool})
	}

	rows, err := query.RunWith(DB).Query()
	if err != nil {
		return nil, err
	}

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

func InsertIPFSGateway(gateway IPFSGateway) error {
	fetchHeadersBytes, err := json.Marshal(gateway.FetchHeaders)
	if err != nil {
		return err
	}

	result, err := DB.Exec(`
		INSERT INTO ipfs_gateways (name,endpoint,fetch_headers,active)
		VALUES (?,?,?,?);
	`, gateway.Name, gateway.Endpoint, string(fetchHeadersBytes), gateway.Active)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	gateway.ID = id
	return nil
}

func UpdateIPFSGateway(gateway IPFSGateway) error {
	fetchHeadersBytes, err := json.Marshal(gateway.FetchHeaders)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		UPDATE ipfs_gateways
		SET name = ?, endpoint = ?, fetch_headers = ?, active = ?
		WHERE id = ?;
	`, gateway.Name, gateway.Endpoint, string(fetchHeadersBytes), gateway.Active, gateway.ID)
	return err
}

func DelIPFSGateway(id int64) error {
	_, err := DB.Exec(`
		DELETE FROM ipfs_gateways
		WHERE id = ?;
	`, id)
	return err
}

func IPFSFetch(cId string) (*http.Response, error) {
	gateways, err := GetIPFSGateways(GetIPFSGatewaysParams{
		Active: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	for _, gateway := range gateways {
		res, err := gateway.Fetch(cId)
		if err != nil {
			continue
		}

		if res.StatusCode != 200 {
			continue
		}

		return res, nil
	}

	return nil, fmt.Errorf("unavailable")
}

func (i IPFSGateway) Fetch(cId string) (*http.Response, error) {
	endpoint := strings.Replace(i.Endpoint, "{cid}", cId, -1)

	client := new(http.Client)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (i IPFSGateway) TestFetch() error {
	cId := "bafybeibozpulxtpv5nhfa2ue3dcjx23ndh3gwr5vwllk7ptoyfwnfjjr4q"
	res, err := i.Fetch(cId)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%d", res.StatusCode)
	}

	return nil
}
