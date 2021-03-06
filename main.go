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

/*****************working with authorization***************/

type User struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

var userList = make(map[string]User)

//get user from basic auth header or cookie
func checkAuth(r *http.Request) *User {
	username, pass, ok := r.BasicAuth()

	mu.Lock()
	defer mu.Unlock()

	if ok == false {
		//check cookie now
		//return nil
		cookie, err := r.Cookie("login")

		if err != nil {
			return nil
		}

		user, ok := userList[cookie.Value]

		if ok == false {
			return nil
		}

		return &user
	}

	if userList[username].Password == pass {
		user := userList[username]
		return &user
	}

	return nil
}

//get user from json data
func getUser(r *http.Request) *User {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var user User
	err := decoder.Decode(&user)

	if err != nil || user.UserName == "" || user.Password == "" {
		return nil
	}

	return &user
}

//register POST method
func checkRegisterHandler(w http.ResponseWriter, r *http.Request) {
	//if already logged in, can't register until logged out
	if user := checkAuth(r); user != nil {
		http.Redirect(w, r, r.URL.Host+"/", 301)
		return
	}

	user := getUser(r)
	if user == nil {
		authResponse(w, false, "invalid info", "", 400)
		return
	}
	mu.Lock()
	defer mu.Unlock()

	_, ok := userList[user.UserName]
	if ok == true {
		authResponse(w, false, "username exists", user.UserName, 400)
		return
	}

	authResponse(w, true, "registration successfull", user.UserName, 200)

	userList[user.UserName] = *user
}

//register GET method
func registerHandler(w http.ResponseWriter, r *http.Request) {
	//if already logged in, can't register until logged out
	if user := checkAuth(r); user != nil {
		http.Redirect(w, r, r.URL.Host+"/", 301)
		return
	}

	authResponse(w, true, "provide name, username and password now", "", 200)
}

//login POST method
func checkLoginHandler(w http.ResponseWriter, r *http.Request) {
	//if already logged in, can't login until logged out
	if user := checkAuth(r); user != nil {
		http.Redirect(w, r, r.URL.Host+"/", 301)
		return
	}

	user := getUser(r)

	if user == nil {
		authResponse(w, false, "invalid user", "", 400)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, ok := userList[user.UserName]
	if ok == false {
		authResponse(w, false, "username doesn't exist", user.UserName, 400)
		return
	}

	if user.Password == userList[user.UserName].Password {
		//set cookie here
		cookie := http.Cookie{Name: "login", Value: user.UserName, Path: "/"}
		http.SetCookie(w, &cookie)
		authResponse(w, true, "login successfull", user.UserName, 200)

		//after login, cookie is set with response
		//the browser will save this cookie and send it with all next requests
		//so, i need to check if next requests contains cookie or not
	} else {
		authResponse(w, false, "invalid pass", user.UserName, 401)
	}
}

//login GET method
func loginHandler(w http.ResponseWriter, r *http.Request) {
	//if already logged in, can't login until logged out
	if user := checkAuth(r); user != nil {
		http.Redirect(w, r, r.URL.Host+"/", 301)
		return
	}
	authResponse(w, true, "provide username and password now", "", 200)
}

//logout
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	//if not logged in, can't logout until logged in
	user := checkAuth(r)
	if user == nil {
		http.Redirect(w, r, r.URL.Host+"/", 301)
		return
	}

	authResponse(w, true, "successfully logged out", user.UserName, 200)
}

type AuthJson struct {
	Success  bool
	Message  string
	Username string
}

func authResponse(w http.ResponseWriter, status bool, m string, user string, statusCode int) {
	resp := AuthJson{
		Success:  status,
		Message:  m,
		Username: user,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

/*****************working with book***************/

type Book struct {
	Name string
	Auth string
	Id   int
}

var storage = make(map[int]Book)

//add books via post
func addHandler(w http.ResponseWriter, r *http.Request) {
	user := checkAuth(r)

	if user == nil {
		authResponse(w, false, "unauthorized", "", 401)
		return
	}

	book := getBook(r)

	if book == nil {
		commonResponse(w, false, "invalid information", []Book{}, 400)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	counter++
	book.Id = counter

	storage[book.Id] = *book

	commonResponse(w, true, "book added successfully", []Book{*book}, 201)
}

//list books via get
func listHandler(w http.ResponseWriter, r *http.Request) {
	user := checkAuth(r)

	if user == nil {
		authResponse(w, false, "unauthorized", "", 401)
		return
	}

	var listBooks []Book

	mu.Lock()
	defer mu.Unlock()

	for _, books := range storage {
		listBooks = append(listBooks, books)
	}

	commonResponse(w, true, "showing all books", listBooks, 200)
}

//remove books via update
func removeHandler(w http.ResponseWriter, r *http.Request) {
	user := checkAuth(r)

	if user == nil {
		authResponse(w, false, "unauthorized", "", 401)
		return
	}

	s := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(s)

	if err != nil {
		commonResponse(w, false, "invalid link", []Book{}, 400)
		return
	}

	book, ok := storage[id]

	if ok == false {
		commonResponse(w, false, "invalid information", []Book{}, 404)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	delete(storage, id)
	commonResponse(w, true, "deleted book successfully", []Book{book}, 200)
}

//update book via put
func updateHandler(w http.ResponseWriter, r *http.Request) {
	user := checkAuth(r)

	if user == nil {
		authResponse(w, false, "unauthorized", "", 401)
		return
	}

	s := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(s)

	if err != nil {
		commonResponse(w, false, "invalid link", []Book{}, 400)
		return
	}

	book := getBook(r)

	if book == nil {
		commonResponse(w, false, "invalid information", []Book{}, 404)
		return
	}

	_, okay := storage[id]

	if okay == false {
		commonResponse(w, false, "invalid information", []Book{}, 404)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	storage[id] = Book{book.Auth, book.Name, id}
	commonResponse(w, true, "updated book successfully", []Book{storage[id]}, 200)
}

type JsonResponse struct {
	Success  bool
	Message  string
	BookList []Book
}

func commonResponse(w http.ResponseWriter, status bool, m string, list []Book, statusCode int) {
	resp := JsonResponse{
		Success:  status,
		Message:  m,
		BookList: list,
	}
	w.WriteHeader(statusCode)
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

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome")
}

//the client has to send the Authorization header along with every request it makes
//we are working on server side, we've to deal with response, not with the request
//only checking if the request has correct Authorization header or not

//each request will come with cookie

func main() {
	m := pat.New()

	m.Get("/", http.HandlerFunc(homePage))

	m.Get("/book/", http.HandlerFunc(listHandler))
	m.Post("/book/", http.HandlerFunc(addHandler))
	m.Del("/book/:id", http.HandlerFunc(removeHandler))
	m.Put("/book/:id", http.HandlerFunc(updateHandler))

	//when a user cliks on login in get method server checks if he is logged in or not
	//if not logged in, the client then sends username and pass in post method, then server checks credentials
	m.Get("/login/", http.HandlerFunc(loginHandler))
	m.Post("/login/", http.HandlerFunc(checkLoginHandler))

	//when user clicks registration link in get method then he get to the page to fill the reg form
	//after filling the form, the form data will be sent to server using post method
	m.Get("/register/", http.HandlerFunc(registerHandler))
	m.Post("/register/", http.HandlerFunc(checkRegisterHandler))

	m.Get("/logout/", http.HandlerFunc(logoutHandler))

	http.Handle("/", m)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
