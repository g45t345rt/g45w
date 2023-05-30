package settings

import (
	"log"
	"path/filepath"

	"gioui.org/app"
)

type Settings struct {
	AppDir  string
	NodeDir string
}

func LoadSettings() (*Settings, error) {
	dataDir, err := app.DataDir()
	if err != nil {
		log.Fatal(err)
	}

	appDir := filepath.Join(dataDir, "g45w")
	nodeDir := filepath.Join(appDir, "node")

	settings := &Settings{
		AppDir:  appDir,
		NodeDir: nodeDir,
	}

	return settings, nil
}

func (s *Settings) Save() {

}
