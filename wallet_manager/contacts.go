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
	_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS contacts (
				addr VARCHAR PRIMARY KEY,
				name VARCHAR UNIQUE,
				note VARCHAR,
				timestamp BIGINT,
				list_order INTEGER
			);
		`)
	return err
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

	var contact Contact
	for rows.Next() {
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

	return &contact, nil
}

func (w *Wallet) StoreContact(contact Contact) error {
	_, err := w.DB.Exec(`
		INSERT INTO contacts (addr,name,note,timestamp)
		VALUES (?,?,?,?)
		ON CONFLICT (addr) DO UPDATE SET
		name = excluded.name,
		note = excluded.note;
	`, contact.Addr, contact.Name, contact.Note, time.Now().UnixMilli())
	return err
}

func (w *Wallet) DelContact(addr string) error {
	_, err := w.DB.Exec(`
		DELETE FROM contacts
		WHERE addr = ?;
	`, addr)
	return err
}
