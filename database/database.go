package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/arcticstorm9/flightowl-api/helpers"
	"github.com/arcticstorm9/flightowl-api/types"
	_ "github.com/mattn/go-sqlite3"
)

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
		CREATE TABLE IF NOT EXISTS users
		(
			user_id INTEGER PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			sex TEXT,
			date_joined TEXT NOT NULL,
			admin INTEGER DEFAULT 0 NOT NULL
		);
		CREATE TABLE IF NOT EXISTS flight_offers
		(
			offer_id INTEGER PRIMARY KEY,
			date_saved TEXT NOT NULL,
			offer TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func SelectAllUsers() ([]types.User, error) {
	conn := connectToDB()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []types.User{}
	for rows.Next() {
		user := types.User{}
		err = rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Sex, &user.DateJoined, &user.Admin)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func SelectUser(email string) (types.User, error) {
	conn := connectToDB()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM users WHERE email = ?;", email)
	if err != nil {
		return types.User{}, err
	}
	defer rows.Close()

	user := types.User{}
	if rows.Next() {
		err = rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Sex, &user.DateJoined, &user.Admin)
		if err != nil {
			return types.User{}, errors.New("not found")
		}
	}

	if user.Email != email {
		return types.User{}, errors.New("not found")
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

func InsertFlightOffer(user_id int64, body string) error {
	currentTime := helpers.GetFormattedTime(time.Now())
	conn := connectToDB()
	defer conn.Close()

	_, err := conn.Exec(`
		INSERT INTO flight_offers (date_saved, offer, user_id)
		VALUES(?, ?, ?);
	`, currentTime, body, user_id)
	if err != nil {
		return err
	}

	return nil
}

func SelectFlightOffers(user_id int64) ([]types.StoredOffer, error) {
	conn := connectToDB()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM flight_offers WHERE user_id = ?;", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offers := []types.StoredOffer{}
	for rows.Next() {
		var offer types.StoredOffer
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

func SelectFlightOffer(offer_id int64, user_id int64) (types.Offer, error) {
	conn := connectToDB()
	defer conn.Close()

	rows, err := conn.Query(
		"SELECT * FROM flight_offers WHERE offer_id = ? AND user_id = ?;",
		offer_id, user_id)
	if err != nil {
		return types.Offer{}, err
	}
	defer rows.Close()

	storedOffer := types.StoredOffer{}
	if rows.Next() {
		var rawOfferData string
		err = rows.Scan(&storedOffer.OfferId, &storedOffer.DateSaved, &rawOfferData, &storedOffer.UserId)

		var offerData types.FlightOfferData
		json.Unmarshal([]byte(rawOfferData), &offerData)
		storedOffer.FlightOffer = offerData
	}

	return storedOffer.FlightOffer.Data, nil
}
