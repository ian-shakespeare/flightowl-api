package database

import (
	"database/sql"
	"time"

	"flightowl.app/api/helpers"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id         int64
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
		panic("could not connect to database")
	}

	return conn
}

func SelectAllUsers() ([]User, error) {
	conn := connectToDB()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Sex, &user.DateJoined)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func SelectUser(email string) (User, error) {
	conn := connectToDB()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM users WHERE email = ?;", email)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	user := User{}
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Sex, &user.DateJoined)
		if err != nil {
			panic("could not get user from database")
		}
	}

	return user, nil
}

func InsertUser(firstName string, lastName string, email string, password string, sex string) (int64, error) {
	currentTime := helpers.GetFormattedTime(time.Now())
	conn := connectToDB()
	defer conn.Close()

	res, err := conn.Exec(`
		INSERT INTO users (first_name, last_name, email, password, sex, date_joined)
		VALUES(?, ?, ?, ?, ?, ?);
		`, firstName, lastName, email, password, sex, currentTime)
	if err != nil {
		panic("could not insert user into databse")
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic("could not get id of inserted user")
	}

	return id, nil
}
