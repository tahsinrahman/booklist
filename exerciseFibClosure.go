package main

import "fmt"

func fib() func() int {
	x, y := 0, 1

	return func() int {
		ret := x
		x, y = y, x+y
		return ret
	}
}

func main() {
	f := fib()

	for i := 0; i < 10; i++ {
		fmt.Println(f())
		//fmt.Println("%T\n", fib())
	}
}
