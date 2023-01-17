package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func handleUsers(w http.ResponseWriter, r *http.Request) {
	u, err := getUserByEmail("ian@shakespeare.dev")
	if err != nil {
		panic(fmt.Sprintf("Could not find user with email: %s\n", "ian@shakespeare.dev"))
	}
	fmt.Printf("GET /users request: %s %s\n", u.first_name, u.last_name)
	io.WriteString(w, fmt.Sprintf("Hello, %s!\n", u.first_name))
}

func handleFlights(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET /flights request\n")
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
