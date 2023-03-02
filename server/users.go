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

type SafeUserData struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Sex        string `json:"sex"`
	DateJoined string `json:"dateJoined"`
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := loadSession(r)
	if err != nil {
		handleUnauthorized(w)
		return
	}

	user, err := database.SelectUser(id)
	if err != nil {
		handleNotFound(w)
		return
	}

	var safeUser SafeUserData
	safeUser.FirstName = user.FirstName
	safeUser.LastName = user.LastName
	safeUser.Email = user.Email
	safeUser.Sex = user.Sex
	safeUser.DateJoined = user.DateJoined

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(safeUser)

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
	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    sessionId,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Domain:   ".flightowl.app",
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

	user, err := database.SelectUserByEmail(credentials.Email)
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
	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    sessionId,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Domain:   ".flightowl.app",
	}
	http.SetCookie(w, &cookie)

	handleCreated(w)
}
