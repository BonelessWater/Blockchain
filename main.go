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

type Transaction struct { 
	Timestamp string
	Sender [2]int
	Reciever [2]int
	Amt int
	Signature []int64
}

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

var Blockchain []Block // later this will be initiliazed by json package

func main() {
	// Create genesis block
	var genesisBlock Block

	genesisBlock.Index = 0
	genesisBlock.Timestamp = time.Now().String()
	genesisBlock.Data = 0
	genesisBlock.PrevHash = ""
	genesisBlock.Target = 1
	genesisBlock.Nonce, genesisBlock.Hash = calculateHash(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	// Create new block with fake transaction data
	pk_pair, _ := keys123.Get_keys()
	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], pk_pair)
	if err != nil {
        fmt.Println("Error:", err)
    }
	Blockchain = append(Blockchain, newBlock)
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
	record := currMerkle + string(newTransaction.Sender[0]) + string(newTransaction.Sender[1]) + string(newTransaction.Reciever[0]) + string(newTransaction.Reciever[1]) + string(newTransaction.Amt) + newTransaction.Timestamp
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	hash := hex.EncodeToString(hashed)
	return hash
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, pk_pair [2]int) (Block, error) {
	var newBlock Block
	var MinerReward Transaction
	Merkleroot := oldBlock.Merkleroot

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
	MinerReward.Sender = pk_pair	
	MinerReward.Reciever = pk_pair // public key
	MinerReward.Amt = 50 // coins
	MinerReward.Timestamp = time.Now().String()
	MinerReward.Signature = keys123.Encrypt(string(MinerReward.Sender[0]) + string(MinerReward.Sender[1]) + string(MinerReward.Reciever[0]) + string(MinerReward.Reciever[1]) + string(MinerReward.Amt) + MinerReward.Timestamp)
	newBlock.Merkleroot = merkleHash(Merkleroot, MinerReward)
	newBlock = verifySign(newBlock)
	newBlock.NewTransactions = append(newBlock.NewTransactions, MinerReward)

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Target = oldBlock.Target // the target remains the same
	newBlock.Nonce, newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isHashValid(hash string, target int) bool {
	prefix := strings.Repeat("0", target)
	return strings.HasPrefix(hash, prefix)
}

func verifySign(block Block) Block {
	transactions := block.NewTransactions
	// iterate through transactions later
	transaction := transactions[2]
	sign := transaction.Signature // miner
	decrypted := keys123.Decrypt(sign)

	if decrypted == string(transaction.Sender[0]) + string(transaction.Sender[1]) + string(transaction.Reciever[0]) + string(transaction.Reciever[1]) + string(transaction.Amt) + transaction.Timestamp{
		fmt.Print("nice")
		return block
	} else {
		fmt.Print("not nice")
		return block // remove transaction here
	}
}