package settings

import (
	"log"
	"path/filepath"

	"gioui.org/app"
)

type Settings struct {
	AppDir     string
	NodeDir    string
	WalletsDir string
}

var Instance *Settings

var Version = "" // replaced by -ldflags
var BuildTime = ""
var GitVersion = ""

func NewSettings() *Settings {
	s := &Settings{}
	Instance = s
	return s
}

func (s *Settings) LoadSettings() error {
	dataDir, err := app.DataDir()
	if err != nil {
		log.Fatal(err)
	}

	appDir := filepath.Join(dataDir, "g45w")
	nodeDir := filepath.Join(appDir, "node")
	walletsDir := filepath.Join(appDir, "wallets")

	s.AppDir = appDir
	s.NodeDir = nodeDir
	s.WalletsDir = walletsDir
	return nil
}

func (s *Settings) Save() {

}
