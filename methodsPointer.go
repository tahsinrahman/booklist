package main

import (
	"fmt"
	"math"
)

type vertex struct {
	x, y float64
}

func (v *vertex) abs() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *vertex) scale(x float64) {
	v.x *= x
	v.y *= x
}

func Test(x ...int) {
	for _, val := range x {
		fmt.Println(val)
	}
}

func main() {
	v := vertex{4, 3}
	fmt.Println(v, v.abs())
	v.scale(10)
	fmt.Println(v, v.abs())

	p := &v
	fmt.Println(*p, p.abs())
	p.scale(10)
	fmt.Println(*p, p.abs())

	Test(1, 2, 3, 4, 5, 6, 7)
}
