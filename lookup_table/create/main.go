package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deroproject/derohe/walletapi"
)

// go run ./lookup_table/create

func main() {
	walletapi.Initialize_LookupTable(1, 1<<21)

	fmt.Println("Creating lookup table...")
	data, err := walletapi.Balance_lookup_table.Serialize()
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./lookup_table/lookup_table", data, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
