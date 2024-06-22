package main

import (
	//"crypto/sha256"

	"blockchain/keys123"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

type Message struct {
	newTransactions []Transaction
}

var Blockchain []Block 
var mutex = &sync.Mutex{}


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		var genesisBlock Block

		genesisBlock.Index = 0
		genesisBlock.Timestamp = time.Now().String()
		genesisBlock.Data = 0
		genesisBlock.PrevHash = ""
		genesisBlock.Target = 1
		genesisBlock.Nonce, genesisBlock.Hash = calculateHash(genesisBlock)
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())
}

// web server
func run() error {
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// write blockchain when we receive an http request
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// takes JSON payload as an input for heart rate transaction request
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	//ensure atomicity when creating new block
	mutex.Lock()

	pk_pair, _ := keys123.Get_keys()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], pk_pair, m.newTransactions)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	_, newHash := calculateHash(newBlock)
	if newHash != newBlock.Hash {
		return false
	}

	return true
}

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
	block.Nonce = nonce
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
func generateBlock(oldBlock Block, pk_pair [2]int, newTransactions []Transaction) (Block) {
	var newBlock Block
	var MinerReward Transaction
	Merkleroot := oldBlock.Merkleroot

//	// Create fake data; otherwise fetch new data from here
//	for i := 0; i < 3; i++ {
//		var transaction Transaction
//		transaction.Sender = [2]int{1, 1} // this is the pk_key pair of a sender
//		transaction.Reciever = [2]int{1, 1} // this it eh sk_key pair of the receiver; doesn't have to be verified
//		transaction.Amt = 50
//		transaction.Timestamp = time.Now().String()
//		transaction.Signature = keys123.Encrypt(string(transaction.Sender[0]) + string(transaction.Sender[1]) + string(transaction.Reciever[0]) + string(transaction.Reciever[1]) + string(transaction.Amt) + transaction.Timestamp)
//		//newBlock.Merkleroot = merkleHash(Merkleroot, transaction)
//		newTransactions = append(newTransactions, transaction)
//	}

	// MinerReward.sender does not have to be defined as they create new money themselves
	MinerReward.Sender = pk_pair	
	MinerReward.Reciever = pk_pair // public key
	MinerReward.Amt = 50 // coins
	MinerReward.Timestamp = time.Now().String()
	MinerReward.Signature = keys123.Encrypt(string(MinerReward.Sender[0]) + string(MinerReward.Sender[1]) + string(MinerReward.Reciever[0]) + string(MinerReward.Reciever[1]) + string(MinerReward.Amt) + MinerReward.Timestamp)
	
	var acceptedTransactions []Transaction
	for i := range newTransactions{
		valid := verifySign(newTransactions[i])
		if valid == true {
			acceptedTransactions = append(acceptedTransactions, newTransactions[i])
			newBlock.Merkleroot = merkleHash(Merkleroot, newTransactions[i])
		}
	}
	newBlock.NewTransactions = append(newBlock.NewTransactions, acceptedTransactions...)
	newBlock.Merkleroot = merkleHash(Merkleroot, MinerReward)
	newBlock.NewTransactions = append(newBlock.NewTransactions, MinerReward)

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Target = oldBlock.Target // the target remains the same
	newBlock.Nonce, newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

func isHashValid(hash string, target int) bool {
	prefix := strings.Repeat("0", target)
	return strings.HasPrefix(hash, prefix)
}

func verifySign(transaction Transaction) bool {
	sign := transaction.Signature // miner
	decrypted := keys123.Decrypt(sign)

	if decrypted == string(transaction.Sender[0]) + string(transaction.Sender[1]) + string(transaction.Reciever[0]) + string(transaction.Reciever[1]) + string(transaction.Amt) + transaction.Timestamp{
		return true
	} else {
		return false
	}
}