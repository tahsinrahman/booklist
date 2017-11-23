package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

var storage = make(map[int]Book)

type JsonResponse struct {
	Success  bool
	Message  string
	BookList []Book
}

//add books via post
func addHandler(w http.ResponseWriter, r *http.Request) {
	book := getBook(r)

	if book == nil {
		commonResponse(w, false, "invalid information", []Book{})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	counter++
	book.Id = counter

	storage[book.Id] = *book

	commonResponse(w, true, "book added successfully", []Book{*book})
}

//list books via get
func listHandler(w http.ResponseWriter, r *http.Request) {
	var listBooks []Book

	mu.Lock()
	defer mu.Unlock()

	for _, books := range storage {
		listBooks = append(listBooks, books)
	}

	commonResponse(w, true, "showing all books", listBooks)
}

//remove books via update
func removeHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(s)

	if err != nil {
		commonResponse(w, false, "invalid information", []Book{})
		return
	}

	book, ok := storage[id]

	if ok == false {
		commonResponse(w, false, "invalid information", []Book{})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	delete(storage, id)
	commonResponse(w, true, "deleted book successfully", []Book{book})
}

//update book via put
func updateHandler(w http.ResponseWriter, r *http.Request) {
	book := getBook(r)
	fmt.Println(book)

	if book == nil {
		commonResponse(w, false, "invalid information", []Book{})
		return
	}

	s := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(s)

	if err != nil {
		commonResponse(w, false, "invalid information", []Book{})
		return
	}

	_, ok := storage[id]

	if ok == false {
		commonResponse(w, false, "invalid information", []Book{})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	storage[id] = Book{book.Auth, book.Name, id}
	commonResponse(w, true, "updated book successfully", storage[id])
}

func commonResponse(w http.ResponseWriter, status bool, m string, list []Book) {
	resp := JsonResponse{
		Success:  status,
		Message:  m,
		BookList: list,
	}
	json.NewEncoder(w).Encode(resp)
}

func getBook(r *http.Request) *Book {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var book Book
	err := decoder.Decode(&book)

	if err != nil {
		return nil
	}
	return &book
}

func main() {
	m := pat.New()
	m.Get("/book/", http.HandlerFunc(listHandler))
	m.Post("/book/", http.HandlerFunc(addHandler))
	m.Del("/book/:id", http.HandlerFunc(removeHandler))
	m.Put("/book/:id", http.HandlerFunc(updateHandler))

	http.Handle("/", m)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
