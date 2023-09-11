package order_column

import (
	"database/sql"
	"fmt"
)

// Be careful with TableName,ColumnName they are NOT sql injection safe
type Orderer struct {
	TableName  string
	ColumnName string
}

func (o Orderer) GetNewOrderNumber(tx *sql.Tx) (orderNumber int, err error) {
	row := tx.QueryRow(fmt.Sprintf(`
	SELECT %s
	FROM %s
	ORDER BY %s DESC
	LIMIT 1;
`, o.ColumnName, o.TableName, o.ColumnName))

	err = row.Err()
	if err != nil {
		return
	}

	err = row.Scan(&orderNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil // return 0 as first ordering number
		}
		return
	}

	orderNumber++
	return
}

func (o Orderer) Insert(tx *sql.Tx, orderNumber int) error {
	_, err := tx.Exec(fmt.Sprintf(`
	UPDATE %s
	SET %s = %s + 1
	WHERE %s >= ?;
`, o.TableName, o.ColumnName, o.ColumnName, o.ColumnName), orderNumber)
	if err != nil {
		return err
	}

	return nil
}

func (o Orderer) Update(tx *sql.Tx, currentOrderNumber int, newOrderNumber int) error {
	if newOrderNumber > currentOrderNumber {
		_, err := tx.Exec(fmt.Sprintf(`
		UPDATE %s
		SET %s = %s - 1
		WHERE %s >= ? AND %s <= ?;
		`, o.TableName, o.ColumnName, o.ColumnName, o.ColumnName, o.ColumnName), currentOrderNumber, newOrderNumber)
		if err != nil {
			return err
		}

		return nil
	} else {
		_, err := tx.Exec(fmt.Sprintf(`
		UPDATE %s
		SET %s = %s + 1
		WHERE %s >= ? AND %s <= ?;
		`, o.TableName, o.ColumnName, o.ColumnName, o.ColumnName, o.ColumnName), newOrderNumber, currentOrderNumber)
		if err != nil {
			return err
		}

		return nil
	}
}

func (o Orderer) Delete(tx *sql.Tx, orderNumber int) error {
	_, err := tx.Exec(fmt.Sprintf(`
	UPDATE %s
	SET %s = %s - 1
	WHERE %s >= ?;
`, o.TableName, o.ColumnName, o.ColumnName, o.ColumnName), orderNumber)
	if err != nil {
		return err
	}

	return nil
}
