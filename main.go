package main

import (
	//"crypto/sha256"

	"crypto/sha256"
	"encoding/json"
	"time"
)

type Block struct {
	Index int
	Timestamp int64
	//NewTransactions []Transaction
	Merkleroot [32]byte
	PrevHash [32]byte
	Hash [32]byte
}
type Transaction struct { 
	sender [4]byte
	reciever [4]byte
	amt [4]byte
}
var Blockchain []Block // later this will be initiliazed by json package

func main() {

	if len(Blockchain) == 0 { // Create genesis block if Blockchain is empty
		b := Block {
			Index: 0,
			Timestamp: time.Now().Unix(),
		}
		 // Convert struct to byte slice (JSON representation)
		jsonData, err := json.Marshal(b)
		if err != nil {
			panic(err)
		}
		b.Hash = sha256.Sum256(jsonData)
	}

	// Fetch json
	
	


	// create new block


	//var m_value []int 
	//var rsa_bl []int
	//pointer := 0
	//for i := 0; i < len(asciiValues); i++{ // makes sure the values of m do not exceed n
	//	if m_value[pointer]
	//		m_value[pointer] = asciiValues[i]
	//}


	//hash := sha256.New() // Get the hash of the message
	//hash.Write([]byte(msg))
	//hashSum := hash.Sum(nil) 

	//fmt.Println(hashSum)
}

// make method for verifying if a foreign node added the correct amount of coins to their account