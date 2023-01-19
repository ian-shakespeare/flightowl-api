package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"flightowl.app/api/database"
)

func usersGet(id int, w http.ResponseWriter, r *http.Request) {
	u, err := database.GetUser(id)
	if err != nil || u.Id != id {
		handleError(w, "User Not Found", 404)
	}
	fmt.Printf("- %d\n", 200)
	io.WriteString(w, u.FirstName)
}

func usersCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, "No Request Body", 400)
	}

	var user database.User
	readError := json.Unmarshal(body, &user)
	if readError != nil {
		handleError(w, "Bad Request Body", 400)
	}

	id, err := database.CreateUser(user.FirstName, user.LastName, user.Email, user.Password, user.Sex)
	if err != nil {
		handleError(w, "Could not commit user to database", 500)
		panic(err)
	}
	fmt.Printf("- %d\n", 200)
	io.WriteString(w, fmt.Sprintf("user created with %d", id))
}
