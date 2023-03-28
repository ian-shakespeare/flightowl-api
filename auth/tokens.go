package auth

import (
	"errors"
	"flightowl-api/helpers"
	"time"

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

func CreateJWT(id int64) []byte {
	return []byte(generateJWT(id))
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
