package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Backend 3 response")
		w.WriteHeader(http.StatusOK)
		log.Print("3")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK 3")
	})

	fmt.Println("Server 3 starting on :9003")
	http.ListenAndServe(":9003", nil)
}
