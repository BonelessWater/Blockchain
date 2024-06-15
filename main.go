package main

import (
	//"crypto/sha256"
	"blockchain/crypto"
	"fmt"
)

type Block struct {
	Index int
	Timestamp string
	Transactions int // later change to a merkle root (i think)
	PrevHash string
	Hash string
}

func main() {
	pk_pair, sk_pair := crypto.Get_keys()

	var msg string = "Hello" 

	enc_msg := crypto.Encrypt(msg, pk_pair)

	dec_msg := crypto.Decrypt(enc_msg, sk_pair)

	fmt.Println(msg)
	fmt.Println(enc_msg)
	fmt.Println(dec_msg)
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
