package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type blockchainService interface {
	// Add a new node to the list
	RegisterNode(address []byte) bool

	// Create a new block in the blockchain
	NewBlock(proof []byte, previousHash []byte) *Block

	// Returns the last block in the chain
	LastBlock() Block

	// Simple Proof of Work
	ProofOfWork(lastProof []byte)

	AddBlock(data []byte)
}

type Block struct {
	Index     int64
	Data      []byte
	Nonce     int64
	Timestamp int64
	Previous  []byte
	Proof     []byte
}
type Blockchain struct {
	chain []*Block
}

type Transaction struct {
}

type Token struct {
}

/**
This calculates the proof of the block. This is a total hashvalue of all the properties
*/

var nonce_gl int64

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	r := regexp.MustCompile(`(^[0]{1})`)

	//fmt.Printf("%x\n", r)
	nonce := []byte(strconv.FormatInt(rand.Int63n(78000), 10))
	headers := bytes.Join([][]byte{nonce, b.Previous, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	// Will be used as a mine function later on. This is just for testing purposes
	for !r.Match(hash[:]) {
		nonce = []byte(strconv.FormatInt(rand.Int63n(78000), 10))
		headers = bytes.Join([][]byte{nonce, b.Previous, b.Data, timestamp}, []byte{})
		hash = sha256.Sum256(headers)
		fmt.Printf("%x = %t\n", hash, r.Match(hash[:]))
	}

	//b.Nonce = nonce
	b.Proof = hash[:]
}

var index int64 = 0

func NewBlock(data string, previousHash []byte) *Block {
	index++
	rand.Seed(time.Now().UTC().UnixNano())
	//nonce := rand.Int63n(78000)
	//fmt.Printf("%d\n", nonce)
	block := &Block{index, []byte(data), 10, time.Now().Unix(), previousHash, []byte{}}
	block.SetHash()
	return block
}

func (chain *Blockchain) AddBlock(data string) {
	previousBlock := chain.chain[len(chain.chain)-1]
	newBlock := NewBlock(data, previousBlock.Proof)
	chain.chain = append(chain.chain, newBlock)
}

// We need a starting point, which is why this function is implemented. This is to generate a first time block

func GenerateFirstBlock() *Block {
	return NewBlock("Generic Starting Block", []byte{})
}

func NewBlockChain() *Blockchain {
	return &Blockchain{[]*Block{GenerateFirstBlock()}}
}

func main() {
	createChain := NewBlockChain()

	createChain.AddBlock("This is first addition to the block")
	createChain.AddBlock("This is the second addition to the block")

	for _, block := range createChain.chain {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Prev hash: %x\n", block.Previous)
		fmt.Printf("Block Data: %s\n", block.Data)
		fmt.Printf("Block Proof: %x\n", block.Proof)
		fmt.Println()
	}
}
