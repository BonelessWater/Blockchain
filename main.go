package main

import (
	//"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"unicode/utf8"

	"github.com/fxtlabs/primes"
)

func gcf(a, b int) int {
	if b == 0 {
		return a
	}
	return gcf(b, a%b)
}

func find_primes(scale float64) (int, int) {
	boundary := rand.Float64()
	upr := int(math.Round(scale*boundary) + scale)
	arr := primes.Sieve(upr)
	p := arr[len(arr)-1]

	upr = int(math.Round(scale*boundary) + scale)
	arr = primes.Sieve(upr)
	q := arr[len(arr)-2]
	return p, q
}

func floor(a, b int) int { // when dividing integer, go already gives your floor division my default
	return a / b
}

func extendedGCD(a, b int) int { // https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
	old_r, r := a, b
	old_s, s := 1, 0
	old_t, t := 0, 1

	var quotient int
	for r != 0{
		quotient = floor(old_r, r)
		old_r, r = r, old_r - quotient * r
		old_s, s = s, old_s - quotient * s
		old_t, t = t, old_t - quotient * t
	}
	if old_s > 0 {
		return old_s
	} else {
		return old_s + b // adding the final value with lamda if d is a negative number
	}
}

func get_keys() ([2]int, [2]int) {
	scale := 1000.0 // must use numbers with magnitude of four or greater or (e*d)%lambda will yeild negative numbers. Larger numebrs are slower but more secure

	var e, d, phi, n, lambda int

	p, q := find_primes(scale)
	n = p * q
	phi = (p - 1) * (q - 1)

	lambda = phi/gcf(p-1, q-1) // lcm(a, b) where a = p-1 and b=q-1; lcm(a, b) is also equal to abs(a*b)/gcd(a,b)
	
	e = 65537 // e = 2^16 + 1. We uses this number because smaller numbers are known to be less secure in spite of being faster
	fmt.Println(lambda)
	d = extendedGCD(e, lambda) // d ≡ e^-1 (mod λ(n))

	// p, q and lambda must be kept secret because they can be used to compute the sk_pair
	return [2]int{e, n}, [2]int{d, n}
}

func main() {
	pk_pair, sk_pair := get_keys()
	fmt.Printf("Public key pair: %v, Secret key pair: %v\n", pk_pair, sk_pair)

	msg := "Hello" // encrypt and decrypt a letter for starters

	// Convert each character to its ASCII value
	utf := make([]byte, utf8.RuneCountInString(msg)) //  bytesare equivalent to unsigned 8 bit integers in all ways. bytes are just used for convention
	i := 0
	for _, r := range msg {
		utf[i] = byte(r)
		i++
	}

	fmt.Printf("Original string: %s\n", msg)
	fmt.Printf("ASCII values: %v\n", utf)

	enc_msg := make([]int64, len(utf)) // encryption
	exponent := big.NewInt(int64(pk_pair[0])) // use "math/big" package to prevent overflows
	c := new(big.Int)
	for i := 0; i < len(enc_msg); i++ {
		base := big.NewInt(int64(utf[i]))
		c.Exp(base, exponent, nil)
		c.Mod(c, big.NewInt(int64(pk_pair[1])))
		enc_msg[i] = c.Int64()
	}
	fmt.Println(enc_msg)

	dec_msg := make([]int64, len(utf)) // decryption
	exponent = big.NewInt(int64(sk_pair[0]))
	m := new(big.Int)
	for i := 0; i < len(dec_msg); i++ {
		base := big.NewInt(int64(enc_msg[i]))
		m.Exp(base, exponent, nil)
		m.Mod(m, big.NewInt(int64(sk_pair[1])))
		dec_msg[i] = m.Int64()
	}

	fmt.Println(dec_msg)
	var reversedMsg string
	for _, ascii := range dec_msg {
		reversedMsg += string(ascii)
	}
	fmt.Println(reversedMsg)
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
