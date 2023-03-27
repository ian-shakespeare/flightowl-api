package main

import (
	"flightowl-api/auth"
	"flightowl-api/database"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func testAuth(c *gin.Context) {
	userId := c.MustGet("UserId").(int64)
	c.JSON(200, userId)
}

func RequireAuth() gin.HandlerFunc {
	return func (c *gin.Context) {
		header, exists := c.Request.Header["Authorization"]
		if !exists {
			c.Status(401)
			return
		}

		tokenString := strings.Split(header[0], "Bearer ")[1]
		token, err := auth.DecodeJWT(tokenString)
		if err != nil {
			c.Status(401)
			return
		}

		c.Set("UserId", token.UserId)
	}
}

func main() {
	database.Initialize()

	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(RequireAuth())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://www.flightowl.app", "https://flightowl.app", "http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,	
	}))

	r.POST("/tokens", auth.GetToken)

	authorized := r.Group("/")
	authorized.Use(RequireAuth())
	{
		authorized.GET("/test", testAuth)
	}

	r.Run(":8000")
}
