package order_column

import (
	"database/sql"
	"fmt"
	"strings"
)

// Be careful with TableName,ColumnName they are NOT sql injection safe
type Orderer struct {
	TableName   string
	ColumnName  string
	FilterQuery string
}

func (o Orderer) applyFilterQuery(query string, prefix string) string {
	if o.FilterQuery != "" {
		return strings.Replace(query, "{filter_query}", prefix+o.FilterQuery, 1)
	}
	return strings.Replace(query, "{filter_query}", "", 1)
}

func (o Orderer) GetNewOrderNumber(tx *sql.Tx) (orderNumber int, err error) {
	query := fmt.Sprintf(`
	SELECT %s
	FROM %s
	{filter_query}
	ORDER BY %s DESC
	LIMIT 1;
`, o.ColumnName, o.TableName, o.ColumnName)

	query = o.applyFilterQuery(query, "WHERE ")
	row := tx.QueryRow(query)

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
	query := fmt.Sprintf(`
	UPDATE %s
	SET %s = %s + 1
	WHERE %s >= ? {filter_query};
`, o.TableName, o.ColumnName, o.ColumnName, o.ColumnName)

	query = o.applyFilterQuery(query, "AND ")
	_, err := tx.Exec(query, orderNumber)
	if err != nil {
		return err
	}

	return nil
}

func (o Orderer) Update(tx *sql.Tx, currentOrderNumber int, newOrderNumber int) error {
	if newOrderNumber > currentOrderNumber {
		query := fmt.Sprintf(`
		UPDATE %s
		SET %s = %s - 1
		WHERE %s >= ? AND %s <= ? {filter_query};
		`, o.TableName, o.ColumnName, o.ColumnName, o.ColumnName, o.ColumnName)

		query = o.applyFilterQuery(query, "AND ")
		_, err := tx.Exec(query, currentOrderNumber, newOrderNumber)
		if err != nil {
			return err
		}

		return nil
	} else {
		query := fmt.Sprintf(`
		UPDATE %s
		SET %s = %s + 1
		WHERE %s >= ? AND %s <= ? {filter_query};
		`, o.TableName, o.ColumnName, o.ColumnName, o.ColumnName, o.ColumnName)
		query = o.applyFilterQuery(query, "AND ")

		_, err := tx.Exec(query, newOrderNumber, currentOrderNumber)
		if err != nil {
			return err
		}

		return nil
	}
}

func (o Orderer) Delete(tx *sql.Tx, orderNumber int) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET %s = %s - 1
	WHERE %s >= ? {filter_query};
`, o.TableName, o.ColumnName, o.ColumnName, o.ColumnName)
	query = o.applyFilterQuery(query, "AND ")
	_, err := tx.Exec(query, orderNumber)
	if err != nil {
		return err
	}

	return nil
}
