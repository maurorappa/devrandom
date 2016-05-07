package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s - Received request\n", time.Now())
	fmt.Fprintf(w, "Hello")
}

func main() {
	http.HandleFunc("/", handler)
	go http.ListenAndServe(":8080", nil)

	fmt.Scanln()
}
