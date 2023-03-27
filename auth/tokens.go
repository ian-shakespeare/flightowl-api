package auth

import (
	"encoding/json"
	"errors"
	"flightowl-api/database"
	"flightowl-api/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
}

var secretKey = []byte(helpers.GetRequiredEnv("JWT_SECRET"))

func generateJWT(userId int64) string {
	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "https://api.flightowl.app",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func getSignedToken(email string, password string) ([]byte, error) {
	user, err := database.SelectUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if helpers.HashString(password) != user.Password {
		return nil, errors.New("unauthorized")
	}

	return []byte(generateJWT(user.UserId)), nil
}

func DecodeJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func (token *jwt.Token) (interface {}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("unauthorized")
	}

	return claims, nil
}

func GetToken(c *gin.Context) {
	type request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Token string `json:"token"`
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(400, nil)
		return
	}

	var cred request
	err = json.Unmarshal(jsonData, &cred)
	if err != nil {
		c.Status(400)
		return
	}

	token, err := getSignedToken(cred.Email, cred.Password)
	if err != nil {
		switch (err.Error()) {
		case "unauthorized":
			c.Status(401)
			return
		default:
			panic(err)
		}
	}

	var res response
	res.Token = string(token)
	c.JSON(201, res)
}
