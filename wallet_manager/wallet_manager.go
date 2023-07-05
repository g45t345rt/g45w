package wallet_manager

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/settings"
)

type Wallet struct {
	Info   *WalletInfo
	Memory *walletapi.Wallet_Memory
}

type WalletInfo struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	//Data      []byte `json:"data"`
	Timestamp int64 `json:"timestamp"`
	ListOrder int   `json:"order"` // save item list ordering
}

var Wallets map[string]*WalletInfo
var OpenedWallet *Wallet

func Load() error {
	walletsDir := settings.WalletsDir
	Wallets = make(map[string]*WalletInfo)

	err := os.MkdirAll(walletsDir, os.ModePerm)
	if err != nil {
		return err
	}

	return filepath.Walk(walletsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if walletsDir == path {
			return nil
		}

		if info.IsDir() {
			addr := info.Name()

			path := filepath.Join(walletsDir, addr, "info.json")
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			walletInfo := &WalletInfo{}
			err = json.Unmarshal(data, walletInfo)
			if err != nil {
				return err
			}

			Wallets[walletInfo.Addr] = walletInfo
		}

		return nil
	})
}

func GetWalletInfo(addr string) (*WalletInfo, error) {
	walletInfo, ok := Wallets[addr]
	if !ok {
		return nil, fmt.Errorf("wallet [%s] does not exists", addr)
	}

	return walletInfo, nil
}

func OpenWallet(addr string, password string) (*walletapi.Wallet_Memory, *WalletInfo, error) {
	walletInfo, err := GetWalletInfo(addr)
	if err != nil {
		return nil, nil, err
	}

	walletsDir := settings.WalletsDir
	path := filepath.Join(walletsDir, addr, "wallet.db")
	walletData, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	wallet, err := walletapi.Open_Encrypted_Wallet_Memory(password, walletData)
	if err != nil {
		return nil, nil, err
	}

	return wallet, walletInfo, nil
}

func SetOpenWallet(memory *walletapi.Wallet_Memory, info *WalletInfo) {
	OpenedWallet = &Wallet{
		Memory: memory,
		Info:   info,
	}
}

func DeleteWallet(addr string, password string) error {
	_, _, err := OpenWallet(addr, password)
	if err != nil {
		return err
	}

	walletsDir := settings.WalletsDir
	path := filepath.Join(walletsDir, addr)
	err = os.RemoveAll(path)
	if err != nil {
		return err
	}

	delete(Wallets, addr)
	return nil
}

func CreateWalletFromPath(name string, password string, path string) error {
	wallet, err := walletapi.Open_Encrypted_Wallet(path, password)
	if err != nil {
		return err
	}

	return saveWallet(wallet.Wallet_Memory, name)
}

func CreateWalletFromSeed(name string, password string, seed string) error {
	wallet, err := walletapi.Create_Encrypted_Wallet_From_Recovery_Words_Memory(password, seed)
	if err != nil {
		return err
	}

	return saveWallet(wallet, name)
}

func CreateWalletFromHexSeed(name string, password, hexSeed string) error {
	seed, err := hex.DecodeString(hexSeed)
	if err != nil {
		return err
	}

	if len(seed) != 32 {
		return fmt.Errorf("hex seed must be 64 chars")
	}

	eSeed := new(crypto.BNRed).SetBytes(seed)
	wallet, err := walletapi.Create_Encrypted_Wallet_Memory(password, eSeed)
	if err != nil {
		return err
	}

	return saveWallet(wallet, name)
}

func CreateRandomWallet(name string, password string) error {
	wallet, err := walletapi.Create_Encrypted_Wallet_Random_Memory(password)
	if err != nil {
		return err
	}

	return saveWallet(wallet, name)
}

func RenameWallet(addr string, newName string) error {
	walletInfo := Wallets[addr]
	walletInfo.Name = newName
	return saveWalletInfo(addr, walletInfo)
}

func OrderWallet(addr string, newOrder int) error {
	walletInfo := Wallets[addr]
	walletInfo.ListOrder = newOrder
	return saveWalletInfo(addr, walletInfo)
}

func ChangePassword(addr string, password string, newPassword string) error {
	memory, _, err := OpenWallet(addr, password)
	if err != nil {
		return err
	}

	seed := memory.GetAccount().Keys.Secret
	newMemory, err := walletapi.Create_Encrypted_Wallet_Memory(newPassword, seed)
	if err != nil {
		return err
	}

	return saveWalletData(newMemory)
}

func saveWallet(wallet *walletapi.Wallet_Memory, name string) error {
	addr := wallet.GetAddress().String()
	walletInfo := &WalletInfo{
		Addr:      addr,
		Name:      name,
		Timestamp: time.Now().Unix(),
	}

	err := saveWalletInfo(addr, walletInfo)
	if err != nil {
		return err
	}

	err = saveWalletData(wallet)
	if err != nil {
		return err
	}

	return nil
}

func saveWalletInfo(addr string, walletInfo *WalletInfo) error {
	data, err := json.Marshal(walletInfo)
	if err != nil {
		return err
	}

	walletsDir := settings.WalletsDir

	path := filepath.Join(walletsDir, addr)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	path = filepath.Join(walletsDir, addr, "info.json")
	err = os.WriteFile(path, data, fs.ModePerm)
	if err != nil {
		return err
	}

	Wallets[addr] = walletInfo
	return nil
}

func saveWalletData(wallet *walletapi.Wallet_Memory) error {
	walletData := wallet.Get_Encrypted_Wallet()
	addr := wallet.GetAddress().String()

	walletsDir := settings.WalletsDir

	path := filepath.Join(walletsDir, addr)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	path = filepath.Join(walletsDir, addr, "wallet.db")
	err = os.WriteFile(path, walletData, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
