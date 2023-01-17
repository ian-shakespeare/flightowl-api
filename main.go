package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func handleUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("go /users request\n")
	io.WriteString(w, "Hello, Users!\n")
}

func handleFlights(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /flights request\n")
	io.WriteString(w, "Hello, Flights!\n")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", handleUsers)
	mux.HandleFunc("/flights", handleFlights)

	port := "8000"

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	fmt.Printf("Listening on port %s\n", port)

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
