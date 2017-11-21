package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Print("Go runs on ")

	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("osx")
	case "linux":
		fmt.Println("linux")
	default:
		fmt.Println("windows")
	}
}
