package server

import (
	"net/http"
)

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

var routes = []route{
	{"GET", "/user", getUser},
	{"GET", "/flights/saved", getSavedFlights},
	{"POST", "/flights", getFlights},
	{"POST", "/flights/check", checkSavedFlight},
	{"POST", "/users", createUser},
	{"POST", "/sessions", authenticateUser},
	{"POST", "/flights", saveFlight},
	{"DELETE", "/tests", deleteTestUser},
}
