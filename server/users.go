package server

import (
	"encoding/json"
	"io"
	"net/http"

	"flightowl.app/api/database"
	"flightowl.app/api/helpers"
)

type AuthFields struct {
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

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	id := helpers.GetPathResource(r.URL.Path)

// 	user, err := database.SelectUser(id)
// 	if err != nil {
// 		handleError(w, "user not found", 404)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(user)
// 	logSuccess(200)
// }

func createUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w)
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
	writeCookie(w, "sessionId", sessionId)

	handleCreated(w)
}

func authenticateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w)
		return
	}

	var credentials AuthFields
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
	writeCookie(w, "sessionId", sessionId)

	handleCreated(w)
}
