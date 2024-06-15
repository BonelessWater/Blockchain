package main

import (
	//"crypto/sha256"

	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type Block struct {
	Index int
	Timestamp string
	Data int
	//NewTransactions []Transaction
	//Merkleroot [32]byte
	PrevHash string
	Hash string
}
type Transaction struct { 
	sender [4]byte
	reciever [4]byte
	amt [4]byte
}
var Blockchain []Block // later this will be initiliazed by json package

var bcServer chan []Block // bcServer handles incoming concurrent Blocks
var mutex = &sync.Mutex{}


func main() {

	if len(Blockchain) == 0 { // Create genesis block if Blockchain is empty
		b := Block {
			Index: 0,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		b.Hash = calculateHash(b)
		spew.Dump(b)
		Blockchain = append(Blockchain, b)
	}

	tcpPort := os.Getenv("PORT") // gets environment variable

	server, err := net.Listen("tcp", ":"+tcpPort) // start TCP and serve TCP server
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP  Server Listening on port :", tcpPort)
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
	// Fetch json
	
	


	// create new block

	//hash := sha256.New() // Get the hash of the message
	//hash.Write([]byte(msg))
	//hashSum := hash.Sum(nil) 

	//fmt.Println(hashSum)
}

// make method for verifying if a foreign node added the correct amount of coins to their account

func handleConn(conn net.Conn) {

	defer conn.Close()

	io.WriteString(conn, "Who do you want to send money to: ")

	scanner := bufio.NewScanner(conn)

	// take in data from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			data, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], data)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

			bcServer <- Blockchain
			io.WriteString(conn, "\nEnter a new number:")
		}
	}()

	// simulate receiving broadcast
	go func() {
		for {
			time.Sleep(30 * time.Second)
			mutex.Lock()
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range bcServer {
		spew.Dump(Blockchain)
	}

}

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// make sure the chain we're checking is longer than the current blockchain
func replaceChain(newBlocks []Block) {
	mutex.Lock()
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
	mutex.Unlock()
}

// SHA256 hasing
func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.Data) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, Data int) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = Data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}