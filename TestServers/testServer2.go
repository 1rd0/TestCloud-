package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Backend 2 response")
		w.WriteHeader(http.StatusOK)
		log.Print("2")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK 2")
	})

	fmt.Println("Server 2 starting on :9002")
	http.ListenAndServe(":9002", nil)
}
