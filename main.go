package main

import (
	"log"
	"net/http"

	"waste-EcoTech/handlers"
)

func main() {
	file := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", file))

	http.HandleFunc("/", handlers.HomeHandler)
	// http.HandleFunc("/resident-register", handlers.ResidentRegisterHandler)
	// http.HandleFunc("/resident-login", handlers.ResidentLoginHandler)
	http.HandleFunc("/collection", handlers.CollectionProcessing)
	http.HandleFunc("/view-data", handlers.ViewRecyclesHandler)
	http.HandleFunc("/segregation", handlers.Segregation)
	http.HandleFunc("/pro-dashboard", handlers.AddManufacturerItems)

	log.Println("server running at : http://localhost:1234")
	http.ListenAndServe(":1234", nil)

}
