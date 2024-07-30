package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
	"waste-EcoTech/database"
	// "waste-EcoTrack/blockchain"
	// "waste-EcoTrack/database"
)

var (
	muSync   sync.Mutex
	recycles []database.Recycle
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/home" {
		temp := template.Must(template.ParseFiles("templates/manufacturer.html"))
		e := temp.Execute(w, nil)
		if e != nil {
			log.Fatalln("Internal server error")
			fmt.Fprint(w, "oops something went wrong")
		}
		return
	}
}
func AddManufacturerItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/manufacturer.html"))
		temp.Execute(w, nil)
		return
	} else if r.Method == http.MethodPost {
		muSync.Lock()
		defer muSync.Unlock()
		r.ParseForm()
		request := database.Recycle{
			ID:        len(recycles) + 1,
			Producer:  r.FormValue("batch-name"),
			Type:      r.FormValue("bottle-count"),
			Code:      r.FormValue("manufacturer-name"),
			CreatedAt: time.Now().String(),
		}
		recycles = append(recycles, request)
		println("data")
		if err := database.SaveRecycle(recycles); err != nil {
			http.Error(w, "Error Saving the Request", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/pro-dashboard", http.StatusSeeOther)
	}
}
func ViewRecyclesHandler(w http.ResponseWriter, r *http.Request) {
	muSync.Lock()
	defer muSync.Unlock()

	temp := template.Must(template.ParseFiles("templates/view-recycle.html"))
	if err := temp.Execute(w, recycles); err != nil {
		log.Fatalln("Internal server error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
func CollectionProcessing(w http.ResponseWriter, r *http.Request) {

}
