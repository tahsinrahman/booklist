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
var storage = make(map[string][]string)

type Book struct {
	Name string
	Auth string
}

func getNames(r *http.Request) (string, string) {
	method := r.Method
	fmt.Println(method)

	var a, b string

	if method == "GET" {
		a = r.URL.Query().Get(":name")
		b = r.URL.Query().Get(":auth")
	} else {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var book Book
		err := decoder.Decode(&book)

		fmt.Println(book)

		if err != nil {
			fmt.Println(err)
			return "", ""
		}
		a, b = book.Name, book.Auth
	}
	return a, b
}

//add books via get
func addHandler(w http.ResponseWriter, r *http.Request) {
	name, author := getNames(r)

	fmt.Println(name, author)

	mu.Lock()
	storage[author] = append(storage[author], name)
	mu.Unlock()

	fmt.Fprintln(w, "book added successfully!")
}

//list books
func listHandler(w http.ResponseWriter, r *http.Request) {
	for name, author := range storage {
		fmt.Fprintln(w, name, author)
	}
}

//search book
func searchBook(name, auth string) (int, bool) {
	for i, x := range storage[auth] {
		if x == name {
			return i, true
		}
	}
	return 0, false
}

//remove books
func removeHandler(w http.ResponseWriter, r *http.Request) {
	name, author := getNames(r)

	//first search, if found, then remove
	i, ok := searchBook(name, author)

	if ok == false {
		//do something
		//fmt.Fprint(w, "not found")
		return
	}

	//fmt.Fprintln(w, storage[author])
	mu.Lock()
	storage[author] = append(storage[author][:i], storage[author][i+1:]...)
	mu.Unlock()
	fmt.Fprintln(w, "removed successfully")
	//fmt.Fprintln(w, storage[author])
}

//update author name
func updateAuthHandler(w http.ResponseWriter, r *http.Request) {
	prev, new := getNames(r)
	fmt.Println(1, prev, new)

	mu.Lock()
	storage[new] = append(storage[new], storage[prev]...)
	delete(storage, prev)
	mu.Unlock()

	fmt.Println(storage)

	fmt.Fprintln(w, "updated successfully")
}

//update book name
func updateNameHandler(w http.ResponseWriter, r *http.Request) {
	prev, new := getNames(r)

	for auth, _ := range storage {
		i, ok := searchBook(prev, auth)
		if ok == true {
			mu.Lock()
			storage[auth][i] = new
			mu.Unlock()
			fmt.Fprintln(w, "updated successfully")
			return
		}
	}

	//not found, do something
}

func main() {
	m := pat.New()
	m.Get("/add/:name/:auth", http.HandlerFunc(addHandler))
	m.Get("/list/", http.HandlerFunc(listHandler))
	m.Get("/remove/:name/:auth", http.HandlerFunc(removeHandler))
	m.Get("/updateName/:name/:auth", http.HandlerFunc(updateNameHandler))
	m.Get("/updateAuth/:name/:auth", http.HandlerFunc(updateAuthHandler))
	m.Get("/updateAuth/:name/:auth", http.HandlerFunc(updateAuthHandler))

	m.Post("/add/", http.HandlerFunc(addHandler))
	m.Post("/list/", http.HandlerFunc(listHandler))
	m.Post("/remove/", http.HandlerFunc(removeHandler))
	m.Post("/updateName/", http.HandlerFunc(updateNameHandler))
	m.Post("/updateAuth/", http.HandlerFunc(updateAuthHandler))
	m.Post("/updateAuth/", http.HandlerFunc(updateAuthHandler))

	// Register this pat with the default serve mux so that other packages
	// may also be exported. (i.e. /debug/pprof/*)
	http.Handle("/", m)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
