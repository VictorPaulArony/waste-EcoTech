package database

import (
	"encoding/json"
	"os"
)

type Recycle struct {
	ID        int    `json:"id"`
	Producer  string `json:"batch-name"`
	Type      string `json:"bottle-count"`
	Code      string `json:"manufacturer-name"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Hash      string `json:"hash"`
}

type Collection struct {
	ID        int    `json:"id"`
	Producer  string `json:"producer"`
	Type      string `json:"type"`
	Code      string `json:"code"`
	TimeStamp string `json:"timestamp"`
	Status    string `json:"status"`
}

var (
	recycleFile    = "recycle.json"
	collectionFile = "collections.json"
	recycles       []Recycle
	collections    []Collection
)

func SaveRecycle(recycles []Recycle) error {
	data, err := json.MarshalIndent(recycles, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(recycleFile, data, 0o644)
}

func LoadRecycle() error {
	file, err := os.ReadFile(recycleFile)
	if err != nil {
		if os.IsNotExist(err) {
			recycles = []Recycle{}
		}
		return err
	}
	return json.Unmarshal(file, &recycles)
}

func SaveCollection(collections []Collection) error {
	data, err := json.MarshalIndent(collections, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(collectionFile, data, 0o644)
}

func LoadCollection() error {
	file, err := os.ReadFile(collectionFile)
	if err != nil {
		if os.IsNotExist(err) {
			collections = []Collection{}
		}
		return err
	}
	return json.Unmarshal(file, &collections)
}
