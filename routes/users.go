package routes

import (
	"encoding/json"
	"fmt"

	"flightowl-api/auth"
	"flightowl-api/database"
	"flightowl-api/helpers"

	"github.com/gin-gonic/gin"
)

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

func GetUser(c *gin.Context) {
	id := c.MustGet("UserId").(int64)

	user, err := database.SelectUser(id)
	if err != nil {
		c.Status(404)
		return
	}

	var safeUser SafeUserData
	safeUser.FirstName = user.FirstName
	safeUser.LastName = user.LastName
	safeUser.Email = user.Email
	safeUser.Sex = user.Sex
	safeUser.DateJoined = user.DateJoined

	c.JSON(200, safeUser)
}

func RegisterUser(c *gin.Context) {
	type request struct {
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Sex        string `json:"sex"`
	}
	type response struct {
		Token string `json:"token"`
	}

	body, err := c.GetRawData()
	if err != nil {
		c.Status(400)
		return
	}

	var req request
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.Status(400)
		return
	}
	req.Password = helpers.HashString(req.Password)

	id, err := database.InsertUser(req.FirstName, req.LastName, req.Email, req.Password, req.Sex)
	if err != nil {
		switch err.Error() {
		case "conflict":
			c.Status(409)
			return
		default:
			panic(err)
		}
	}

	token := auth.CreateJWT(id)

	var res response
	res.Token = string(token)
	c.JSON(201, res)
}

func LoginUser(c *gin.Context) {
	type request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Token string `json:"token"`
	}

	body, err := c.GetRawData()
	if err != nil {
		c.Status(400)
		return
	}

	var req request
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.Status(400)
		return
	}

	user, err := database.SelectUserByEmail(req.Email)
	if err != nil {
		fmt.Println(err)
		c.Status(404)
		return
	}

	if helpers.HashString(req.Password) != user.Password {
		c.Status(401)
		return
	}

	token := auth.CreateJWT(user.UserId)

	var res response
	res.Token = string(token)
	c.JSON(201, res)
}
