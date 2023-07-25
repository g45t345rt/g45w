package app_data

import (
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
)

type IPFSGateway struct {
	ID           int64
	Name         string
	Endpoint     string
	FetchHeaders map[string]string
}

func initDatabaseIPFSGateways() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS ipfs_gateways (
			id INTEGER AUTOINCREMENT,
			name VARCHAR,
			endpoint VARCHAR,
			fetch_headers VARCHAR
		);
	`)
	return err
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

		var fetchHeaders string
		err = rows.Scan(
			&ipfsGateway.Name,
			&ipfsGateway.Endpoint,
			&fetchHeaders,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(fetchHeaders), &ipfsGateway.FetchHeaders)
		if err != nil {
			return nil, err
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
