package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", test)
	mux.HandleFunc("POST /add", createTodo)

	fmt.Println("Hello API")

	//Listening
	log.Fatal((http.ListenAndServe(":8000", mux)))
}

// struct for TODO
type TODO struct {
	Title string `json:"title"`
	Do    string `json:"do"`
}

// todo cache. local temp db for practice
var todoCache = make(map[int]TODO)

// to test the server
func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API WORKING")
}

// for limit access
var mutex sync.RWMutex

// to post todo
func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo TODO
	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if todo.Do == "" || todo.Title == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	todoCache[len(todoCache)+1] = todo
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Println(todoCache)
}
