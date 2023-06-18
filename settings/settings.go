package settings

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"gioui.org/app"
)

type Settings struct {
	AppDir     string
	NodeDir    string
	WalletsDir string
}

var Instance *Settings

// vars below are replaced by -ldflags during build
var Version = "development"
var BuildTime = fmt.Sprint(time.Now().Unix())
var GitVersion = "development"

func Instantiate() *Settings {
	Instance = &Settings{}
	return Instance
}

func (s *Settings) Load() error {
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
