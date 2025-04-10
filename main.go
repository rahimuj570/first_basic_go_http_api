package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", test)
	fmt.Println("Hello API")

	//Listening
	http.ListenAndServe(":8000", mux)
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API WORKING")
}
