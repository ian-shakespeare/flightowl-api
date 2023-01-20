package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"flightowl.app/api/helpers"
)

var sessionStore = map[string]string{}

func createSessionId() string {
	return helpers.GetRandomString(16)
}

func loadSession(r *http.Request) (int64, error) {
	if len(sessionStore) < 1 {
		return 0, errors.New("session store empty")
	}

	sessionId, err := r.Cookie("sessionId")
	if err != nil {
		return 0, err
	}

	sessionData := sessionStore[sessionId.Value]
	if sessionData == "" {
		return 0, errors.New("could not find sessionId")
	}

	id, err := strconv.ParseInt(sessionData, 10, 64)
	if err != nil {
		fmt.Printf("id: '%s'\n", sessionData)
		panic("error converting session data")
	}

	return id, nil
}

func createSession(id int64) string {
	sessionId := createSessionId()
	sessionStore[sessionId] = fmt.Sprint(id)
	return sessionId
}
