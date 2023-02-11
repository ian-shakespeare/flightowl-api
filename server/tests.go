package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"flightowl-api/database"
	"flightowl-api/helpers"
)

type TestCredentials struct {
	Token string `json:"token"`
}

const TESTING_KEY = "\xba\x85)\xbb\x8d|\xf2B\x80z{\xb1\xb8\xab\u00a2\xda\a\x860\xde\xddY\x1a\xe6\xf7\x17\vRA\x99\xd4"

func checkToken(token string) bool {
	return helpers.HashString(token) == strconv.QuoteToASCII(TESTING_KEY)
}

func deleteTestUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w)
		return
	}

	var cred TestCredentials
	err = json.Unmarshal(body, &cred)
	if err != nil {
		handleBadRequest(w)
		return
	}

	if !checkToken(cred.Token) {
		handleUnauthorized(w)
		return
	}

	err = database.DeleteTestUser()
	if err != nil {
		panic("could not delete test user")
	}

	handleOK(w)
}

func deleteTestFlight(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w)
		return
	}

	var cred TestCredentials
	err = json.Unmarshal(body, &cred)
	if err != nil {
		handleBadRequest(w)
		return
	}

	if !checkToken(cred.Token) {
		handleUnauthorized(w)
		return
	}

	err = database.DeleteTestFlight()
	if err != nil {
		panic("could not delete test flight")
	}

	handleOK(w)
}
