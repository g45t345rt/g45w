package wallet_manager

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Contact struct {
	Name      string
	Addr      string
	Note      string
	Timestamp int64
	ListOrder sql.NullInt32
}

func initDatabaseContacts(db *sql.DB) error {
	dbTx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = dbTx.Exec(`
			CREATE TABLE IF NOT EXISTS contacts (
				addr VARCHAR PRIMARY KEY,
				name VARCHAR UNIQUE,
				note VARCHAR,
				timestamp BIGINT,
				list_order INTEGER
			);
		`)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(dbTx)
}

type GetContactsParams struct {
}

func (w *Wallet) GetContacts(params GetContactsParams) ([]Contact, error) {
	query := sq.Select("*").From("contacts")

	rows, err := query.RunWith(w.DB).Query()
	if err != nil {
		return nil, err
	}

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		err = rows.Scan(
			&contact.Addr,
			&contact.Name,
			&contact.Note,
			&contact.Timestamp,
			&contact.ListOrder,
		)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (w *Wallet) GetContact(addr string) (*Contact, error) {
	query := sq.Select("*").From("contacts").Where(sq.Eq{"addr": addr})

	rows, err := query.RunWith(w.DB).Query()
	if err != nil {
		return nil, err
	}

	var contact *Contact
	for rows.Next() {
		contact = &Contact{}
		err := rows.Scan(
			&contact.Addr,
			&contact.Name,
			&contact.Note,
			&contact.Timestamp,
			&contact.ListOrder,
		)
		if err != nil {
			return nil, err
		}
	}

	return contact, nil
}

func (w *Wallet) StoreContact(contact Contact) error {
	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO contacts (addr,name,note,timestamp)
		VALUES (?,?,?,?)
		ON CONFLICT (addr) DO UPDATE SET
		name = excluded.name,
		note = excluded.note;
	`, contact.Addr, contact.Name, contact.Note, time.Now().UnixMilli())
	if err != nil {
		return err
	}

	return handleDatabaseCommit(tx)
}

func (w *Wallet) DelContact(addr string) error {
	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM contacts
		WHERE addr = ?;
	`, addr)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(tx)
}
