package blockchain

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"
)

type Collection struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	TimeStamp string `json:"timestamp"`
	PrevHash  string `json:"prevhash"`
	Hash      string `json:"hash"`
}
type Recycle struct {
	ID        int    `json:"id"`
	Producer  string `json:"batch-name"`
	Type      string `json:"bottle-count"`
	CreatedAt string `json:"created_at"`
	PrevHash  string `json:"prevhash"`
	Hash      string `json:"hash"`
}

type Blockchain struct {
	sync.Mutex
	Collections []Collection `json:"collections"`
	Recycles    []Recycle    `json:"Recycle"`
}

var fileName = "blocks.json"
var recycleFile = "recycalBlocks.json"

func CreateHash(col Collection) string {
	res := strconv.Itoa(col.ID) + col.Code + col.TimeStamp + col.PrevHash + col.Hash
	hash := sha512.New()
	hash.Write([]byte(res))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}
func CreateRecycleHash(rec Recycle) string {
	res := strconv.Itoa(rec.ID) + rec.Producer + rec.Type + rec.CreatedAt + rec.PrevHash + rec.Hash
	hash := sha512.New()
	hash.Write([]byte(res))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateGenesis() Collection {
	genesis := Collection{0, "Genesis Collection", time.Now().String(), "", ""}
	genesis.Hash = CreateHash(genesis)
	return genesis
}
func GenerateGenesisRecycle() Recycle {
	genesis := Recycle{0, "genesis Producer", "", time.Now().String(), "", " "}
	genesis.Hash = CreateRecycleHash(genesis)
	return genesis
}

func (bc *Blockchain) AddBlock(code string) string {
	bc.Lock()
	defer bc.Unlock()

	prevBlock := bc.Collections[len(bc.Collections)-1]
	newCollection := Collection{
		ID:        prevBlock.ID + 1,
		Code:      code,
		TimeStamp: time.Now().String(),
		PrevHash:  prevBlock.Hash,
	}
	newCollection.Hash = CreateHash(newCollection)
	bc.Collections = append(bc.Collections, newCollection)
	return newCollection.Hash
}
func (bc *Blockchain) AddrecBlock(code string, typer string) string {
	bc.Lock()
	defer bc.Unlock()

	prevBlock := bc.Collections[len(bc.Collections)-1]
	newCollection := Recycle{
		ID:        prevBlock.ID + 1,
		Producer:  code,
		Type:      typer,
		CreatedAt: time.Now().String(),
		PrevHash:  prevBlock.Hash,
	}
	newCollection.Hash = CreateRecycleHash(newCollection)
	bc.Recycles = append(bc.Recycles, newCollection)
	return newCollection.Hash
}

// Function to save blockchain to the JSON file
func (bc *Blockchain) SaveBlock() error {
	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(fileName, data, 0o644); err != nil {
		return err
	}
	return nil
}

// Function to load the blockchain from the JSON file
func (bc *Blockchain) LoadBlock() error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			bc.Collections = []Collection{GenerateGenesis()}
			return bc.SaveBlock()
		}
		return err
	}
	err = json.Unmarshal(file, bc)
	if err != nil {
		return err
	}
	if len(bc.Collections) == 0 {
		bc.Collections = []Collection{GenerateGenesis()}
	}
	return nil
}

func (bc *Blockchain) SaveBlockrecycle() error {
	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(recycleFile, data, 0o644); err != nil {
		return err
	}
	return nil
}

// Function to load the blockchain from the JSON file
func (bc *Blockchain) LoadBlockRecycle() error {
	file, err := os.ReadFile(recycleFile)
	if err != nil {
		if os.IsNotExist(err) {
			bc.Recycles = []Recycle{GenerateGenesisRecycle()}
			return bc.SaveBlock()
		}
		return err
	}
	err = json.Unmarshal(file, bc)
	if err != nil {
		return err
	}
	if len(bc.Recycles) == 0 {
		bc.Recycles = []Recycle{GenerateGenesisRecycle()}
	}
	return nil
}
