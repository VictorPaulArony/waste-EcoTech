package main

import (
	"log"
	"net/http"

	"waste_Eco_Track/handlers"
)

func main() {
	file := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", file))

	http.HandleFunc("/", handlers.HomeHandler)
	// http.HandleFunc("/resident-register", handlers.ResidentRegisterHandler)
	// http.HandleFunc("/resident-login", handlers.ResidentLoginHandler)
	// http.HandleFunc("/staff-register", handlers.StaffRegistrationHandler)
	// http.HandleFunc("/staff-login", handlers.StaffLoginHandler)
	// http.HandleFunc("/company-dashboard", handlers.StaffDshboardHandler)
	//http.HandleFunc("/resident-dash", handlers.)

	log.Println("server running at : http://localhost:1234")
	http.ListenAndServe(":1234", nil)


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
