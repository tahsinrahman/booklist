package main

import "fmt"

type newName int

func (x newName) Abs() int {
	if x > 0 {
		return int(x)
	} else {
		return int(-x)
	}
}

func main() {
	var x, y newName
	x, y = 30, -30
	fmt.Println(x, y)
	fmt.Println(x.Abs(), y.Abs())
}
