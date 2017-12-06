package main

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
	"math"
	"encoding/binary"
    "log"
)

type blockchainService interface {
	// Add a new node to the list
	RegisterNode(address []byte) bool

	// Create a new block in the blockchain
	NewBlock(data string, previousHash []byte) *Block

	// Returns the last block in the chain
	LastBlock() Block

	// Simple Proof of Work
	NewProofOfWork(block *Block) *ProofOfWork

	AddBlock(data string)
}

type Block struct {
	Index 			int64
	Timestamp		int64
	Data			[]byte
	PrevHash		[]byte
	Hash			[]byte
	Nonce			int
}

type Blockchain struct {
	blocks 			[]*Block
}

type ProofOfWork struct {
	block *Block
	target *big.Int
}

var index int64 = 0 
const targetBits = 16


/**
This calculates the proof of the block. This is a total hashvalue of all the properties
*/
func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{block.PrevHash, block.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	block.Hash = hash[:]
}
// This function creates a new Block 
func NewBlock(data string, prevHash []byte) *Block {
	index++
	block := &Block{index, time.Now().Unix(), []byte(data), prevHash, []byte{}, 0}
	proofOfWork := NewProofOfWork(block)
	nonce, hash := proofOfWork.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block	
}

// Here we add the Block to a Blockchain with the previous Block's hash 
func (blockchain *Blockchain) AddBlock(data string) {
	prevBlock := blockchain.blocks[len(blockchain.blocks) -1]
	newBlock := NewBlock(data, prevBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, newBlock)
}

// Here we create the first (Genesis) Block
func NewGenesisBlock() *Block {
	return NewBlock("This is the Genesis Block", []byte{})
}

// And here we use that Genesis Block on the Blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	proofOfWork := &ProofOfWork{block, target}

	return proofOfWork
}

func IntToHex(num int64) []byte {
    buff := new(bytes.Buffer)
    err := binary.Write(buff, binary.BigEndian, num)
    if err != nil {
        log.Panic(err)
    }

    return buff.Bytes()
}

func (proofOfWork *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.block.PrevHash,
			proofOfWork.block.Data,
			IntToHex(proofOfWork.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (proofOfWork *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var maxNonce = math.MaxInt64
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", proofOfWork.block.Data)
	for nonce < maxNonce {
		data := proofOfWork.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(proofOfWork.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("\n\n")

	return nonce, hash[:]
} 

func (proofOfWork *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := proofOfWork.prepareData(proofOfWork.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(proofOfWork.target) == -1

	return isValid
}

func main() {
	blockchain := NewBlockchain()
	
		blockchain.AddBlock("This is first addition to the block")
		blockchain.AddBlock("This is the second addition to the block")
		blockchain.AddBlock("This is the third addition to the blockf")		
	
		for _, block := range blockchain.blocks {
			fmt.Printf("Timestamp: %d\n", block.Timestamp)
			fmt.Printf("Index: %d\n", block.Index)			
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
			fmt.Printf("Data: %s\n", block.Data)
			fmt.Printf("Hash: %x\n", block.Hash)
			proofOfWork := NewProofOfWork(block)
			fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(proofOfWork.Validate()))
			fmt.Println()
		}
}