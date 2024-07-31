package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"waste-EcoTech/blockchain"
	"waste-EcoTech/database"
)

var (
	muSync      sync.Mutex
	recycles    []database.Recycle
	collections []database.Collection
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/home" {
		temp := template.Must(template.ParseFiles("templates/manufacturer.html"))
		if err := temp.Execute(w, nil); err != nil {
			log.Fatalln("Internal server error")
			fmt.Fprint(w, "Oops, something went wrong")
		}
		return
	}
}
func Segregation(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("templates/segregation.html"))
	if err := temp.Execute(w, nil); err != nil {
		log.Fatalln("Internal server error")
		fmt.Fprint(w, "Oops, something went wrong")
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
			Status:    "Pending",
		}
		recycles = append(recycles, request)

		// Initialize blockchain
		bc := blockchain.Blockchain{}
		if err := bc.LoadBlock(); err != nil {
			http.Error(w, fmt.Sprintf("Error loading blockchain: %v", err), http.StatusInternalServerError)
			return
		}

		// Add to blockchain
		bc.AddBlock(fmt.Sprintf("Recycle ID: %d, Code: %s", request.ID, request.Code))

		if err := bc.SaveBlock(); err != nil {
			http.Error(w, fmt.Sprintf("Error saving blockchain: %v", err), http.StatusInternalServerError)
			return
		}

		if err := database.SaveRecycle(recycles); err != nil {
			http.Error(w, fmt.Sprintf("Error saving recycle: %v", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
}
func AddSegregatedItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/segregation.html"))
		temp.Execute(w, nil)
		return
	} else if r.Method == http.MethodPost {
		muSync.Lock()
		defer muSync.Unlock()

		r.ParseForm()
		hashValue := r.FormValue("hash")

		// Validate the hash value
		valid := false
		for _, rec := range recycles {
			if rec.Hash == hashValue {
				valid = true
				break
			}
		}

		if !valid {
			http.Error(w, "Invalid hash value", http.StatusBadRequest)
			return
		}

		request := database.Collection{
			ID:        len(collections) + 1,
			Producer:  r.FormValue("producer"),
			Type:      r.FormValue("type"),
			Code:      r.FormValue("code"),
			TimeStamp: time.Now().String(),
			Status:    "Pending",
		}
		collections = append(collections, request)

		// Initialize blockchain
		bc := blockchain.Blockchain{}
		if err := bc.LoadBlockRecycle(); err != nil {
			http.Error(w, fmt.Sprintf("Error loading blockchain: %v", err), http.StatusInternalServerError)
			return
		}

		// Add to blockchain
		bc.AddrecBlock(fmt.Sprintf("Producer: %s, Type: %s", request.Producer, request.Type), request.Type)

		if err := bc.SaveBlock(); err != nil {
			http.Error(w, fmt.Sprintf("Error saving blockchain: %v", err), http.StatusInternalServerError)
			return
		}

		if err := database.SaveCollection(collections); err != nil {
			http.Error(w, fmt.Sprintf("Error saving collection: %v", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/segregation", http.StatusFound)
		return
	}
}

func ViewRecyclesHandler(w http.ResponseWriter, r *http.Request) {
	muSync.Lock()
	defer muSync.Unlock()

	temp := template.Must(template.ParseFiles("templates/view.html"))
	if err := temp.Execute(w, recycles); err != nil {
		log.Fatalln("Internal server error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}


func CollectionProcessing(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	muSync.Lock()
	defer muSync.Unlock()

	var recycle *database.Recycle
	for i := range recycles {
		if recycles[i].ID == id {
			recycle = &recycles[i]
			break
		}
	}
	if recycle == nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}
	if recycle.Status != "Pending" {
		http.Error(w, "Request already processed", http.StatusBadRequest)
		return
	}

	recycle.Status = "Completed"
	collection := database.Collection{
		ID:        len(collections) + 1,
		Producer:  recycle.Producer,
		Type:      recycle.Type,
		Code:      recycle.Code,
		TimeStamp: time.Now().String(),
		Status:    recycle.Status,
	}

	// Load blockchain, add new collection, and save
	bc := blockchain.Blockchain{}
	if err := bc.LoadBlock(); err != nil {
		http.Error(w, fmt.Sprintf("Error loading blockchain: %v", err), http.StatusInternalServerError)
		return
	}
	bc.AddBlock(collection.Code)
	if err := bc.SaveBlock(); err != nil {
		http.Error(w, fmt.Sprintf("Error saving blockchain: %v", err), http.StatusInternalServerError)
		return
	}
	collections = append(collections, collection)

	if err := database.SaveCollection(collections); err != nil {
		http.Error(w, "Error saving certificates", http.StatusInternalServerError)
		return
	}
}

func RequestCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/manufacturer.html"))
		temp.Execute(w, nil)
		return
	} else if r.Method == http.MethodPost {
		muSync.Lock()
		defer muSync.Unlock()

		r.ParseForm()

		collection := database.Collection{
			ID:        len(recycles) + 1,
			Producer:  r.FormValue("batch-name"),
			Type:      r.FormValue("bottle-count"),
			Code:      r.FormValue("manufacturer-name"),
			TimeStamp: time.Now().String(),
			Status:    "pending",
		}
		collections = append(collections, collection)

		if err := database.SaveCollection(collections); err != nil {
			http.Error(w, "Error saving the request", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/student-dashboard", http.StatusSeeOther)
	}
}
