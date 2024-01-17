package app_db

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/deroproject/derohe/globals"
	"github.com/g45t345rt/g45w/app_db/order_column"
	"github.com/g45t345rt/g45w/app_db/schema_version"
	"github.com/g45t345rt/g45w/settings"
)

type NodeConnection struct {
	ID          int64
	Endpoint    string
	Name        string
	Integrated  bool
	OrderNumber int
}

func GetIntegratedNode() NodeConnection {
	endpoint := fmt.Sprintf("ws://127.0.0.1:%d/ws", globals.Config.RPC_Default_Port)
	return NodeConnection{
		ID:         -1,
		Endpoint:   endpoint,
		Name:       "Integrated",
		Integrated: true,
	}
}

func GetLocalNode() NodeConnection {
	endpoint := fmt.Sprintf("ws://127.0.0.1:%d/ws", globals.Config.RPC_Default_Port)
	return NodeConnection{
		ID:         -2,
		Endpoint:   endpoint,
		Name:       "Local",
		Integrated: false,
	}
}

func GetTrustedRemoveNodes() []NodeConnection {
	if settings.App.Testnet {
		return []NodeConnection{}
	}

	return []NodeConnection{
		{Endpoint: "wss://node.deronfts.com/ws", Name: "DeroNFTs"},
		{Endpoint: "wss://dero-node.mysrv.cloud/ws", Name: "MySrvCloud"},
		{Endpoint: "ws://derostats.io:10102/ws", Name: "DeroStats"},
		{Endpoint: "ws://node.derofoundation.org:11012/ws", Name: "Foundation"},
		{Endpoint: "ws://wallet.friendspool.club:10102/ws", Name: "Friendspool"},
	}
}

var nodeOrderer = order_column.Orderer{
	TableName:  "nodes",
	ColumnName: "order_number",
}

func initTableNodes() error {
	version, err := schema_version.GetVersion(DB, "nodes")
	if err != nil {
		return err
	}

	if version == 0 {
		_, err := DB.Exec(`
			CREATE TABLE IF NOT EXISTS nodes (
				id INTEGER PRIMARY KEY,
				endpoint VARCHAR,
				name VARCHAR
			);
		`)
		if err != nil {
			return err
		}

		_, err = DB.Exec(`
			ALTER TABLE nodes
			ADD COLUMN order_number INT NOT NULL DEFAULT 0;
		`)
		if err != nil {
			return err
		}

		err = ReOrderNodes(DB)
		if err != nil {
			return err
		}

		version = 1
		err = schema_version.StoreVersion(DB, "nodes", version)
		if err != nil {
			return err
		}
	}

	return err
}

func ResetNodeConnections() error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM nodes`)
	if err != nil {
		tx.Rollback()
		return err
	}

	for i, node := range GetTrustedRemoveNodes() {
		_, err = tx.Exec(`
			INSERT INTO nodes (endpoint, name, order_number)
			VALUES (?,?,?);
		`, node.Endpoint, node.Name, i)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetNodeConnections() ([]NodeConnection, error) {
	query := sq.Select("*").From("nodes").OrderBy("order_number ASC")

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
			&node.OrderNumber,
		)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, err
}

func GetNodeConnection(id int64) (node NodeConnection, err error) {
	row := DB.QueryRow(`
		SELECT * FROM nodes
		WHERE id = ?;
	`, id)

	err = row.Scan(
		&node.ID,
		&node.Endpoint,
		&node.Name,
		&node.OrderNumber,
	)
	if err != nil {
		return
	}

	return
}

func GetNodeConnectionByEndpoint(endpoint string) (node NodeConnection, err error) {
	query := sq.Select("*").From("nodes").Where(sq.Eq{"endpoint": endpoint})

	row := query.RunWith(DB).QueryRow()

	err = row.Scan(
		&node.ID,
		&node.Endpoint,
		&node.Name,
		&node.OrderNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}

		return
	}

	return
}

func InsertNodeConnection(node NodeConnection) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	node.OrderNumber, err = nodeOrderer.GetNewOrderNumber(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO nodes (endpoint,name,order_number)
		VALUES (?,?,?);
	`, node.Endpoint, node.Name, node.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func UpdateNodeConnection(node NodeConnection) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	currentNode, err := GetNodeConnection(node.ID)
	if err != nil {
		return err
	}

	err = nodeOrderer.Update(tx, currentNode.OrderNumber, node.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE nodes
		SET name = ?, endpoint = ?, order_number = ?
		WHERE id = ?;
	`, node.Name, node.Endpoint, node.OrderNumber, node.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DelNodeConnection(id int64) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	nodeConnection, err := GetNodeConnection(id)
	if err != nil {
		return err
	}

	err = nodeOrderer.Delete(tx, nodeConnection.OrderNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM nodes
		WHERE id = ?;
	`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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

func ReOrderNodes(db *sql.DB) error {
	nodes, err := GetNodeConnections()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for i, node := range nodes {
		_, err = tx.Exec(`
				UPDATE nodes
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
