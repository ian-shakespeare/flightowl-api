package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/arcticstorm9/flightowl-api/database"
	"github.com/arcticstorm9/flightowl-api/server"
)

func main() {
	err := database.Init()
	if err != nil {
		panic("could not create database")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.AssignRoutes)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	fmt.Println("Listening...")

	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
