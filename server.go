package main

import (
	"fmt"
	"net/http"
	"strings"
)

//same author can have multiple books, so using slice
var storage = make(map[string][]string)

/*
//book info
type Book struct {
	Name, Author string
}
*/

//default handler
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "visit localhost:8080/add/Name/Author to add book")
	fmt.Fprintln(w, "visit localhost:8080/remove/Name/Author to remove book")
	fmt.Fprintln(w, "visit localhost:8080/search/Name/Author to search book")
	fmt.Fprintln(w, "visit localhost:8080/list to list books")
}

//list books
func listHandler(w http.ResponseWriter, r *http.Request) {
	for name, author := range storage {
		fmt.Fprintln(w, name, author)
	}
}

func parseNameAuth(r *http.Request, title string) (string, string, bool) {
	bookInfo := strings.Split(r.URL.Path[len(title)+2:], "/")

	if len(bookInfo) != 2 {
		return "", "", false
	}
	name, author := bookInfo[0], bookInfo[1]
	return name, author, true
}

//add books via get
func addHandler(w http.ResponseWriter, r *http.Request) {
	name, author, ok := parseNameAuth(r, "add")

	if ok == false {
		//do some work
		return
	}

	storage[author] = append(storage[author], name)

	fmt.Fprintln(w, "book added successfully!")
}

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
	name, author, ok := parseNameAuth(r, "remove")

	if ok == false {
		//do some work
		return
	}

	//first search, if found, then remove
	i, ok := searchBook(name, author)

	if ok == false {
		//do something
		//fmt.Fprint(w, "not found")
		return
	}

	//fmt.Fprintln(w, storage[author])
	storage[author] = append(storage[author][:i], storage[author][i+1:]...)
	fmt.Fprintln(w, "removed successfully")
	//fmt.Fprintln(w, storage[author])
}

//update author name
func updateAuthHandler(w http.ResponseWriter, r *http.Request) {
	prevAuth, newAuth, ok := parseNameAuth(r, "updateAuth")

	if ok == false {
		//do some work
		fmt.Fprintln(w, "not found")
		return
	}

	fmt.Fprintln(w, storage)

	storage[newAuth] = storage[prevAuth]
	fmt.Fprintln(w, storage)
	delete(storage, prevAuth)
	fmt.Fprintln(w, storage)
	fmt.Fprintln(w, "updated successfully")
}

//update book name
func updateNameHandler(w http.ResponseWriter, r *http.Request) {
	prevName, newName, ok := parseNameAuth(r, "updateName")

	if ok == false {
		//do some work
		return
	}

	for auth, _ := range storage {
		i, ok := searchBook(prevName, auth)
		if ok == true {
			storage[auth][i] = newName
			fmt.Fprintln(w, "updated successfully")
			return
		}
	}

	//not found, do something
}

func main() {
	fmt.Println("starting server on port 8080... ")

	http.HandleFunc("/", handler)
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/list/", listHandler)
	http.HandleFunc("/remove/", removeHandler)
	http.HandleFunc("/updateName/", updateNameHandler)
	http.HandleFunc("/updateAuth/", updateAuthHandler)
	http.ListenAndServe(":8080", nil)
}
