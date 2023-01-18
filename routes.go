package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u, err := getUser(1)
		if err != nil {
			panic("Could not find user with specified id")
		}
		fmt.Printf("GET /users request: %s %s\n", u.FirstName, u.LastName)
		io.WriteString(w, fmt.Sprintf("Hello, %s!\n", u.FirstName))
	case "POST":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic("Empty request body")
		}

		user := User{}
		readError := json.Unmarshal(body, &user)
		if readError != nil {
			panic("Could not read request body")
		}
		fmt.Printf("New User: %s %s %s\n", user.FirstName, user.LastName, user.Email)
		fmt.Printf("New user with email: %s\n", user.Email)
	}
}

func handleFlights(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET /flights request\n")
	io.WriteString(w, "Hello, Flights!\n")
}
