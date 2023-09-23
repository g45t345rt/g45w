package wallet_manager

import (
	"database/sql"
	"math"
	"sort"
	"sync"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
)

type Entry struct {
	rpc.Entry
	SCID crypto.Hash
}

type SCCallParams struct {
	SCID       sql.NullString
	Entrypoint sql.NullString
}

type GetEntriesParams struct {
	In                       sql.NullBool
	Out                      sql.NullBool
	Coinbase                 sql.NullBool
	Sender                   sql.NullString
	Receiver                 sql.NullString
	BurnGreaterOrEqualThan   sql.NullInt64
	AmountGreaterOrEqualThan sql.NullInt64
	TXID                     sql.NullString
	BlockHash                sql.NullString
	SC_CALL                  *SCCallParams
	Offset                   sql.NullInt64
	Limit                    sql.NullInt64
}

func filterEntries(allEntries []rpc.Entry, start, end int, entrySCID crypto.Hash, params GetEntriesParams, entryChan chan<- Entry, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := start; i < end; i++ {
		e := allEntries[i]
		add := true

		if params.Coinbase.Valid {
			add = e.Coinbase == params.Coinbase.Bool
		}

		if params.In.Valid {
			add = e.Incoming == params.In.Bool
		}

		if params.Out.Valid {
			add = (!e.Incoming && !e.Coinbase) == params.Out.Bool
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

		if params.SC_CALL != nil {
			add = false
			for _, arg := range e.SCDATA {
				if params.SC_CALL.SCID.Valid && arg.Name == "SC_ID" {
					scId, ok := arg.Value.(string)
					if ok && scId == params.SC_CALL.SCID.String {
						add = true
					}
				}

				if params.SC_CALL.Entrypoint.Valid && arg.Name == "entrypoint" {
					entrypoint, ok := arg.Value.(string)
					if ok && entrypoint == params.SC_CALL.Entrypoint.String {
						add = true
						break
					}
				}
			}
		}

		if add {
			entry := Entry{Entry: e, SCID: entrySCID}
			entryChan <- entry
		}
	}
}

func (w *Wallet) GetEntries(SCID *crypto.Hash, params GetEntriesParams) []Entry {
	w.Memory.Lock()
	defer w.Memory.Unlock()

	account := w.Memory.GetAccount()
	var wg sync.WaitGroup
	entryChan := make(chan Entry)

	var filteredEntries []Entry
	done := make(chan bool)
	go func() {
		for e := range entryChan {
			filteredEntries = append(filteredEntries, e)
		}

		done <- true
	}()

	workSize := 100
	for entrySCID, entries := range account.EntriesNative {
		if SCID != nil && entrySCID != *SCID {
			continue
		}

		workers := int(math.Max(1, float64(len(entries)/workSize)))
		for i := 0; i < workers; i++ {
			start := i * workSize
			end := (i + 1) * workSize
			if i == workers-1 {
				end = len(entries)
			}

			wg.Add(1)
			go filterEntries(entries, start, end, entrySCID, params, entryChan, &wg)
		}
	}

	wg.Wait()
	close(entryChan)
	<-done

	sort.Slice(filteredEntries, func(a, b int) bool {
		return filteredEntries[a].Time.Unix() > filteredEntries[b].Time.Unix()
	})

	if params.Offset.Valid {
		offset := params.Offset.Int64
		if len(filteredEntries) > int(offset) {
			filteredEntries = filteredEntries[offset:]
		}
	}

	if params.Limit.Valid {
		limit := params.Limit.Int64
		if len(filteredEntries) > int(limit) {
			filteredEntries = filteredEntries[:limit]
		}
	}

	return filteredEntries
}
