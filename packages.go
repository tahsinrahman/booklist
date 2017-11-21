package main

import "fmt"
import "math/rand"

func main() {
	rand.Seed(100)
	for i := 0; i < 10; i++ {
		fmt.Println(rand.Intn(10))
	}
}
