package wallet_manager

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	AskToStoreDEXTokens               bool `json:"ask_to_store_dex_tokens"`
	NotifyXSWDMobileBackgroundService bool `json:"notify_xswd_mobile_background_service"`
}

func (w *Wallet) settingsPath() string {
	return filepath.Join(w.FolderPath, "settings.json")
}

func (w *Wallet) LoadSettings() error {
	settingsPath := w.settingsPath()

	w.Settings = Settings{
		AskToStoreDEXTokens:               true,
		NotifyXSWDMobileBackgroundService: true,
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
