package main

import (
	"errors"
	"fmt"
	"net/http"

	"flightowl.app/api/server"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.AssignRoutes)

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
