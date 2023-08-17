package schema_version

import "database/sql"

func Init(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			schema VARCHAR PRIMARY KEY,
			version INTEGER
		);
	`)
	return err
}

func GetVersion(db *sql.DB, schema string) (version int, err error) {
	row := db.QueryRow(`
		SELECT version FROM schema_version
		WHERE schema = ?;
	`, schema)

	err = row.Scan(&version)
	if err == sql.ErrNoRows {
		err = nil // default version = 0
	}

	return
}

func StoreVersion(db *sql.DB, schema string, version int) error {
	_, err := db.Exec(`
		INSERT INTO schema_version (schema, version)
		VALUES (?,?)
		ON CONFLICT (schema) DO UPDATE SET
		version = excluded.version;
	`, schema, version)
	return err
}
