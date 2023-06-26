package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gioui.org/app"
)

type Settings struct {
	AppDir     string `json:"-"`
	NodeDir    string `json:"-"`
	WalletsDir string `json:"-"`

	HideBalance  bool `json:"hide_balance"`
	SendRingSize int  `json:"send_ring_size"`
}

var Instance *Settings

// vars below are replaced by -ldflags during build
var Version = "development"
var BuildTime = ""
var GitVersion = "development"

func Instantiate() *Settings {
	Instance = &Settings{}
	return Instance
}

func (s *Settings) Load() error {
	dataDir, err := app.DataDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(dataDir, "g45w")
	nodeDir := filepath.Join(appDir, "node")
	walletsDir := filepath.Join(appDir, "wallets")

	s.AppDir = appDir
	s.NodeDir = nodeDir
	s.WalletsDir = walletsDir

	path := filepath.Join(s.AppDir, "settings.json")

	_, err = os.Stat(path)
	if err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Settings) Save() error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	path := filepath.Join(s.AppDir, "settings.json")
	return os.WriteFile(path, data, os.ModePerm)
}
