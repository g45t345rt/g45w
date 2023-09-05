package wallet_manager

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/deroproject/derohe/config"
	"github.com/deroproject/derohe/cryptography/bn256"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/transaction"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_db/schema_version"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/utils"

	"database/sql"

	_ "modernc.org/sqlite" // use CGO free port
)

type Wallet struct {
	Info   *WalletInfo
	Memory *walletapi.Wallet_Disk
	DB     *sql.DB
}

type WalletInfo struct {
	Name              string `json:"name"`
	Addr              string `json:"addr"`
	RegistrationTxHex string `json:"registration_tx_hex"`
	//Data      []byte `json:"data"`
	Timestamp int64 `json:"timestamp"`
	ListOrder int   `json:"order"` // save item list ordering
}

var Wallets map[string]*WalletInfo
var WalletsErr map[string]error
var OpenedWallet *Wallet

func Load() error {
	walletsDir := settings.WalletsDir
	Wallets = make(map[string]*WalletInfo)
	WalletsErr = make(map[string]error)

	err := os.MkdirAll(walletsDir, os.ModePerm)
	if err != nil {
		return err
	}

	filepath.Walk(walletsDir, func(path string, info fs.FileInfo, fileErr error) error {
		if walletsDir == path {
			return nil
		}

		if info.IsDir() {
			addr := info.Name()

			_, err := globals.ParseValidateAddress(addr)
			if err != nil {
				return nil
			}

			fileInfo := filepath.Join(walletsDir, addr, "info.json")
			data, err := os.ReadFile(fileInfo)
			if err != nil {
				WalletsErr[addr] = err
				return nil
			}

			walletInfo := &WalletInfo{}
			err = json.Unmarshal(data, walletInfo)
			if err != nil {
				WalletsErr[addr] = err
				return nil
			}

			Wallets[addr] = walletInfo
		}

		return nil
	})

	return nil
}

func CloseOpenedWallet() {
	if OpenedWallet != nil {
		wallet := OpenedWallet
		go func() {
			close(wallet.Memory.Quit) // make sure to close goroutines when wallet is in online mode
			wallet.Memory.Close_Encrypted_Wallet()
		}()
		wallet.DB.Close()
		OpenedWallet = nil
	}
}

func GetWallet(addr string) (*WalletInfo, error) {
	for _, walletInfo := range Wallets {
		if walletInfo.Addr == addr {
			return walletInfo, nil
		}
	}

	return nil, fmt.Errorf("wallet [%s] not found", addr)
}

func OpenWallet(addr string, password string) error {
	walletInfo, err := GetWallet(addr)
	if err != nil {
		return err
	}

	walletsDir := settings.WalletsDir
	walletPath := filepath.Join(walletsDir, addr, "wallet.db")

	bkCopied := false
open_wallet:
	memory, err := walletapi.Open_Encrypted_Wallet(walletPath, password)
	if err != nil {
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

	err = initDatabaseOutgoingTxs(db)
	if err != nil {
		return err
	}

	err = initDatabaseTokens(db)
	if err != nil {
		return err
	}

	err = initDatabaseContacts(db)
	if err != nil {
		return err
	}

	account := memory.GetAccount()
	// fix: looks like EntriesNative is not instantiated on startup but only in InsertReplace func???
	if account.EntriesNative == nil {
		account.EntriesNative = make(map[crypto.Hash][]rpc.Entry)
	}

	wallet := &Wallet{
		Info:   walletInfo,
		Memory: memory,
		DB:     db,
	}

	OpenedWallet = wallet
	return nil
}

func DeleteWallet(addr string) error {
	walletsDir := settings.WalletsDir
	path := filepath.Join(walletsDir, addr)
	err := os.RemoveAll(path)
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

func CreateWalletFromData(name string, password string, data []byte) error {
	walletMemory, err := walletapi.Open_Encrypted_Wallet_Memory(password, data)
	if err != nil {
		return err
	}

	return saveWallet(walletMemory, name)
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

func (w *Wallet) Rename(newName string) error {
	w.Info.Name = newName
	return saveWalletInfo(w.Info.Addr, w.Info)
}

func (w *Wallet) OrderList(newOrder int) error {
	w.Info.ListOrder = newOrder
	return saveWalletInfo(w.Info.Addr, w.Info)
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

type RingMembers struct {
	RingsBalances [][][]byte
	Rings         [][]*bn256.G1
	RingsAddrs    []map[string]bool
	MaxBits       int
}

func (w *Wallet) BuildRingMembers(transfers []rpc.Transfer, ringsize uint64) (ringMembers RingMembers, err error) {
	walletAddr := w.Memory.GetAddress().String()
	account := w.Memory.GetAccount()

	for _, transfer := range transfers {
		var ring []*bn256.G1
		var ringBalances [][]byte
		ringAddrs := make(map[string]bool, 0)

		bitsNeeded := make([]int, ringsize)

		var selfEncryptedBalance *crypto.ElGamal
		bitsNeeded[0], _, _, selfEncryptedBalance, err = w.Memory.GetEncryptedBalanceAtTopoHeight(transfer.SCID, -1, walletAddr)
		if err != nil {
			return
		}

		ringBalances = append(ringBalances, selfEncryptedBalance.Serialize())
		ring = append(ring, account.Keys.Public.G1())
		ringAddrs[walletAddr] = true

		if transfer.Destination == walletAddr {
			err = fmt.Errorf("can't send to self")
			return
		}

		var destAddr *rpc.Address
		destAddr, err = rpc.NewAddress(transfer.Destination)
		if err != nil {
			return
		}

		var destEncryptedBalance *crypto.ElGamal
		bitsNeeded[1], _, _, destEncryptedBalance, err = w.Memory.GetEncryptedBalanceAtTopoHeight(transfer.SCID, -1, destAddr.String())
		if err != nil {
			return
		}

		ringBalances = append(ringBalances, destEncryptedBalance.Serialize())
		ring = append(ring, destAddr.PublicKey.G1())
		ringAddrs[transfer.Destination] = true

	loadRingMembers:
		var addrList []string
		addrList, err = w.GetRandomAddresses(transfer.SCID)
		if err != nil {
			return
		}

		if len(addrList) < len(ringBalances) {
			// get addr list from base asset if we don't have enough for current asset
			addrList, err = w.GetRandomAddresses(crypto.ZEROHASH)
			if err != nil {
				return
			}
		}

		for _, addr := range addrList {
			_, addrExists := ringAddrs[addr] // avoid adding same addr including wallet addr and destination addr
			if len(ringBalances) == int(ringsize) ||
				addrExists {
				continue
			}

			var memberAddr *rpc.Address
			memberAddr, err = rpc.NewAddress(addr)
			if err != nil {
				return
			}

			var memberEncryptedBalance *crypto.ElGamal
			bitsNeeded[len(ringBalances)], _, _, memberEncryptedBalance, err = w.Memory.GetEncryptedBalanceAtTopoHeight(transfer.SCID, -1, addr)
			if err != nil {
				return
			}

			ringBalances = append(ringBalances, memberEncryptedBalance.Serialize())
			ring = append(ring, memberAddr.PublicKey.G1())
			ringAddrs[memberAddr.String()] = true
		}

		if len(ringBalances) < int(ringsize) {
			goto loadRingMembers
		}

		ringMembers.RingsBalances = append(ringMembers.RingsBalances, ringBalances)
		ringMembers.Rings = append(ringMembers.Rings, ring)
		ringMembers.RingsAddrs = append(ringMembers.RingsAddrs, ringAddrs)

		for i := range bitsNeeded {
			if ringMembers.MaxBits < bitsNeeded[i] {
				ringMembers.MaxBits = bitsNeeded[i]
			}
		}
	}

	ringMembers.MaxBits += 6 // extra 6 bits for unknown reasons?

	return
}

func (w *Wallet) GetGasEstimate(transfers []rpc.Transfer, ringsize uint64, scArgs rpc.Arguments) (uint64, error) {
	signer := w.Memory.GetAddress().String()

	var result rpc.GasEstimate_Result
	err := walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetGasEstimate", rpc.GasEstimate_Params{
		Transfers: transfers,
		SC_RPC:    scArgs,
		Ringsize:  ringsize,
		Signer:    signer,
	}, &result)
	if err != nil {
		return 0, err
	}

	return result.GasStorage, nil
}

func (w *Wallet) BuildTransaction(transfers []rpc.Transfer, ringsize uint64, scArgs rpc.Arguments, dryRun bool) (tx *transaction.Transaction, txFees uint64, gasFees uint64, err error) {
	if len(scArgs) > 0 {
		// smart contract call to test if it can succeed and return gas fees for the dry run
		gasFees, err = w.GetGasEstimate(transfers, ringsize, scArgs)
		if err != nil {
			return
		}
	}

	// need at least one Dero transfers
	hasBase := false
	for _, t := range transfers {
		if t.SCID.IsZero() {
			hasBase = true
			break
		}
	}

	if !hasBase {
		var randomAddr string
		randomAddr, err = w.GetRandomAddress(crypto.ZEROHASH)
		if err != nil {
			return
		}

		transfers = append(transfers, rpc.Transfer{
			SCID:        crypto.ZEROHASH,
			Destination: randomAddr,
			Amount:      0,
		})
	}

	var ringMembers RingMembers
	ringMembers, err = w.BuildRingMembers(transfers, ringsize)
	if err != nil {
		return
	}

	topoHeight := w.Memory.Get_TopoHeight()

	var encryptedBalance rpc.GetEncryptedBalance_Result
	encryptedBalance, err = w.Memory.GetSelfEncryptedBalanceAtTopoHeight(crypto.ZEROHASH, topoHeight)
	if err != nil {
		return
	}

	height := uint64(encryptedBalance.Height)
	blockHash := encryptedBalance.BlockHash
	treeHash := encryptedBalance.Merkle_Balance_TreeHash

	var treeHashRaw []byte
	treeHashRaw, err = hex.DecodeString(treeHash)
	if err != nil {
		return
	}

	// get len of for Dero transfers - the fees are only applied to Dero asset statements and not other tokens
	deroTransfers := 0
	assetsAmount := make(map[crypto.Hash]uint64, 0)
	for _, transfer := range transfers {
		_, ok := assetsAmount[transfer.SCID]
		if !ok {
			assetsAmount[transfer.SCID] = 0
		}

		assetsAmount[transfer.SCID] += transfer.Amount + transfer.Burn
		if transfer.SCID.IsZero() {
			deroTransfers++
		}
	}

	for asset, amount := range assetsAmount {
		balance, _ := w.Memory.Get_Balance_scid(asset)
		if amount > balance {
			if asset.IsZero() {
				err = fmt.Errorf("you don't have enough Dero")
			} else {
				err = fmt.Errorf("you don't have enough asset funds of [%s]", utils.ReduceTxId(asset.String()))
			}

			return
		}
	}

	// build a dry transaction to get transaction size and calculate fees
	// set fees to 1 to avoid automatic fees in statement
	// fee value is store in tx but its too small for making any adjustments
	tx = w.Memory.BuildTransaction(
		transfers, ringMembers.RingsBalances, ringMembers.Rings, blockHash,
		height, scArgs, treeHashRaw, ringMembers.MaxBits, 1)
	if tx == nil {
		err = fmt.Errorf("can't build transaction")
		return
	}

	txSize := uint64(len(tx.Serialize()))
	txFees = w.CalculateTxFees(txSize)
	totalFees := txFees + gasFees

	if dryRun {
		return
	}

	// set fees in BuildTransaction applies it to all Dero transfers -_-
	// split the fees amongst all Dero transfers
	feesPerTransfer := uint64(math.Ceil(float64(totalFees) / float64(deroTransfers)))

	tx = w.Memory.BuildTransaction(
		transfers, ringMembers.RingsBalances, ringMembers.Rings, blockHash,
		height, scArgs, treeHashRaw, ringMembers.MaxBits, feesPerTransfer)
	if tx == nil {
		err = fmt.Errorf("can't build transaction")
		return
	}

	return
}

func (w *Wallet) CalculateTxFees(sizeInBytes uint64) (fees uint64) {
	size := sizeInBytes / 1024

	if size%1024 != 0 {
		size += 1 // add full kb for any rest
	}

	return size*config.FEE_PER_KB + 1 // add plus one Deri because mempool fee check is using > instead of >=
}

func StoreRegistrationTx(addr string, tx *transaction.Transaction) error {
	txHex := hex.EncodeToString(tx.Serialize())
	walletInfo, err := GetWallet(addr)
	if err != nil {
		return err
	}

	walletInfo.RegistrationTxHex = txHex
	return saveWalletInfo(addr, walletInfo)
}

func saveWallet(wallet *walletapi.Wallet_Memory, name string) error {
	wallet.SetNetwork(globals.IsMainnet())

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

	infoPath := filepath.Join(walletsDir, addr, "info.json")
	err = os.WriteFile(infoPath, data, fs.ModePerm)
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

func (w *Wallet) GetRandomAddresses(scId crypto.Hash) ([]string, error) {
	var result rpc.GetRandomAddress_Result

	err := walletapi.RPC_Client.Call("DERO.GetRandomAddress", rpc.GetRandomAddress_Params{
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
