package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/bmizerany/pat"
)

var mu sync.Mutex
var counter int

type Book struct {
	Name string
	Auth string
	Id   int
}

var storage = make(map[string][]Book)

type JsonResponse struct {
	Success  int
	Message  string
	BookList []Book
}

//add books via get
func addHandler(w http.ResponseWriter, r *http.Request) {
	name, author := getNames(r)

	mu.Lock()
	defer mu.Unlock()

	counter++
	book := Book{name, author, counter}

	storage[author] = append(storage[author], book)

	resp := JsonResponse{
		Success:  1,
		Message:  "book added successfully",
		BookList: []Book{book},
	}
	json.NewEncoder(w).Encode(resp)
}

//list books
func listHandler(w http.ResponseWriter, r *http.Request) {
	var listBooks []Book

	for _, books := range storage {
		listBooks = append(listBooks, books[:]...)
	}

	resp := JsonResponse{
		Success:  1,
		Message:  "showing all books",
		BookList: listBooks,
	}
	json.NewEncoder(w).Encode(resp)
}

//remove books
func removeHandler(w http.ResponseWriter, r *http.Request) {
	name, author := getNames(r)

	//first search, if found, then remove
	i, ok := searchBook(name, author)

	if ok == false {
		resp := JsonResponse{
			Success: 0,
			Message: "the book isn't available",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	book := storage[author][i]

	storage[author] = append(storage[author][:i], storage[author][i+1:]...)

	resp := JsonResponse{
		Success:  1,
		Message:  "book removed successfully",
		BookList: []Book{book},
	}
	json.NewEncoder(w).Encode(resp)
}

func getNames(r *http.Request) (string, string) {
	method := r.Method

	var a, b string

	if method == "GET" {
		a = r.URL.Query().Get(":name")
		b = r.URL.Query().Get(":auth")
	} else {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var book Book
		err := decoder.Decode(&book)

		if err != nil {
			fmt.Println(err)
			return "", ""
		}
		a, b = book.Name, book.Auth
	}
	return a, b
}

//search book
func searchBook(name, auth string) (int, bool) {
	mu.Lock()
	defer mu.Unlock()
	for i, x := range storage[auth] {
		if x.Name == name {
			return i, true
		}
	}
	return 0, false
}

/*
//update author name
func updateAuthHandler(w http.ResponseWriter, r *http.Request) {
	prev, new := getNames(r)
	fmt.Println(1, prev, new)

	mu.Lock()
	defer mu.Unlock()

	storage[new] = append(storage[new], storage[prev]...)
	delete(storage, prev)

	fmt.Fprintln(w, "updated successfully")
}

//update book name
func updateNameHandler(w http.ResponseWriter, r *http.Request) {
	prev, new := getNames(r)

	for auth, _ := range storage {
		i, ok := searchBook(prev, auth)
		if ok == true {
			mu.Lock()
			defer mu.Unlock()
			storage[auth][i] = new
			fmt.Fprintln(w, "updated successfully")
			return
		}
	}

	//not found, do something
}
*/

func main() {
	m := pat.New()
	m.Get("/add/:name/:auth", http.HandlerFunc(addHandler))
	m.Get("/list/", http.HandlerFunc(listHandler))
	m.Get("/remove/:name/:auth", http.HandlerFunc(removeHandler))
	//m.Get("/updateName/:name/:auth", http.HandlerFunc(updateNameHandler))
	//m.Get("/updateAuth/:name/:auth", http.HandlerFunc(updateAuthHandler))
	//m.Get("/updateAuth/:name/:auth", http.HandlerFunc(updateAuthHandler))

	m.Post("/add/", http.HandlerFunc(addHandler))
	m.Post("/list/", http.HandlerFunc(listHandler))
	m.Post("/remove/", http.HandlerFunc(removeHandler))
	//m.Post("/updateName/", http.HandlerFunc(updateNameHandler))
	//m.Post("/updateAuth/", http.HandlerFunc(updateAuthHandler))
	//m.Post("/updateAuth/", http.HandlerFunc(updateAuthHandler))

	// Register this pat with the default serve mux so that other packages
	// may also be exported. (i.e. /debug/pprof/*)
	http.Handle("/", m)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
