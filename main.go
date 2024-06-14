package main

import (
	//"crypto/sha256"
	"fmt"
	"math"
	"math/rand"

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
	return old_s
}

func get_keys() ([2]int, [2]int) {
	scale := 1000.0

	var e, d, phi, n, lambda int

	p, q := find_primes(scale)
	n = p * q
	phi = (p - 1) * (q - 1)

	lambda = phi/gcf(p-1, q-1) // lcm(a, b) where a = p-1 and b=q-1; lcm(a, b) is also equal to abs(a*b)/gcd(a,b)
	fmt.Println(lambda)
	e = 7

	d = extendedGCD(e, phi)

	return [2]int{e, n}, [2]int{d, n}
}

func main() {
	pk_pair, sk_pair := get_keys()
	fmt.Printf("Public key pair: %v, Secret key pair: %v\n", pk_pair, sk_pair)

	msg := "Hello"

	asciiValues := make([]int, len(msg))

	// Convert each character to its ASCII value
	for i, char := range msg {
		asciiValues[i] = int(char)
	}

	fmt.Printf("Original string: %s\n", msg)
	fmt.Printf("ASCII values: %v\n", asciiValues)


	fmt.Println((pk_pair[0]*sk_pair[0])%pk_pair[1] == 1)

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
