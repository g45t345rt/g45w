package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gioui.org/app"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/theme"
)

type AppSettings struct {
	LanguageKey  string `json:"language"`
	HideBalance  bool   `json:"hide_balance"`
	SendRingSize int    `json:"send_ring_size"`
	NodeEndpoint string `json:"node_endpoint"`
	TabBarsKey   string `json:"tab_bars_key"`
	ThemeKey     string `json:"theme"`
}

var (
	AppDir            string
	IntegratedNodeDir string
	WalletsDir        string
	CacheDir          string
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
	cacheDir := filepath.Join(appDir, "cache")

	AppDir = appDir
	IntegratedNodeDir = integratedNodeDir
	WalletsDir = walletsDir
	CacheDir = cacheDir

	settingsPath := filepath.Join(AppDir, "settings.json")

	// settings with default values
	appSettings := AppSettings{
		LanguageKey:  "en",
		HideBalance:  false,
		SendRingSize: 16,
		NodeEndpoint: "",
		TabBarsKey:   "tokens",
		ThemeKey:     "light",
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

	language := lang.Get(appSettings.LanguageKey)
	if language != nil {
		lang.Current = language.Key
	} else {
		lang.Current = "en"
		appSettings.LanguageKey = "en"
	}

	currentTheme := theme.Get(appSettings.ThemeKey)
	if currentTheme != nil {
		theme.Current = *currentTheme
	} else {
		theme.Current = theme.Light
		appSettings.ThemeKey = "light"
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
