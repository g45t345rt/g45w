package lookup_table

import (
	"bytes"
	_ "embed"
	"encoding/gob"

	"github.com/deroproject/derohe/walletapi"
)

//go:embed lookup_table
var LOOKUP_TABLE_BYTES []byte

func Load() error {
	var lookupTable walletapi.LookupTable
	err := Deserialize(LOOKUP_TABLE_BYTES, &lookupTable)
	if err != nil {
		return err
	}

	walletapi.Balance_lookup_table = &lookupTable

	return nil
}

func Serialize(table *walletapi.LookupTable) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(table)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Deserialize(buf []byte, table *walletapi.LookupTable) error {
	buffer := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(buffer)

	err := dec.Decode(&table)
	if err != nil {
		return err
	}

	return nil
}
