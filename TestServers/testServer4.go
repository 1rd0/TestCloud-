package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Backend 4 response")
		w.WriteHeader(http.StatusOK)
		log.Print("4")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK 4")
	})

	fmt.Println("Server 3 starting on :9004")
	http.ListenAndServe(":9004", nil)
}
