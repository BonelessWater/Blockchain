package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/fxtlabs/primes"
)

func gcf(n1 int, n2 int) int {
   if (n2 != 0) {
      return gcf(n2, n1 % n2);
   } else {
      return n1;
   }
}

func find_primes(scale float64) (int, int){
	boundary := rand.Float64() // Generates a random integer between 0 and 1
	
	upr := int(math.Round(scale * boundary) + scale) // gets upper bound of first prime
	arr := primes.Sieve(upr)

	p := arr[len(arr) - 1]

	upr = int(math.Round(scale * boundary) + scale) // gets upper bound of first prime
	arr = primes.Sieve(upr)

	q := arr[len(arr) - 1]

	return p, q
}

func rsa() ([2]int, [2]int) {
	var scale float64 = 100.0
    
	valid_primes := false
	var e int
	var phi int
	var n int

	for !valid_primes{

		p, q := find_primes(scale)
		n = p * q

		phi = (p-1) * (q-1)

		for i := phi-1; i >= 1 ; i-- { // finds a number e where, 1 < e < phi and the gfc of e and pdi is 1 (e and phi are coprime )
			gfc := gcf(phi, i)
			if gfc == 1{
				valid_primes = true
				e = i
				break
			}
		}
	}
	
	mult_inverse_fd := false // finding modular multiplicative inverse of e modulo phi. there is only one value that satifies this equation: d*e≡1 (mod ϕ(n))
	var d int
	for !mult_inverse_fd{
		d = 1
		if (d * e) % phi == (1 % phi){
			mult_inverse_fd = true
		}
	}

	var pk_pair [2]int = [2]int{e, n}
	var sk_pair [2]int = [2]int{d, n}

	return pk_pair, sk_pair
}

func main() {
	pk_pair, sk_pair := rsa()
	fmt.Printf("Public key pair: %v, Secret key pair: %v ", pk_pair, sk_pair)
}