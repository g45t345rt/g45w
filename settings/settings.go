package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gioui.org/app"
)

type AppSettings struct {
	Language     string `json:"language"`
	HideBalance  bool   `json:"hide_balance"`
	SendRingSize int    `json:"send_ring_size"`
	NodeEndpoint string `json:"node_endpoint"`
}

var (
	AppDir            string
	IntegratedNodeDir string
	WalletsDir        string
)

var App AppSettings

// vars below are replaced by -ldflags during build
var Version = "development"
var BuildTime = ""
var GitVersion = "development"

func Load() error {
	dataDir, err := app.DataDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(dataDir, "g45w")
	integratedNodeDir := filepath.Join(appDir, "node")
	walletsDir := filepath.Join(appDir, "wallets")

	AppDir = appDir
	IntegratedNodeDir = integratedNodeDir
	WalletsDir = walletsDir

	path := filepath.Join(AppDir, "settings.json")

	_, err = os.Stat(path)
	if err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &App)
		if err != nil {
			return err
		}
	}

	return nil
}

func Save() error {
	data, err := json.Marshal(App)
	if err != nil {
		return err
	}

	path := filepath.Join(AppDir, "settings.json")
	return os.WriteFile(path, data, os.ModePerm)
}
