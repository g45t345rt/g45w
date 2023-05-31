package wallet_manager

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/settings"
)

type WalletInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Addr string `json:"addr"`
	Data []byte `json:"data"`
	Path string
}

var Instance *WalletManager

type WalletManager struct {
	Wallets map[string]*WalletInfo
}

func NewWalletManager() *WalletManager {
	w := &WalletManager{
		Wallets: make(map[string]*WalletInfo),
	}
	Instance = w
	return Instance
}

func (w *WalletManager) LoadWallets() error {
	walletsDir := settings.Instance.WalletsDir
	w.Wallets = make(map[string]*WalletInfo)

	err := filepath.WalkDir(walletsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		walletInfo := &WalletInfo{}
		err = json.Unmarshal(data, walletInfo)
		if err != nil {
			return err
		}

		walletInfo.Path = path
		w.Wallets[walletInfo.ID] = walletInfo
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (w *WalletManager) GetWalletInfo(id string) (*WalletInfo, error) {
	walletInfo, ok := w.Wallets[id]
	if !ok {
		return nil, fmt.Errorf("wallet [%s] does not exists", id)
	}

	return walletInfo, nil
}

func (w *WalletManager) OpenWallet(id string, password string) (*walletapi.Wallet_Memory, *WalletInfo, error) {
	walletInfo, err := w.GetWalletInfo(id)
	if err != nil {
		return nil, nil, err
	}

	wallet, err := walletapi.Open_Encrypted_Wallet_Memory(password, walletInfo.Data)
	if err != nil {
		return nil, nil, err
	}

	return wallet, walletInfo, nil
}

func (w *WalletManager) DeleteWallet(id string, password string) error {
	_, walletInfo, err := w.OpenWallet(id, password)
	if err != nil {
		return err
	}

	err = os.Remove(walletInfo.Path)
	if err != nil {
		return err
	}

	delete(w.Wallets, id)
	return nil
}

func (w *WalletManager) CreateWallet(name string, password string, confirmPassword string) error {
	if password != confirmPassword {
		return fmt.Errorf("password don't match")
	}

	walletsDir := settings.Instance.WalletsDir
	wallet, err := walletapi.Create_Encrypted_Wallet_Random_Memory(password)
	if err != nil {
		return err
	}

	walletData := wallet.Get_Encrypted_Wallet()

	id := fmt.Sprint(time.Now().Unix())
	walletInfo := &WalletInfo{
		ID:   id,
		Name: name,
		Addr: wallet.GetAddress().String(),
		Data: walletData,
	}

	data, err := json.Marshal(walletInfo)
	if err != nil {
		return err
	}

	err = os.MkdirAll(walletsDir, fs.ModePerm)
	if err != nil {
		return err
	}

	walletPath := filepath.Join(walletsDir, fmt.Sprintf("%s.json", id))
	err = os.WriteFile(walletPath, data, fs.ModePerm)

	w.Wallets[id] = walletInfo
	return err
}
