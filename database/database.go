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

func connectToDB() *sql.DB {
	dbURL := helpers.GetRequiredEnv("DB_URL")

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("could not connect to database")
		panic(err)
	}

	return conn
}

func Init() error {
	conn := connectToDB()
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

	rows, err := conn.Query("SELECT * FROM users WHERE email = $1", email)
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
	conn := connectToDB()
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
	conn := connectToDB()
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

func SelectFlightOffers(user_id int64) ([]types.StoredOffer, error) {
	conn := connectToDB()
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM flight_offers WHERE user_id = $1;", user_id)
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
		"SELECT * FROM flight_offers WHERE offer_id = $1 AND user_id = $2;",
		offer_id, user_id)
	if err != nil {
		return types.Offer{}, err
	}
	defer rows.Close()

	storedOffer := types.StoredOffer{}
	if rows.Next() {
		var rawOfferData string
		err = rows.Scan(&storedOffer.OfferId, &storedOffer.DateSaved, &rawOfferData, &storedOffer.UserId)
		if err != nil {
			return types.Offer{}, err
		}

		var offerData types.FlightOfferData
		json.Unmarshal([]byte(rawOfferData), &offerData)
		storedOffer.FlightOffer = offerData
	}

	return storedOffer.FlightOffer.Data, nil
}

func DeleteTestFlight() error {
	conn := connectToDB()
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
