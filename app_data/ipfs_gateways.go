package app_data

import (
	"database/sql"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
)

type IPFSGateway struct {
	ID           int64
	Name         string
	Endpoint     string
	FetchHeaders map[string]string
	Use          bool
}

var TRUSTED_IPFS_GATEWAYS = []IPFSGateway{
	{Name: "deronfts.com", Endpoint: "https://ipfs.deronfts.com/ipfs/"},
	{Name: "ipfs.io", Endpoint: "https://ipfs.io/ipfs/"},
	{Name: "pinata.cloud", Endpoint: "https://gateway.pinata.cloud/ipfs/"},
	{Name: "dweb.link", Endpoint: "https://dweb.link/ipfs/"},
	{Name: "nftstorage.link", Endpoint: "https://nftstorage.link/ipfs"},
}

func initDatabaseIPFSGateways() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS ipfs_gateways (
			endpoint VARCHAR PRIMARY KEY,
			name VARCHAR,
			fetch_headers VARCHAR,
			use BOOL
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
			INSERT INTO ipfs_gateways (endpoint, name, use)
			VALUES (?,?,?)
			ON CONFLICT (endpoint) DO UPDATE SET
			name = ?;
		`, gateway.Endpoint, gateway.Name, true, gateway.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetIPFSGateways() ([]IPFSGateway, error) {
	sq := sq.Select("*").From("ipfs_gateways")

	rows, err := sq.RunWith(DB).Query()
	if err != nil {
		return nil, err
	}

	var ipfsGateways []IPFSGateway
	for rows.Next() {
		var ipfsGateway IPFSGateway

		var fetchHeaders sql.NullString
		err = rows.Scan(
			&ipfsGateway.Endpoint,
			&ipfsGateway.Name,
			&fetchHeaders,
			&ipfsGateway.Use,
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
		INSERT INTO node_list (name,endpoint,fetch_headers)
		VALUES (?,?,?);
	`, gateway.Name, gateway.Endpoint, string(fetchHeadersBytes))
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

func UpdateIPFSGatway(gateway IPFSGateway) error {
	fetchHeadersBytes, err := json.Marshal(gateway.FetchHeaders)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		UPDATE node_list
		SET name = ?, endpoint = ?, fetch_headers = ?
		WHERE id = ?;
	`, gateway.Name, gateway.Endpoint, string(fetchHeadersBytes), gateway.ID)
	return err
}

func DelIPFSGateway(id int64) error {
	_, err := DB.Exec(`
		DELETE FROM ipfs_gateways
		WHERE id = ?;
	`, id)
	return err
}
