package token_manager

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

type TokenManager struct {
	WalletAddr string
}

var MAIN_FOLDER = "__main__"

func New(walletAddr string) *TokenManager {
	return &TokenManager{WalletAddr: walletAddr}
}

func (t *TokenManager) tokensFolder() string {
	walletDir := settings.WalletsDir
	tokensPath := filepath.Join(walletDir, t.WalletAddr, "tokens")
	return tokensPath
}

func (t *TokenManager) GetTokens(folder string) ([]interface{}, error) {
	tokensFolder := filepath.Join(t.tokensFolder(), folder)

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

func (t *TokenManager) GetFolders() ([]string, error) {
	var folders []string
	err := filepath.Walk(t.tokensFolder(), func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			folders = append(folders, info.Name())
		}

		return nil
	})

	return folders, err
}

func (t *TokenManager) AddToken(folder string, scId string, data interface{}) error {
	path := filepath.Join(t.tokensFolder(), folder, fmt.Sprintf("%s.json", scId))

	dataString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(dataString), os.ModePerm)
}

func (t *TokenManager) MoveToken(path string, newPath string) error {
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

func (t *TokenManager) DupToken(path string, newPath string) error {
	return utils.CopyFile(path, newPath)
}

func (t *TokenManager) DelToken(path string) error {
	return os.Remove(path)
}

func (t *TokenManager) AddFolder(name string) error {
	folderPath := filepath.Join(t.tokensFolder(), name)
	return os.Mkdir(folderPath, os.ModePerm)
}

func (t *TokenManager) DelFolder(name string) error {
	folderPath := filepath.Join(t.tokensFolder(), name)
	return os.RemoveAll(folderPath)
}

func (t *TokenManager) RenameFolder(name string, newName string) error {
	oldPath := filepath.Join(t.tokensFolder(), name)
	newPath := filepath.Join(t.tokensFolder(), newName)
	return os.Rename(oldPath, newPath)
}

func (t *TokenManager) ImportTokensFromWallet(walletAddr string) error {
	walletDir := settings.WalletsDir
	tokensFolder := filepath.Join(walletDir, walletAddr)

	files, err := os.ReadDir(tokensFolder)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		sourceFolder := t.tokensFolder()
		sourcePath := filepath.Join(sourceFolder, fileInfo.Name())
		destPath := filepath.Join(t.tokensFolder(), fileInfo.Name())
		err := utils.CopyFile(sourcePath, destPath)
		if err != nil {
			return err
		}
	}

	return nil
}
