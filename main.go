package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", test)
	mux.HandleFunc("POST /add", createTodo)
	mux.HandleFunc("GET /get/{id}", getTodo)
	mux.HandleFunc("DELETE /get/{id}", deleteTodo)

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

// get todo by id
func getTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	todo, ok := todoCache[id]
	mutex.Unlock()

	if ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	} else {
		http.Error(w, "Not Found", http.StatusBadRequest)
		return
	}
}

// to delete todo
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := todoCache[id]; !ok {
		http.Error(w, "Not Found", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	delete(todoCache, id)
	mutex.Unlock()

	type t_err struct {
		Msg    string `json:"msg"`
		Status uint   `json:"status"`
	}

	json.NewEncoder(w).Encode(t_err{"Deleted", http.StatusOK})

}
