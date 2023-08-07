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
	TabBarsKey   string `json:"tab_bars_key"`
}

var (
	AppDir            string
	IntegratedNodeDir string
	WalletsDir        string
)

var App AppSettings
var Name = "G45W"

// vars below are replaced by -ldflags during build
var Version = "development"
var BuildTime = ""
var GitVersion = "development"
var DonationAddress = "dero1qyhunyuk24g9qsjtcr4r0c7rgjquuernqcfnx76kq0jvn4ns98tf2qgj5dq70"

func Load() error {
	dataDir, err := app.DataDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(dataDir, "g45w")
	_, err = os.Stat(appDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(appDir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	integratedNodeDir := filepath.Join(appDir, "node")
	walletsDir := filepath.Join(appDir, "wallets")

	AppDir = appDir
	IntegratedNodeDir = integratedNodeDir
	WalletsDir = walletsDir

	settingsPath := filepath.Join(AppDir, "settings.json")

	// settings with default values
	appSettings := AppSettings{
		Language:     "en",
		HideBalance:  false,
		SendRingSize: 16,
		NodeEndpoint: "",
		TabBarsKey:   "tokens",
	}

	_, err = os.Stat(settingsPath)
	if err == nil {
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &appSettings)
		if err != nil {
			return err
		}
	}

	App = appSettings
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
