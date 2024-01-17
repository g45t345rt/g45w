package wallet_manager

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	AskedToStoreDEXTokens bool
}

func (w *Wallet) settingsPath() string {
	return filepath.Join(w.FolderPath, "settings.json")
}

func (w *Wallet) LoadSettings() error {
	settingsPath := w.settingsPath()

	w.Settings = Settings{
		AskedToStoreDEXTokens: false,
	}

	_, err := os.Stat(settingsPath)
	if err == nil {
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &w.Settings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Wallet) SaveSettings() error {
	data, err := json.MarshalIndent(w.Settings, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(w.settingsPath(), data, os.ModePerm)
}
