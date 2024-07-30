package project

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

const difficulty = 4 // Adjust difficulty as needed (higher number = harder)

type Block struct {
	Index        int
	Timestamp    int64
	Data         string
	PreviousHash []byte
	Hash         []byte
	Nonce        int
}

type Blockchain struct {
	Blocks []*Block
}

func NewBlock(data string, previousHash []byte) *Block {
	block := &Block{
		Index:        len(previousHash),
		Timestamp:    time.Now().Unix(),
		Data:         data,
		PreviousHash: previousHash,
		Hash:         []byte{},
		Nonce:        0,
	}
	block.Hash = block.calculateHash()
	return block
}

func (b *Block) calculateHash() []byte {
	// Concatenate block data using efficient byte manipulation
	data := bytes.Join(
		[][]byte{
			b.PreviousHash,
			[]byte(b.Data),
			[]byte(strconv.Itoa(b.Index)),
			[]byte(strconv.FormatInt(b.Timestamp, 10)),
			[]byte(strconv.Itoa(b.Nonce)),
		},
		[]byte{},
	)

	// Use sha256.Sum256 for a more concise approach
	hash := sha256.Sum256(data)
	return hash[:]
}

func (b *Block) String() string {
	return fmt.Sprintf("Index: %d, Timestamp: %d, Data: %s, PreviousHash: %x, Hash: %x, Nonce: %d", b.Index, b.Timestamp, b.Data, b.PreviousHash, b.Hash, b.Nonce)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) PrintChain() {
	for _, block := range bc.Blocks {
		fmt.Println(block)
	}
}

func (b *Block) mineBlock(difficulty int) {
	for {
		if bytes.Compare(b.Hash[0:difficulty], bytes.Repeat([]byte{0}, difficulty)) == 0 {
			break
		}
		b.Nonce++
		b.Hash = b.calculateHash()
	}
	fmt.Printf("Block mined successfully: %x\n", b.Hash)
}
