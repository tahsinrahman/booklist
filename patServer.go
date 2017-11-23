package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/bmizerany/pat"
)

var mu sync.Mutex
var storage = make(map[string][]string)

//add books via get
func addHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get(":name")
	author := r.URL.Query().Get(":auth")

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
	name := r.URL.Query().Get(":name")
	author := r.URL.Query().Get(":auth")

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
	prev := r.URL.Query().Get(":prev")
	new := r.URL.Query().Get(":new")

	fmt.Fprintln(w, storage)

	mu.Lock()
	storage[new] = storage[prev]
	delete(storage, prev)
	mu.Unlock()
	fmt.Fprintln(w, "updated successfully")
}

//update book name
func updateNameHandler(w http.ResponseWriter, r *http.Request) {
	prev := r.URL.Query().Get(":prev")
	new := r.URL.Query().Get(":new")

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
	m.Get("/updateName/:prev/:new", http.HandlerFunc(updateNameHandler))
	m.Get("/updateAuth/:prev/:new", http.HandlerFunc(updateAuthHandler))
	m.Get("/updateAuth/:prev/:new", http.HandlerFunc(updateAuthHandler))

	// Register this pat with the default serve mux so that other packages
	// may also be exported. (i.e. /debug/pprof/*)
	http.Handle("/", m)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
