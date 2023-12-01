package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gioui.org/app"
	sysTheme "gioui.org/x/pref/theme"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/theme"
)

var (
	MainTabBarsToken = "tokens"
	MainTabBarsTxs   = "txs"
	FolderLayoutGrid = "grid"
	FolderLayoutList = "list"
)

type AppSettings struct {
	Language     string `json:"language"`
	HideBalance  bool   `json:"hide_balance"`
	SendRingSize int    `json:"send_ring_size"`
	NodeEndpoint string `json:"node_endpoint"`
	MainTabBars  string `json:"main_tab_bars"`
	Theme        string `json:"theme"`
	FolderLayout string `json:"folder_layout"`
}

var (
	AppDir            string
	IntegratedNodeDir string
	WalletsDir        string
	CacheDir          string
)

var App AppSettings
var Name = "secret-wallet"

// vars below are replaced by -ldflags during build
var Version = "secret-dev"
var BuildTime = ""
var GitVersion = "secret-dev"

func Load() error {
	dataDir, err := app.DataDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(dataDir, "secret-wallet")
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
		Language:     "en",
		HideBalance:  false,
		SendRingSize: 16,
		NodeEndpoint: "",
		MainTabBars:  MainTabBarsTxs,
		FolderLayout: FolderLayoutGrid,
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

	language := lang.Get(appSettings.Language)
	if language != nil {
		lang.Current = language.Key
	} else {
		lang.Current = "en"
		appSettings.Language = "en"
	}

	currentTheme := theme.Get(appSettings.Theme)
	if currentTheme != nil {
		theme.Current = currentTheme
	} else {
		// check system user theme preference
		isDark, _ := sysTheme.IsDarkMode()
		if isDark {
			appSettings.Theme = "dark"
		} else {
			appSettings.Theme = "light"
		}

		theme.Current = theme.Get(appSettings.Theme)
	}

	App = appSettings
	return nil
}

func Save() error {
	data, err := json.MarshalIndent(App, "", " ")
	if err != nil {
		return err
	}

	path := filepath.Join(AppDir, "settings.json")
	return os.WriteFile(path, data, os.ModePerm)
}
