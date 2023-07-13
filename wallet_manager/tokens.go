package wallet_manager

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/g45t345rt/g45w/sc"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/utils"
)

type TokenInfo struct {
	SCID      string    `json:"scid"`
	Type      sc.SCType `json:"type"`
	JsonData  string    `json:"data"`
	ListOrder int       `json:"order"`
}

var MAIN_FOLDER = "__main__"

func (w *Wallet) TokensFolderPath() string {
	walletDir := settings.WalletsDir
	tokensPath := filepath.Join(walletDir, w.Info.Addr, "tokens")
	return tokensPath
}

func (w *Wallet) GetTokens(folder string) ([]interface{}, error) {
	tokensFolder := filepath.Join(w.TokensFolderPath(), folder)

	var tokens []interface{}
	err := filepath.Walk(tokensFolder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var token interface{}
		err = json.Unmarshal(data, &token)
		if err != nil {
			return err
		}

		tokens = append(tokens, token)
		return nil
	})

	return tokens, err
}

func (w *Wallet) GetFolders() ([]string, error) {
	var folders []string
	err := filepath.Walk(w.TokensFolderPath(), func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			folders = append(folders, info.Name())
		}

		return nil
	})

	return folders, err
}

func (w *Wallet) AddToken(folder string, scId string, data interface{}) error {
	path := filepath.Join(w.TokensFolderPath(), folder, fmt.Sprintf("%s.json", scId))

	dataString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(dataString), os.ModePerm)
}

func (w *Wallet) MoveToken(path string, newPath string) error {
	err := utils.CopyFile(path, newPath)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (w *Wallet) DupToken(path string, newPath string) error {
	return utils.CopyFile(path, newPath)
}

func (w *Wallet) DelToken(path string) error {
	return os.Remove(path)
}

func (w *Wallet) AddFolder(name string) error {
	folderPath := filepath.Join(w.TokensFolderPath(), name)
	return os.Mkdir(folderPath, os.ModePerm)
}

func (w *Wallet) DelFolder(name string) error {
	folderPath := filepath.Join(w.TokensFolderPath(), name)
	return os.RemoveAll(folderPath)
}

func (w *Wallet) RenameFolder(name string, newName string) error {
	oldPath := filepath.Join(w.TokensFolderPath(), name)
	newPath := filepath.Join(w.TokensFolderPath(), newName)
	return os.Rename(oldPath, newPath)
}

func (w *Wallet) ImportTokensFromWallet(walletAddr string) error {
	walletDir := settings.WalletsDir
	tokensFolder := filepath.Join(walletDir, walletAddr)

	files, err := os.ReadDir(tokensFolder)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		sourceFolder := w.TokensFolderPath()
		sourcePath := filepath.Join(sourceFolder, fileInfo.Name())
		destPath := filepath.Join(w.TokensFolderPath(), fileInfo.Name())
		err := utils.CopyFile(sourcePath, destPath)
		if err != nil {
			return err
		}
	}

	return nil
}
