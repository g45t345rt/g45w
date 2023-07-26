package app_data

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type NodeConnection struct {
	ID         int64
	Endpoint   string
	Name       string
	Integrated bool
}

var INTEGRATED_NODE_CONNECTION = NodeConnection{
	ID:         -1,
	Endpoint:   "ws://127.0.0.1:10102/ws",
	Name:       "Integrated",
	Integrated: true,
}

var TRUSTED_NODE_CONNECTIONS = []NodeConnection{
	{Endpoint: "wss://node.deronfts.com/ws", Name: "DeroNFTs"},
	{Endpoint: "wss://dero-node.mysrv.cloud/ws", Name: "MySrvCloud"},
	{Endpoint: "ws://derostats.io:10102/ws", Name: "DeroStats"},
	{Endpoint: "ws://node.derofoundation.org:11012/ws", Name: "Foundation"},
	{Endpoint: "ws://wallet.friendspool.club:10102/ws", Name: "Friendspool"},
}

func initDatabaseNodes() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS nodes (
			id INTEGER PRIMARY KEY,
			endpoint VARCHAR,
			name VARCHAR
		);
	`)
	return err
}

func StoreTrustedNodeConnections() error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	for _, nodeConn := range TRUSTED_NODE_CONNECTIONS {
		_, err = tx.Exec(`
			INSERT INTO nodes (endpoint, name)
			VALUES (?,?);
		`, nodeConn.Endpoint, nodeConn.Name, nodeConn.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetNodeConnections() ([]NodeConnection, error) {
	query := sq.Select("*").From("nodes")

	rows, err := query.RunWith(DB).Query()
	if err != nil {
		return nil, err
	}

	var nodes []NodeConnection
	for rows.Next() {
		var node NodeConnection
		err = rows.Scan(
			&node.ID,
			&node.Endpoint,
			&node.Name,
		)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, err
}

func GetNodeConnection(endpoint string) (*NodeConnection, error) {
	query := sq.Select("*").From("nodes").Where(sq.Eq{"endpoint": endpoint})

	row := query.RunWith(DB).QueryRow()

	var nodeConn NodeConnection
	err := row.Scan(
		&nodeConn.ID,
		&nodeConn.Endpoint,
		&nodeConn.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &nodeConn, nil
}

func InsertNodeConnection(nodeConn NodeConnection) error {
	_, err := DB.Exec(`
		INSERT INTO nodes (endpoint,name)
		VALUES (?,?);
	`, nodeConn.Endpoint, nodeConn.Name)
	return err
}

func UpdateNodeConnection(nodeConn NodeConnection) error {
	_, err := DB.Exec(`
		UPDATE nodes
		SET name = ?, endpoint = ?
		WHERE id = ?;
	`, nodeConn.Name, nodeConn.Endpoint, nodeConn.ID)
	return err
}

func DelNodeConnection(id int64) error {
	_, err := DB.Exec(`
		DELETE FROM nodes
		WHERE id = ?;
	`, id)
	return err
}

func ClearNodeConnections() error {
	_, err := DB.Exec(`
		DELETE FROM nodes
	`)
	return err
}

func GetNodeCount() (int, error) {
	query := sq.Select("COUNT(*)").From("nodes")

	row := query.RunWith(DB).QueryRow()

	var count int
	err := row.Scan(&count)
	return count, err
}
