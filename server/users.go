package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"flightowl-api/database"
	"flightowl-api/helpers"
	"flightowl-api/types"
)

type Credentials struct {
	Email    string
	Password string
}

type SessionInfo struct {
	SessionId string `json:"sessionId"`
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

	var user types.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		handleBadRequest(w)
		return
	}
	user.Password = helpers.HashString(user.Password)

	id, err := database.InsertUser(user.FirstName, user.LastName, user.Email, user.Password, user.Sex)
	if err != nil {
		switch err.Error() {
		case "conflict":
			handleConflict(w)
			return
		default:
			panic(err)
		}
	}

	sessionId := createSession(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SessionInfo{SessionId: sessionId})
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
		fmt.Println(err)
		handleNotFound(w)
		return
	}

	if helpers.HashString(credentials.Password) != user.Password {
		handleUnauthorized(w)
		return
	}

	sessionId := createSession(user.UserId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SessionInfo{SessionId: sessionId})
}
