package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/lookup_table"
)

// go run ./lookup_table/create

func main() {
	table := walletapi.Initialize_LookupTable(1, 1<<21)

	fmt.Println("Creating lookup table...")
	data, err := lookup_table.Serialize(table)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./lookup_table/lookup_table_2", data, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
