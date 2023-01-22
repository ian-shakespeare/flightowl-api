package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type AccessToken struct {
	token        string
	duration     int64
	timeReceived int64
}

type LocationDate struct {
	IATACode string `json:"iataCode"`
	Terminal string `json:"terminal"`
	At       string `json:"at"`
}

type Segment struct {
	Departure    LocationDate `json:"departure"`
	Arrival      LocationDate `json:"arrival"`
	CarrierCode  string       `json:"carrierCode"`
	FlightNumber string       `json:"number"`
	OperatedBy   struct {
		CarrierCode string `json:"CarrierCode"`
	} `json:"operating"`
}

type Itinerary struct {
	Duration string    `json:"duration"`
	Segments []Segment `json:"segments"`
}

type Offer struct {
	OneWay            bool        `json:"oneWay"`
	LastTicketingDate string      `json:"lastTicketingDate"`
	AvailableSeats    int         `json:"numberOfBookableSeats"`
	Itineraries       []Itinerary `json:"itineraries"`
}

type FlightOffersResponse struct {
	Offers []Offer `json:"data"`
}

const baseURL = "https://test.api.amadeus.com"

var apiKey = os.Getenv("API_KEY")
var apiSecret = os.Getenv("API_SECRET")
var accessToken = AccessToken{"", 0, 0}

func isTokenExpired(token *AccessToken) bool {
	return time.Now().Unix() > token.duration+token.timeReceived
}

func retrieveAccessToken(token *AccessToken) error {
	requestURL := baseURL + "/v1/security/oauth2/token"

	formBody := url.Values{}
	formBody.Set("grant_type", "client_credentials")
	formBody.Set("client_id", apiKey)
	formBody.Set("client_secret", apiSecret)
	body := strings.NewReader(formBody.Encode())

	req, err := http.NewRequest(http.MethodPost, requestURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	rawResBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var resBody AuthResponse
	err = json.Unmarshal(rawResBody, &resBody)
	if err != nil {
		return err
	}
	token.token = resBody.AccessToken
	token.duration = resBody.ExpiresIn
	token.timeReceived = time.Now().Unix()

	return nil
}

func getFlightOffers(originCode string, destinationCode string, departureDate string, numOfAdults int) ([]Offer, error) {
	if isTokenExpired(&accessToken) {
		retrieveAccessToken(&accessToken)
	}

	resourseURL := "/v2/shopping/flight-offers"
	queryParameters := fmt.Sprintf(
		"?originLocationCode=%s&destinationLocationCode=%s&departureDate=%s&adults=%d&currencyCode=USD",
		originCode, destinationCode, departureDate, numOfAdults)
	requestURL := baseURL + resourseURL + queryParameters

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, errors.New("bad request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken.token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("bad request")
	}

	rawResBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("bad request")
	}

	var resBody FlightOffersResponse
	err = json.Unmarshal(rawResBody, &resBody)
	if err != nil {
		return nil, errors.New("bad request")
	}

	offers := resBody.Offers
	if len(offers) < 1 {
		return nil, errors.New("not found")
	} else if len(offers) < 5 {
		return offers[:3], nil
	}
	return offers[:5], nil
}
