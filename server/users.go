package server

import (
	"encoding/json"
	"io"
	"net/http"

	"flightowl.app/api/database"
	"flightowl.app/api/helpers"
)

type Credentials struct {
	Email    string
	Password string
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	_, err := loadSession(r)
	if err != nil {
		handleUnauthorized(w)
		return
	}

	users, err := database.SelectAllUsers()
	if err != nil {
		handleNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	handleOK(w)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w)
		return
	}

	var user database.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		handleBadRequest(w)
		return
	}
	user.Password = helpers.HashString(user.Password)

	id, err := database.InsertUser(user.FirstName, user.LastName, user.Email, user.Password, user.Sex)
	if err != nil {
		panic(err)
	}

	sessionId := createSession(id)
	cookie := http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	http.SetCookie(w, &cookie)

	handleCreated(w)
}

func authenticateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w)
		return
	}

	var credentials Credentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		handleBadRequest(w)
		return
	}

	user, err := database.SelectUser(credentials.Email)
	if err != nil {
		handleNotFound(w)
		return
	}

	if helpers.HashString(credentials.Password) != user.Password {
		handleUnauthorized(w)
		return
	}

	sessionId := createSession(user.Id)
	cookie := http.Cookie{
		Name:  "sessionId",
		Value: sessionId,
	}
	http.SetCookie(w, &cookie)

	handleCreated(w)
}
