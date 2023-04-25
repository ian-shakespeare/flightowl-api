package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"flightowl-api/auth"
	"flightowl-api/database"
	"flightowl-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func (c *gin.Context) {
		header, exists := c.Request.Header["Authorization"]
		if !exists || len(header) == 0 {
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

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://www.flightowl.app", "https://flightowl.app", "http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,	
	}))

	r.POST("/register", routes.RegisterUser)
	r.POST("/login", routes.LoginUser)

	authorized := r.Group("/")
	authorized.Use(RequireAuth())
	{
		authorized.GET("/user", routes.GetUser)
		authorized.GET("/flights/saved", routes.GetSavedFlights)
		authorized.POST("/flights", routes.GetFlights)
		authorized.POST("/flights/check", routes.CheckSavedFlight)
		authorized.POST("/flights/saved", routes.SaveFlight)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r.Run(fmt.Sprintf(":%s", port))
}
