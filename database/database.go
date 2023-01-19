package database

import (
	"database/sql"
	"time"

	"flightowl.app/api/helpers"
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

func GetUser(id int) (User, error) {
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

func CreateUser(firstName string, lastName string, email string, password string, sex string) (int64, error) {
	currentTime := helpers.GetFormattedTime(time.Now())
	conn := connectToDB()
	defer conn.Close()

	res, err := conn.Exec(`
		INSERT INTO users (first_name, last_name, email, password, sex, date_joined)
		VALUES(?, ?, ?, ?, ?, ?);
		`, firstName, lastName, email, password, sex, currentTime)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
