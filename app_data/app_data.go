package app_data

import (
	"database/sql"
	"path/filepath"

	"github.com/g45t345rt/g45w/settings"
)

var DB *sql.DB

func Load() error {
	appDir := settings.AppDir
	dbPath := filepath.Join(appDir, "app.db")

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	initDatabaseNodes()
	initDatabaseIPFSGateways()
	return err
}
