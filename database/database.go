package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"flightowl-api/helpers"
	"flightowl-api/types"

	_ "github.com/lib/pq"
)

type User struct {
	UserId     int64
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Sex        string `json:"sex"`
	DateJoined string
	Admin      int64
}

type SavedOffer struct {
	OfferId     int64           `json:"offer_id"`
	DateSaved   string          `json:"date_saved"`
	FlightOffer types.FlightOfferData `json:"offer"`
	UserId      int64           `json:"user_id"`
}

func getConn() *sql.DB {
	dbURL := helpers.GetRequiredEnv("DATABASE_URL")

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	return conn
}

func Initialize() error {
	conn := getConn()
	defer conn.Close()

	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS users
		(
			id SERIAL PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			sex TEXT NOT NULL,
			date_joined TEXT NOT NULL,
			admin INTEGER DEFAULT 0 NOT NULL
		);
		CREATE TABLE IF NOT EXISTS flight_offers
		(
			offer_id SERIAL PRIMARY KEY,
			date_saved TEXT NOT NULL,
			offer TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
				ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func SelectUser(id int64) (User, error) {
	conn := getConn()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM users WHERE id = $1;", id)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	user := User{}
	if rows.Next() {
		err = rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Sex, &user.DateJoined, &user.Admin)
		if err != nil {
			return User{}, err
		}
	}

	return user, nil
}

func SelectUserByEmail(email string) (User, error) {
	conn := getConn()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	user := User{}
	if rows.Next() {
		err = rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Sex, &user.DateJoined, &user.Admin)
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
	conn := getConn()
	defer conn.Close()

	var userId int64
	err := conn.QueryRow(`
		INSERT INTO users (first_name, last_name, email, password, sex, date_joined)
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id
	`, firstName, lastName, email, password, sex, currentTime).Scan(&userId)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("conflict")
	}

	return userId, nil
}

func DeleteTestUser() error {
	conn := getConn()
	defer conn.Close()

	_, err := conn.Exec(`
		DELETE FROM users
		WHERE email = 'test@email.com'
	`)
	if err != nil {
		return err
	}

	return nil
}

func InsertFlightOffer(user_id int64, body string) error {
	currentTime := helpers.GetFormattedTime(time.Now())
	conn := getConn()
	defer conn.Close()

	_, err := conn.Exec(`
		INSERT INTO flight_offers (date_saved, offer, user_id)
		VALUES($1, $2, $3);
	`, currentTime, body, user_id)
	if err != nil {
		return err
	}

	return nil
}

func SelectFlightOffers(user_id int64) ([]SavedOffer, error) {
	conn := getConn()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM flight_offers WHERE user_id = $1;", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offers := []SavedOffer{}
	for rows.Next() {
		var offer SavedOffer
		var rawOfferData string
		err = rows.Scan(&offer.OfferId, &offer.DateSaved, &rawOfferData, &offer.UserId)
		if err != nil {
			return nil, err
		}

		var offerData types.FlightOfferData
		json.Unmarshal([]byte(rawOfferData), &offerData)
		offer.FlightOffer = offerData
		offers = append(offers, offer)
	}

	return offers, nil
}

func SelectFlightOffer(offer_id int64, user_id int64) (types.FlightOffer, error) {
	conn := getConn()
	defer conn.Close()

	rows, err := conn.Query(
		"SELECT * FROM flight_offers WHERE offer_id = $1 AND user_id = $2;",
		offer_id, user_id)
	if err != nil {
		return types.FlightOffer{}, err
	}
	defer rows.Close()

	storedOffer := SavedOffer{}
	if rows.Next() {
		var rawOfferData string
		err = rows.Scan(&storedOffer.OfferId, &storedOffer.DateSaved, &rawOfferData, &storedOffer.UserId)
		if err != nil {
			return types.FlightOffer{}, err
		}

		var offerData types.FlightOfferData
		json.Unmarshal([]byte(rawOfferData), &offerData)
		storedOffer.FlightOffer = offerData
	}

	return storedOffer.FlightOffer.Data, nil
}

func DeleteTestFlight() error {
	conn := getConn()
	defer conn.Close()

	_, err := conn.Exec(`
		DELETE FROM flights
		JOIN users ON flights.user_id = users.id
		WHERE email = 'test@email.com';
	`)
	if err != nil {
		return err
	}

	return nil
}
