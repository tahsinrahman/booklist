package main

import (
	"fmt"
	"net/http"
)

type Book struct {
	string Name, Author
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello!"+"visit localhost:8080/add/Name/Author to add book")
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	bookInfo := 
}

func main() {
	fmt.Println("starting server...")

	http.HandleFunc("/", handler)
	http.HandleFunc("/add/", addHandler)
	http.ListenAndServe(":8080", nil)
}
