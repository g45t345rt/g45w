package wallet_manager

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/creachadair/jrpc2"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/transaction"
	"github.com/deroproject/derohe/walletapi"
	"github.com/deroproject/derohe/walletapi/xswd"
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/app_db/schema_version"
	"github.com/g45t345rt/g45w/sc"
	"github.com/g45t345rt/g45w/settings"

	"database/sql"

	_ "modernc.org/sqlite" // use CGO free port
)

type Wallet struct {
	Info       app_db.WalletInfo
	Memory     *walletapi.Wallet_Disk
	DB         *sql.DB
	ServerXSWD *xswd.XSWD
	FolderPath string
	Settings   Settings
}

var OpenedWallet *Wallet

func CloseOpenedWallet() {
	if OpenedWallet != nil {
		wallet := OpenedWallet
		go func() {
			close(wallet.Memory.Quit) // make sure to close goroutines when wallet is in online mode
			wallet.Memory.Close_Encrypted_Wallet()
		}()
		wallet.CloseXSWD()
		wallet.DB.Close()
		OpenedWallet = nil
	}
}

func OpenWallet(addr string, password string) error {
	walletInfo, err := app_db.GetWalletInfo(addr)
	if err != nil {
		return err
	}

	walletsDir := settings.WalletsDir
	folderPath := filepath.Join(walletsDir, addr)
	walletPath := filepath.Join(folderPath, "wallet.db")

	bkCopied := false
open_wallet:
	memory, err := walletapi.Open_Encrypted_Wallet(walletPath, password)
	if err != nil {
		if err.Error() == "Invalid Password" {
			return err
		}

		if bkCopied {
			return err
		}

		// maybe the wallet file is corrupt or does not exists
		// we will try to use backup file and copy as last resort
		walletBkPath := filepath.Join(walletsDir, addr, "wallet.db.bak")
		bkFile, err := os.Open(walletBkPath)
		if err != nil {
			return err
		}

		walletFile, err := os.Create(walletPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(walletFile, bkFile)
		if err != nil {
			return err
		}

		bkCopied = true
		goto open_wallet
	}

	memory.SetNetwork(globals.IsMainnet())

	CloseOpenedWallet()
	dbPath := filepath.Join(walletsDir, addr, "data.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	err = schema_version.Init(db)
	if err != nil {
		return err
	}

	err = initTableOutgoingTxs(db)
	if err != nil {
		return err
	}

	err = initTableTokens(db)
	if err != nil {
		return err
	}

	err = initTableContacts(db)
	if err != nil {
		return err
	}

	account := memory.GetAccount()
	// fix: looks like EntriesNative is not instantiated on startup but only in InsertReplace func???
	if account.EntriesNative == nil {
		account.EntriesNative = make(map[crypto.Hash][]rpc.Entry)
	}

	wallet := &Wallet{
		Info:       walletInfo,
		Memory:     memory,
		DB:         db,
		FolderPath: folderPath,
	}

	err = wallet.LoadSettings()
	if err != nil {
		return err
	}

	OpenedWallet = wallet
	return nil
}

func (w *Wallet) OpenXSWD(appHandler func(appData *xswd.ApplicationData) bool, reqHandler func(appData *xswd.ApplicationData, req *jrpc2.Request) xswd.Permission) {
	// create the XSWD server and start listening to incoming calls for authorization
	// XSWD is a secure communication protocol that offers easy interaction between the user wallet and a dApp
	// it was create by Slixe
	w.ServerXSWD = xswd.NewXSWDServer(w.Memory, appHandler, reqHandler)
}

func (w *Wallet) CloseXSWD() {
	if w.ServerXSWD != nil {
		w.ServerXSWD.Stop()
		w.ServerXSWD = nil
	}
}

func (w *Wallet) Delete() error {
	err := w.DB.Close()
	if err != nil {
		return err
	}

	return DeleteWallet(w.Info.Addr)
}

func (w *Wallet) RefreshInfo() error {
	addr := w.Memory.GetAddress()
	walletInfo, err := app_db.GetWalletInfo(addr.String())
	if err != nil {
		return err
	}

	w.Info = walletInfo
	return nil
}

func DeleteWallet(addr string) error {
	walletsDir := settings.WalletsDir
	path := filepath.Join(walletsDir, addr)
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	err = app_db.DelWalletInfo(addr)
	if err != nil {
		return err
	}

	return nil
}

func CreateWalletFromPath(name string, password string, path string) error {
	wallet, err := walletapi.Open_Encrypted_Wallet(path, password)
	if err != nil {
		return err
	}

	return createWallet(wallet.Wallet_Memory, name)
}

func CreateWalletFromData(name string, password string, data []byte) error {
	walletMemory, err := walletapi.Open_Encrypted_Wallet_Memory(password, data)
	if err != nil {
		return err
	}

	return createWallet(walletMemory, name)
}

func CreateWalletFromSeed(name string, password string, seed string) error {
	wallet, err := walletapi.Create_Encrypted_Wallet_From_Recovery_Words_Memory(password, seed)
	if err != nil {
		return err
	}

	return createWallet(wallet, name)
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

	return createWallet(wallet, name)
}

func CreateRandomWallet(name string, password string) error {
	wallet, err := walletapi.Create_Encrypted_Wallet_Random_Memory(password)
	if err != nil {
		return err
	}

	return createWallet(wallet, name)
}

func (w *Wallet) Rename(newName string) error {
	err := app_db.UpdateWalletInfo(w.Info)
	if err != nil {
		return err
	}

	w.Info.Name = newName
	return nil
}

func (w *Wallet) ResetBalanceResult(scId string) {
	account := w.Memory.GetAccount()
	hash := crypto.HashHexToHash(scId)

	account.Lock()
	idxToDelete := -1
	for i, balanceResult := range account.Balance_Result {
		if balanceResult.SCID == hash {
			idxToDelete = i
			break
		}
	}
	if idxToDelete > -1 {
		if idxToDelete == len(account.Balance_Result)-1 {
			account.Balance_Result = account.Balance_Result[:idxToDelete]
		} else {
			account.Balance_Result = append(account.Balance_Result[:idxToDelete], account.Balance_Result[idxToDelete+1:]...)
		}
	}
	account.Unlock()
}

func (w *Wallet) ChangePassword(password string, newPassword string) error {
	seed := w.Memory.GetAccount().Keys.Secret
	newMemory, err := walletapi.Create_Encrypted_Wallet_Memory(newPassword, seed)
	if err != nil {
		return err
	}

	return saveWalletData(newMemory)
}

func StoreRegistrationTx(addr string, tx *transaction.Transaction) error {
	txHex := hex.EncodeToString(tx.Serialize())
	walletInfo, err := app_db.GetWalletInfo(addr)
	if err != nil {
		return err
	}

	walletInfo.RegistrationTxHex = txHex
	return app_db.UpdateWalletInfo(walletInfo)
}

// This function is unused but we can keep it.
// I opted to insert tokens from DEX page instead of this hardcoded version.
// Check askToCreateFolderTokens() in dex_pairs.go
func (w *Wallet) InsertDexTokensFolder() error {
	dexFolder := TokenFolder{
		ParentId: sql.NullInt64{},
		Name:     "DEX Tokens",
	}

	id, err := w.InsertFolderToken(dexFolder)
	if err != nil {
		return err
	}

	folderId := sql.NullInt64{Int64: id, Valid: true}

	// Image for tokens will be loaded automatically when fetching dex data
	tokens := []Token{
		// DUSDT
		Token{
			SCID:           "f93b8d7fbbbf4e8f8a1e91b7ce21ac5d2b6aecc4de88cde8e929bce5f1746fbd",
			Name:           "Dero wrapped Tether USD",
			Decimals:       6,
			StandardType:   sc.DEX_SC_TYPE,
			Symbol:         sql.NullString{String: "DUSDT", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			IsFavorite:     true, //sql.NullBool{Bool: true, Valid: true},
			ImageUrl:       sql.NullString{String: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xdAC17F958D2ee523a2206206994597C13D831ec7/logo.png", Valid: true},
		},
		// DUSDC
		Token{
			SCID:           "bc161c4f65285d5d927e9749fddbd127859748be7e161099f2f6785edc70b3dc",
			Name:           "Dero wrapped USD Coin",
			Decimals:       6,
			StandardType:   sc.DEX_SC_TYPE,
			Symbol:         sql.NullString{String: "DUSDC", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			IsFavorite:     true, //sql.NullBool{Bool: true, Valid: true},
			ImageUrl:       sql.NullString{String: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48/logo.png", Valid: true},
		},
		// DWBTC
		Token{
			SCID:           "b0bb9c1c75fc0e84dd92ce03f0619d1b61737981f0bb796911ea31529a76358c",
			Name:           "Dero wrapped Wrapped BTC",
			Decimals:       7,
			StandardType:   sc.DEX_SC_TYPE,
			Symbol:         sql.NullString{String: "DWBTC", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			IsFavorite:     true, //sql.NullBool{Bool: true, Valid: true},
			ImageUrl:       sql.NullString{String: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599/logo.png", Valid: true},
		},
		// DWETH
		Token{
			SCID:           "fb855d8edd1d95ea94e9544224019c3fe4e636086f7266808879d6134ee2b8f1",
			Name:           "Dero wrapped Wrapped Ether",
			Decimals:       7,
			StandardType:   sc.DEX_SC_TYPE,
			Symbol:         sql.NullString{String: "DWETH", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			IsFavorite:     true, //sql.NullBool{Bool: true, Valid: true},
			ImageUrl:       sql.NullString{String: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png", Valid: true},
		},
		// DST
		Token{
			SCID:           "d74d1bb9968e3947a9bd40c5a9bdf598135f6b07a93bc98ded1fefa6ddd36bf5",
			Name:           "Dero Seals Token",
			Decimals:       5,
			MaxSupply:      sql.NullInt64{Int64: 2_800_000_00000, Valid: true}, // 2,800,000.00000
			StandardType:   sc.G45_FAT_TYPE,
			Symbol:         sql.NullString{String: "DST", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			IsFavorite:     true, //sql.NullBool{Bool: true, Valid: true},
			ImageUrl:       sql.NullString{String: "ipfs://QmboGpusU71C9zBPNjxskrXfY7GX1uoPo83MJ7NiJU2RUP/dero_seals_token.jpg", Valid: true},
		},
		// COCO
		Token{
			SCID:           "a9a977297ed6ed087861bfa117e6fbbd603c2051b0cc1b0d704bc764011aabb6",
			Name:           "Coconuts",
			Decimals:       0,
			StandardType:   sc.UNKNOWN_TYPE,
			Symbol:         sql.NullString{String: "COCO", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			ImageUrl:       sql.NullString{String: "", Valid: true},
		},
		// DLINK
		Token{
			SCID:           "ab8ee3627b212a0b3803c127f3de7c44465fac21ec30692cb7988b14059990bb",
			Name:           "Dero wrapped ChainLink Token",
			Decimals:       7,
			StandardType:   sc.DEX_SC_TYPE,
			Symbol:         sql.NullString{String: "DLINK", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			ImageUrl:       sql.NullString{String: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0x514910771AF9Ca656af840dff83E8264EcF986CA/logo.png", Valid: true},
		},
		// DgOHM
		Token{
			SCID:           "92136ec02ca1e0db8e1767f7d5d221c7951263790fe4ee6616c4dd6c011e65ba",
			Name:           "Dero wrapped Governance OHM",
			Decimals:       7,
			StandardType:   sc.DEX_SC_TYPE,
			Symbol:         sql.NullString{String: "DgOHM", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			ImageUrl:       sql.NullString{String: "https://raw.githubusercontent.com/OlympusDAO/olympus-frontend/develop/src/assets/tokens/token_OHM.svg", Valid: true},
		},
		// DFRAX
		Token{
			SCID:           "f42fd725bc3659a7e6502ce416363afea0951e7f21af4f8f71b42090206e29d4",
			Name:           "Dero wrapped Frax",
			Decimals:       7,
			StandardType:   sc.DEX_SC_TYPE,
			Symbol:         sql.NullString{String: "DFRAX", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			ImageUrl:       sql.NullString{String: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0x853d955aCEf822Db058eb8505911ED77F175b99e/logo.png", Valid: true},
		},
		// DDAI
		Token{
			SCID:           "93707e89ba07f9aafc862ae07df1bfa70f488d5157d37439b85498fb79b6d1e6",
			Name:           "Dero wrapped Dai Stablecoin",
			Decimals:       7,
			StandardType:   sc.DEX_SC_TYPE,
			Symbol:         sql.NullString{String: "DDAI", Valid: true},
			FolderId:       folderId,
			AddedTimestamp: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
			ImageUrl:       sql.NullString{String: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0x6B175474E89094C44Da98b954EedeAC495271d0F/logo.png", Valid: true},
		},
	}

	for _, token := range tokens {
		err = w.InsertToken(token)
		if err != nil {
			break
		}
	}

	if err != nil {
		// remove folder if an error occured while adding tokens
		err2 := w.DelTokenFolder(folderId.Int64)
		if err2 != nil {
			return err2
		}

		return err
	}

	return nil
}

func createWallet(wallet *walletapi.Wallet_Memory, name string) error {
	wallet.SetNetwork(globals.IsMainnet())

	addr := wallet.GetAddress().String()
	walletInfo := app_db.WalletInfo{
		Addr:        addr,
		Name:        name,
		Timestamp:   time.Now().Unix(),
		OrderNumber: -1,
	}

	err := app_db.InsertWalletInfo(walletInfo)
	if err != nil {
		return err
	}

	err = saveWalletData(wallet)
	if err != nil {
		return err
	}

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

func (w *Wallet) GetRandomAddresses(scId crypto.Hash) ([]string, error) {
	var result rpc.GetRandomAddress_Result

	err := RPCCall("DERO.GetRandomAddress", rpc.GetRandomAddress_Params{
		SCID: scId,
	}, &result)
	if err != nil {
		return nil, err
	}

	return result.Address, nil
}

func (w *Wallet) GetRandomAddress(scId crypto.Hash) (string, error) {
	walletAddr := w.Memory.GetAddress().String()

	addresses, err := w.GetRandomAddresses(scId)
	if err != nil {
		return "", err
	}

	if len(addresses) < 10 {
		// sample size too small
		addresses, err = w.GetRandomAddresses(crypto.ZEROHASH)
		if err != nil {
			return "", err
		}
	}

	var addrList []string
	for _, addr := range addresses {
		if addr != walletAddr { // make sure this wallet addr is not part of the randomly selected addr
			addrList = append(addrList, addr)
		}
	}

	index := rand.Intn(len(addrList))
	return addrList[index], nil
}

func RPCCall(method string, params interface{}, result interface{}) error {
	rpcClient := walletapi.GetRPCClient()
	if rpcClient.RPC == nil {
		return fmt.Errorf("node client is not connected")
	}

	return rpcClient.RPC.CallResult(context.Background(), method, params, result)
}
