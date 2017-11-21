package main

import (
	"math/rand"
	"time"

	"golang.org/x/tour/pic"
)

func Pic(dx, dy int) [][]uint8 {
	ret := make([][]uint8, dy)

	for i := 0; i < dy; i++ {
		ret[i] = make([]uint8, dx)
	}

	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			ret[i][j] = uint8(rand.Uint32())
		}
	}

	return ret
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	pic.Show(Pic)
	//Pic(1, 2)
}
