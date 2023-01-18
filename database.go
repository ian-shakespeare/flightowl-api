package main

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id         int
	FirstName  string
	LastName   string
	Email      string
	Password   string
	Sex        string
	DateJoined string
}

const file string = "flightowl.db"

func connectToDB() *sql.DB {
	conn, err := sql.Open("sqlite3", file)
	if err != nil {
		panic("Error: Could not connect to database")
	}
	return conn
}

func getUser(id int) (User, error) {
	conn := connectToDB()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM users WHERE id = ?;", id)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	user := User{}

	if rows.Next() {
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Sex, &user.DateJoined)
		if err != nil {
			return User{}, err
		}
	}

	return user, nil
}

func createUser(firstName string, lastName string, email string, password string, sex string) (int64, error) {
	currentTime := strings.Split(time.Now().String(), " +")[0]
	conn := connectToDB()
	defer conn.Close()

	res, err := conn.Exec("INSERT INTO users VALUES('?', '?', '?', '?', '?', '?')", firstName, lastName, email, password, sex, currentTime)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
