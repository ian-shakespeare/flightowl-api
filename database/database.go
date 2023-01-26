package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/arcticstorm9/flightowl-api/helpers"
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
	Admin      int64
}

const file string = "flightowl.db"

func connectToDB() *sql.DB {
	conn, err := sql.Open("sqlite3", file)

	if err != nil {
		panic("could not connect to database")
	}

	return conn
}

func Init() error {
	conn := connectToDB()
	defer conn.Close()

	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			sex TEXT,
			date_joined TEXT NOT NULL,
			admin INTEGER DEFAULT 0 NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	return nil
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
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Sex, &user.DateJoined, &user.Admin)
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
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Sex, &user.DateJoined, &user.Admin)
		if err != nil {
			return User{}, errors.New("not found")
		}
	}

	if user.Email != email {
		return User{}, errors.New("not found")
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
		return 0, errors.New("conflict")
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic("could not get id of inserted user")
	}

	return id, nil
}

func DeleteTestUser() error {
	conn := connectToDB()
	defer conn.Close()

	_, err := conn.Exec(`
		DELETE FROM users
		WHERE email = 'test@email.com';
	`)
	if err != nil {
		return err
	}

	return nil
}
