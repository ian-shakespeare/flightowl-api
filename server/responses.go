package server

import (
	"fmt"
	"net/http"
)

func handleNotFound(w http.ResponseWriter) {
	fmt.Printf("- 404\n")
	w.WriteHeader(http.StatusNotFound)
}

func handleBadRequest(w http.ResponseWriter) {
	fmt.Printf("- 400\n")
	w.WriteHeader(http.StatusBadRequest)
}

func handleUnauthorized(w http.ResponseWriter) {
	fmt.Printf("- 401\n")
	w.WriteHeader(http.StatusUnauthorized)
}

func handleConflict(w http.ResponseWriter) {
	fmt.Printf("- 409\n")
	w.WriteHeader(http.StatusConflict)
}

func handleOK(w http.ResponseWriter) {
	fmt.Printf("- 200\n")
}

func handleCreated(w http.ResponseWriter) {
	fmt.Printf("- 201\n")
	w.WriteHeader(http.StatusCreated)
}
