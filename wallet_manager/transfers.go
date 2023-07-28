package wallet_manager

import (
	"database/sql"
	"runtime"
	"sort"
	"sync"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
)

type GetTransfersParams struct {
	In                       sql.NullBool
	Out                      sql.NullBool
	Coinbase                 sql.NullBool
	Sender                   sql.NullString
	Receiver                 sql.NullString
	BurnGreaterOrEqualThan   sql.NullInt64
	AmountGreaterOrEqualThan sql.NullInt64
	TXID                     sql.NullString
	BlockHash                sql.NullString
	Offset                   sql.NullInt64
	Limit                    sql.NullInt64
}

func filterEntries(allEntries []rpc.Entry, params GetTransfersParams, start, end int, entryChan chan<- rpc.Entry, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := start; i < end; i++ {
		e := allEntries[i]
		add := true

		if params.Coinbase.Valid {
			add = e.Coinbase == params.Coinbase.Bool
		}

		if params.In.Valid {
			add = (e.Incoming && !e.Coinbase) == params.In.Bool
		}

		if params.Out.Valid {
			add = !(e.Incoming || e.Coinbase) == params.Out.Bool
		}

		if params.Sender.Valid {
			add = e.Sender == params.Sender.String
		}

		if params.Receiver.Valid {
			add = e.Destination == params.Receiver.String
		}

		if params.AmountGreaterOrEqualThan.Valid {
			add = e.Amount >= uint64(params.AmountGreaterOrEqualThan.Int64)
		}

		if params.BurnGreaterOrEqualThan.Valid {
			add = e.Burn >= uint64(params.BurnGreaterOrEqualThan.Int64)
		}

		if params.TXID.Valid {
			add = e.TXID == params.TXID.String
		}

		if params.BlockHash.Valid {
			add = e.BlockHash == params.BlockHash.String
		}

		if add {
			entryChan <- e
		}
	}
}

func (w *Wallet) GetTransfers(scId string, params GetTransfersParams) []rpc.Entry {
	w.Memory.Lock()
	defer w.Memory.Unlock()

	account := w.Memory.GetAccount()
	allEntries := account.EntriesNative[crypto.HashHexToHash(scId)]
	totalEntries := len(allEntries)
	if allEntries == nil || totalEntries < 1 {
		return allEntries
	}

	workers := runtime.NumCPU()
	var wg sync.WaitGroup
	entryChan := make(chan rpc.Entry)
	chunkSize := totalEntries / workers
	for i := 0; i < workers; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if i == workers-1 {
			end = totalEntries
		}

		wg.Add(1)
		go filterEntries(allEntries, params, start, end, entryChan, &wg)
	}

	var entries []rpc.Entry
	go func() {
		for e := range entryChan {
			entries = append(entries, e)
		}
	}()

	wg.Wait()
	close(entryChan)

	sort.Slice(entries, func(a, b int) bool {
		return entries[a].Time.Unix() > entries[b].Time.Unix()
	})

	if params.Offset.Valid {
		offset := params.Offset.Int64
		if len(entries) > int(offset) {
			entries = entries[offset:]
		}
	}

	if params.Limit.Valid {
		limit := params.Limit.Int64
		if len(entries) > int(limit) {
			entries = entries[:limit]
		}
	}

	return entries
}
