package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	x := strings.Split(s, " ")

	for _, val := range x {
		m[val]++
	}

	return m
}

func main() {
	wc.Test(WordCount)
}
