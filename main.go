package main

import (
	//"crypto/sha256"

	"blockchain/keys123"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

var BlockChain []Block

type Block struct {
	Index int
	Timestamp string
	Data int
	NewTransactions []Transaction
	Merkleroot string
	PrevHash string
	Target int
	Nonce int
	Hash string
}
type Transaction struct { 
	Timestamp string
	Sender [2]int
	Reciever [2]int
	Amt int
	Signature []int64
}
var Blockchain []Block // later this will be initiliazed by json package

func main() {
	keys123.Make_keys()
	x, y := keys123.Get_keys()
	fmt.Println(x, y)

	// Create genesis block
	var genesisBlock Block

	t := time.Now()

	genesisBlock.Index = 0
	genesisBlock.Timestamp = t.String()
	genesisBlock.Data = 0
	genesisBlock.PrevHash = ""
	genesisBlock.Target = 1
	genesisBlock.Nonce, genesisBlock.Hash = calculateHash(genesisBlock)

	Blockchain = append(Blockchain, genesisBlock)

	pk_pair, sk_pair := keys123.Get_keys()
	newBlock, err := generateBlock(BlockChain[len(Blockchain)-1], pk_pair, sk_pair)
	if err != nil {
        fmt.Println("Error:", err)
    }
	BlockChain = append(BlockChain, newBlock)
}

// make method for verifying if a foreign node added the correct amount of coins to their account

// SHA256 hashing
func calculateHash(block Block) (int, string) {
	valid := false
	nonce := 0
	var hash string
	for !valid {
		nonce++
		record := string(block.Index) + block.Timestamp + string(block.Data) + block.PrevHash + string(block.Target) + string(nonce)
		h := sha256.New()
		h.Write([]byte(record))
		hashed := h.Sum(nil)
		hash = hex.EncodeToString(hashed)
		valid = isHashValid(hash, block.Target)
	}
	return nonce, hash
}

func merkleHash(currMerkle string, newTransaction Transaction) string {
	record := currMerkle + string(newTransaction.sender) + string(newTransaction.reciever) + string(newTransaction.amt)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	hash = hex.EncodeToString(hashed)
	return hash
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, pk_pair [2]int, sk_pair [2]int) (Block, error) {
	var newBlock Block
	var MinerReward Transaction
	Merkleroot := oldBlock.Merkleroot
	t := time.Now()

	// Create fake data; otherwise fetch new data from here
	for i := 0; i < 3; i++ {
		var transaction Transaction
		transaction.Sender = [2]int{1, 1} // this is the pk_key pair of a sender
		transaction.Reciever = [2]int{1, 1} // this it eh sk_key pair of the receiver; doesn't have to be verified
		transaction.Amt = 50
		transaction.Timestamp = time.Now().String()
		transaction.Signature = keys123.Encrypt(string(transaction.Sender[0]) + string(transaction.Sender[1]) + string(transaction.Reciever[0]) + string(transaction.Reciever[1]) + string(transaction.Amt) + transaction.Timestamp)
		newBlock.Merkleroot = merkleHash(Merkleroot, transaction)
		newBlock.NewTransactions = append(newBlock.NewTransactions, transaction)
	}

	// MinerReward.sender does not have to be defined as they create new money themselves
	MinerReward.Reciever = pk_pair[0] // public key
	MinerReward.Amt = 50 // coins
	MinerReward.Timestamp = time.Now().String()
	newBlock.Merkleroot = merkleHash(Merkleroot, MinerReward)
	newBlock.NewTransactions = append(newBlock.NewTransactions, MinerReward)

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Nonce, newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isHashValid(hash string, target int) bool {
	prefix := strings.Repeat("0", target)
	return strings.HasPrefix(hash, prefix)
}