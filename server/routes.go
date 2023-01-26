package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/arcticstorm9/flightowl-api/helpers"
)

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

var routes = []route{
	{"GET", "/users", getAllUsers},
	{"GET", "/flights", getFlights},
	{"POST", "/users", createUser},
	{"POST", "/sessions", authenticateUser},
	{"DELETE", "/tests", deleteTestUser},
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

	if len(tarElements) != len(canElements) {
		return false
	}

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
