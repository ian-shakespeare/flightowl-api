package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"flightowl.app/api/helpers"
)

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

var routes = []route{
	{"GET", "/users", getAllUsers},
	{"POST", "/users", createUser},
	{"POST", "/sessions", authenticateUser},
}

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

func handleOK(w http.ResponseWriter) {
	fmt.Printf("- 200\n")
}

func handleCreated(w http.ResponseWriter) {
	fmt.Printf("- 201\n")
	w.WriteHeader(http.StatusCreated)
}

func writeCookie(w http.ResponseWriter, name string, value string) {
	cookie := http.Cookie{
		Name:  name,
		Value: value,
	}
	http.SetCookie(w, &cookie)
}

func readCookie(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(key)
	if err != nil {
		return "", errors.New("cookie not found")
	}
	return cookie.Value, nil
}

func filterRoutesByMethod(r []route, method string) ([]route, error) {
	var filteredRoutes = []route{}

	for i, value := range r {
		if value.method == method {
			filteredRoutes = append(filteredRoutes, r[i])
		}
	}
	if len(filteredRoutes) < 1 {
		return nil, errors.New("bad request")
	}

	return filteredRoutes, nil
}

func pathsMatch(target string, candidate string) bool {
	tarElements := strings.Split(target, "/")
	canElements := strings.Split(candidate, "/")

	for i, value := range tarElements {
		if !helpers.Includes(canElements[i], ":") && value != canElements[i] {
			return false
		}
	}

	return true
}

func getMatchingRoute(r *http.Request, path string) (route, error) {
	possibleRoutes, err := filterRoutesByMethod(routes, r.Method)
	if err != nil {
		return route{}, err
	}

	for _, value := range possibleRoutes {
		if pathsMatch(path, value.path) {
			return value, nil
		}
	}

	return route{}, errors.New("not found")
}

func AssignRoutes(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s %s ", helpers.GetFormattedTime(time.Now()), r.Method, r.URL.Path)

	currentRoute, err := getMatchingRoute(r, r.URL.Path)
	if err != nil {
		if err.Error() == "bad request" {
			handleBadRequest(w)
			return
		} else if err.Error() == "not found" {
			handleNotFound(w)
			return
		} else {
			panic("unexpected error")
		}
	}

	currentRoute.handler(w, r)
}
