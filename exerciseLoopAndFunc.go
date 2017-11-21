package main

/*
import (
	"fmt"
	"math"
	"math/rand"
)
*/

//why gofmt doesn't work?
//generate same random numbers inspite of using seed?

import "fmt"
import "math"
import "math/rand"

const Small = 1e-9

func sqrt(x float64) (float64, int) {
	z := 1.0

	itr := 0
	for math.Abs(z*z-x) > Small {
		z -= (z*z - x) / (2 * z)
		itr++
	}

	return z, itr
}

func sqrtWithRanInit(x float64) (float64, int) {
	//z := float64(rand.Intn(int(x)) + 1)
	rand.Seed(100)
	z := rand.Float64()
	fmt.Println(z)
	//fmt.Printf("%T\n", z)

	itr := 0
	for math.Abs(z*z-x) > Small {
		z -= (z*z - x) / (2 * z)
		itr++
	}

	return z, itr
}

func main() {
	fmt.Println(sqrt(2))

	for i := 0; i <= 10; i++ {
		fmt.Println(sqrtWithRanInit(2))
	}
}
