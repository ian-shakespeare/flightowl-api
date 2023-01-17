package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	id          int
	first_name  string
	last_name   string
	email       string
	password    string
	sex         string
	date_joined string
}

const file string = "flightowl.db"

func getUserByEmail(email string) (User, error) {
	conn, err := sql.Open("sqlite3", file)
	if err != nil {
		panic("Error connecting to the database")
	}
	defer conn.Close()

	rows, err := conn.Query(fmt.Sprintf("SELECT * FROM users WHERE email = '%s';", email))
	if err != nil {
		panic("Error accessing row")
	}
	defer rows.Close()

	user := User{}

	if rows.Next() {
		err = rows.Scan(&user.id, &user.first_name, &user.last_name, &user.email, &user.password, &user.sex, &user.date_joined)
		if err != nil {
			panic("Could not find user")
		}
	}
	fmt.Printf("%s\n", user.first_name)
	return user, err
}
