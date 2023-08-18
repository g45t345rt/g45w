package app_db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/g45t345rt/g45w/settings"
)

var DB *sql.DB

func Load() error {
	appDir := settings.AppDir
	dbPath := filepath.Join(appDir, "app.db")

	firstLoad := false
	_, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			firstLoad = true
		} else {
			return err
		}
	}

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	initDatabaseNodes()
	initDatabaseIPFSGateways()

	if firstLoad {
		err = ResetNodeConnections()
		if err != nil {
			return err
		}

		err = ResetIPFSGateways()
		if err != nil {
			return err
		}
	}

	return err
}
