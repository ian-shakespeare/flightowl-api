package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"flightowl.app/api/helpers"
)

func handleError(w http.ResponseWriter, msg string, code int) {
	fmt.Printf("- %d\n", code)
	http.Error(w, msg, code)
}

func doGET(path string, w http.ResponseWriter, r *http.Request) {
	if helpers.StartsWith(path, "/users") {
		id, err := strconv.Atoi(helpers.GetPathCollection(path))
		if err != nil {
			handleError(w, "Bad Request Resource", 400)
			return
		}
		usersGet(id, w, r)
	} else {
		handleError(w, "Not Found", 404)
	}
}

func doPOST(path string, w http.ResponseWriter, r *http.Request) {
	if path == "/users" {
		usersCreate(w, r)
	} else {
		handleError(w, "Not Found", 404)
	}
}

func AssignRoutes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Printf("%s %s %s ", helpers.GetFormattedTime(time.Now()), r.Method, path)

	switch r.Method {
	case "GET":
		doGET(path, w, r)
	case "POST":
		doPOST(path, w, r)
	default:
		handleError(w, "Bad Request Method", 400)
	}
}
