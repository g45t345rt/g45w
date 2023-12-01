package app_db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/g45t345rt/g45w/app_db/schema_version"
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

	err = schema_version.Init(DB)
	if err != nil {
		return err
	}

	err = initDatabaseNodes()
	if err != nil {
		return err
	}

	err = initDatabaseWallets()
	if err != nil {
		return err
	}

	err = delWalletInfoIfNoFolder()
	if err != nil {
		fmt.Println(err)
	}

	if firstLoad {
		err = ResetNodeConnections()
		if err != nil {
			return err
		}

	}

	return err
}
