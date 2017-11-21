package main

import "fmt"

func namedSum(x, y int) (sum int) {
	sum = x + y
	return
}

func main() {
	fmt.Println(namedSum(10, 20))
}
