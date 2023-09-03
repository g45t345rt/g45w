package lookup_table

import (
	_ "embed"

	"github.com/deroproject/derohe/walletapi"
)

//go:embed lookup_table
var LOOKUP_TABLE []byte

func Load() error {
	var lookupTable walletapi.LookupTable
	err := lookupTable.Deserialize(LOOKUP_TABLE)
	if err != nil {
		return err
	}

	walletapi.Balance_lookup_table = &lookupTable

	return nil
}
