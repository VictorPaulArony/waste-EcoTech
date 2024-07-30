package database

import (
	"encoding/json"
	"os"
)

type Recycle struct {
	ID        int    `json:"id"`
	Producer  string `json:"producer"`
	Type      string `json:"type"`
	Code      string `json:"code"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

var (
	recycleFile = "recycle.json"
	recycles    []Recycle
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

// func generateBarcode() {
// 	res := Math.floor(Math.random()*1000000).toString().padStart(6, '0')
// 	fmt.Sprintf("%v", res)

// 	return
// }
