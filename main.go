package main

import (
	"fmt"
	"net/http"

	o "blockchain/project"
)

func main() {
	http.HandleFunc("/blocks", AdBlock)

	fmt.Println("Server listening on port http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func AdBlock(w http.ResponseWriter, r *http.Request) {
	// Create a new blockchain
	bc :=&o.Blockchain{}
	// Create the genesis block
	genesisBlock := o.NewGenesisBlock()
	bc.Blocks = append(bc.Blocks, genesisBlock)

	// Add additional blocks (replace with your desired data)
	bc.AddBlock("Send 1 BTC to Alice")
	bc.AddBlock("Send 2 ETH to Bob")

	// Print the blockchain
	bc.PrintChain()
}
