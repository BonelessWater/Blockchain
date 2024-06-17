package keys123

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"unicode/utf8"

	"github.com/fxtlabs/primes"
)

type Transaction struct { 
	Timestamp string
	Sender [2]int
	Reciever [2]int
	Amt int
	Signature []int64
}

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

func extendedGCD(a, b int) int { // https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
	old_r, r := a, b
	old_s, s := 1, 0
	old_t, t := 0, 1

	var quotient int
	for r != 0{
		quotient = old_r / r // note this is a floor division operator because of the int type
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

func Make_keys() {
	scale := 100.0 // must use numbers with magnitude of four or greater or (e*d)%lambda will yeild negative numbers. Larger numebrs are slower but more secure

	var e, d, phi, n, lambda int

	p, q := find_primes(scale)
	n = p * q
	phi = (p - 1) * (q - 1)

	lambda = phi/gcf(p-1, q-1) // lcm(a, b) where a = p-1 and b=q-1; lcm(a, b) is also equal to abs(a*b)/gcd(a,b)
	
	e = 65537 // e = 2^16 + 1. We uses this number because smaller numbers are known to be less secure in spite of being faster

	d = extendedGCD(e, lambda) // d ≡ e^-1 (mod λ(n))
	// p, q and lambda must be kept secret because they can be used to compute the sk_pair
	line := strconv.Itoa(e) + "," + strconv.Itoa(d) + "," + strconv.Itoa(n)

	f, err := os.Create("keys.txt")
	if err != nil {
		fmt.Println(err)
            	f.Close()
		return
	}
	l, err := f.WriteString(line)
	if err != nil {
		fmt.Println(err)
        f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Get_keys() ([2]int, [2]int) {
	filepath := "keys.txt"
    f, err := os.Open(filepath)
    if err != nil {
        log.Fatal("Unable to read input file " + filepath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for " + filepath, err)
    }
	e, err := strconv.Atoi(records[0][0])
	if err != nil {
        fmt.Println("Error:", err)
    }
	d, err := strconv.Atoi(records[0][1])
	if err != nil {
        fmt.Println("Error:", err)
    }
	n, err := strconv.Atoi(records[0][2])
	if err != nil {
        fmt.Println("Error:", err)
    }
	return [2]int{e, n}, [2]int{d, n}
}

func Encrypt(data string) []int64 {
	// Convert each character to its ASCII value
	pk_pair := [2]int{}

	utf := make([]byte, utf8.RuneCountInString(data)) //  bytesare equivalent to unsigned 8 bit integers in all ways. bytes are just used for convention
	i := 0
	for _, r := range data {
		utf[i] = byte(r)
		i++
	}
	
	enc_msg := make([]int64, len(utf)) // encryption
	exponent := big.NewInt(int64(pk_pair[0])) // use "math/big" package to prevent overflows
	c := new(big.Int)
	for i := 0; i < len(enc_msg); i++ {
		base := big.NewInt(int64(utf[i]))
		c.Exp(base, exponent, nil)
		c.Mod(c, big.NewInt(int64(pk_pair[1])))
		enc_msg[i] = c.Int64()
	}
	return enc_msg
}

func decrypt(value int64, sk_pair [2]int) int64 {
	exponent := big.NewInt(int64(sk_pair[0]))
	m := new(big.Int)
	base := big.NewInt(int64(value))

	m.Exp(base, exponent, nil)
	m.Mod(m, big.NewInt(int64(sk_pair[1])))
	return m.Int64()
}

func Decrypt(enc_msg []int64, sk_pair [2]int) string {
	var wg sync.WaitGroup
	dec_msg := make([]int64, len(enc_msg)) // decryption

	for i := range enc_msg {
		// Increment the WaitGroup counter
		wg.Add(1)
		
		// Capture the index and value for the goroutine
		index := i
		value := enc_msg[i]

		go func() {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes
			dec_msg[index] = decrypt(value, sk_pair)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	var reversedMsg string
	for _, ascii := range dec_msg {
		reversedMsg += string(ascii)
	}
	return reversedMsg
}