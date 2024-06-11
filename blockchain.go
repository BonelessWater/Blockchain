package main // this ie needed for an excecutable program

import ( // must use imported libraries or Go throws an error
	"fmt" // library for input and output streams
	"time"
)
func add(a int, b int) int { // input varaibles need their types declared. return type also needed
    return a + b
}

func divide(a, b int) (int, int) { // this is how to return multiple types
    quotient := a / b
    remainder := a % b
    return quotient, remainder
}

type Person struct {
    Name string
    Age  int
}

func (p Person) Greet() { // (self, name of struct) the variable p is acting similar to self in python classes
    fmt.Printf("Hello, my name is %s and I am %d years old.\n", p.Name, p.Age)
}

func printNumbers() {
    for i := 1; i <= 5; i++ {
        fmt.Println(i)
        time.Sleep(100 * time.Millisecond)
    }
}

func sum(a, b int, c chan int) {
    c <- a + b
}

func main() {
    fmt.Println("hello world") // notice how we treat fmt as a class

	var name string = "Go" // we must declare types when new varaibles are initialized. notice how the type goes after the varaible name unlike c
    var age int = 10
    var version float32 = 1.15
    isAwesome := true // shorthand for declaring and initializing a variable

    fmt.Println(name, age, version, isAwesome)

	year := 20 // we can use := to implicitly define types in Go. just like in Godot :)

    if year < 18 {
        fmt.Println("You are a minor.")
    } else {
        fmt.Println("You are an adult.")
    }

	for i := 0; i < 5; i++ { // for loops in godot are identical to those in js but without the ()
        fmt.Println(i)
    }

    // while-style loop
    n := 1
    for n < 5 {
        fmt.Println(n)
        n++
    }

    // infinite loop
    /* Uncomment to run infinite loop
    for {
        fmt.Println("Infinite loop")
    }
    */

	result := add(3, 4)
    fmt.Println("Sum:", result)

	q, r := divide(10, 3)
    fmt.Println("Quotient:", q, "Remainder:", r)

	var arr [5]int // declare an array of 5 integers. declare type after defining array length
    arr[0] = 1 //arrays are indexable (thank god)
    fmt.Println(arr) // arrays can only be one type like c

	slice := []int{1, 2, 3, 4, 5} // this is how to make an array without having to define it first in another line
    fmt.Println(slice)
    slice = append(slice, 6) // claaassicc python stuff
    fmt.Println(slice)

	m := make(map[string]int) // umm looks like an overcomplicated dictionary declaration
    m["one"] = 1
    m["two"] = 2

    fmt.Println(m)
    delete(m, "one")
    fmt.Println(m)

	person := Person{Name: "Alice", Age: 30} // LETSS GOOO we have structs ladies and gentlemen
    fmt.Println(person)

    person.Greet() // we have struct-specific functions too? interesting how we might be able to use those later

	go printNumbers() // what... you can run multiple things at the same time
    fmt.Println("Goroutine started")
    time.Sleep(1 * time.Second)
    fmt.Println("Main function ends")

	c := make(chan int, 2) // c is a channel and is supposed to be an integer. the number 2 is to indicate the capacity buffer of the channel
    go sum(1, 2, c)
    go sum(3, 4, c)

    result1 := <-c // the value of c changes over time
    result2 := <-c

    fmt.Println("Results:", result1, result2)
}