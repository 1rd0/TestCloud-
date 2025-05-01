package main

import (
	"fmt"
	"net/http"
)

// dima loh
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Backend 1 response")
		w.WriteHeader(http.StatusOK)

	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK 1")
	})

	fmt.Println("Server 1 starting on :9001")
	http.ListenAndServe(":9001", nil)
}
