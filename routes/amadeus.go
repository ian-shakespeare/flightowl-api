package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"flightowl-api/types"
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
	Value        string
	Duration     int64
	TimeReceived int64
}


const baseURL = "https://test.api.amadeus.com"

var apiKey = os.Getenv("API_KEY")
var apiSecret = os.Getenv("API_SECRET")
var accessToken = AccessToken{Value: "", Duration: 0, TimeReceived: 0}

func isTokenExpired(token *AccessToken) bool {
	return time.Now().Unix() > token.Duration+token.TimeReceived
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
	token.Value = resBody.AccessToken
	token.Duration = resBody.ExpiresIn
	token.TimeReceived = time.Now().Unix()

	return nil
}

func getFlightOffers(originCode string, destinationCode string, departureDate string, numOfAdults int) ([]types.FlightOffer, error) {
	if isTokenExpired(&accessToken) {
		retrieveAccessToken(&accessToken)
	}

	resourceURL := "/v2/shopping/flight-offers"
	queryParameters := fmt.Sprintf(
		"?originLocationCode=%s&destinationLocationCode=%s&departureDate=%s&adults=%d&currencyCode=USD",
		originCode, destinationCode, departureDate, numOfAdults)
	requestURL := baseURL + resourceURL + queryParameters

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, errors.New("bad request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken.Value))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("bad request")
	}

	rawResBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("bad request")
	}

	var resBody types.FlightOffersResponse
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

func getUpdatedFlightOffer(previousOffer types.FlightOffer) (types.FlightOffer, error) {
	if isTokenExpired(&accessToken) {
		retrieveAccessToken(&accessToken)
	}

	resourceURL := "/v1/shopping/flight-offers/pricing"
	requestURL := baseURL + resourceURL

	var body types.FlightOfferPrice
	body.Data.Type = "flight-offers-pricing"
	body.Data.FlightOffers = append(body.Data.FlightOffers, previousOffer)

	reqBody, err := json.Marshal(body)
	if err != nil {
		return types.FlightOffer{}, errors.New("bad request")
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return types.FlightOffer{}, errors.New("bad request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken.Value))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.FlightOffer{}, errors.New("bad request")
	}

	rawResBody, err := io.ReadAll(res.Body)
	if err != nil {
		return types.FlightOffer{}, errors.New("bad request")
	}

	var resBody types.FlightOfferPrice
	err = json.Unmarshal(rawResBody, &resBody)
	if err != nil {
		return types.FlightOffer{}, errors.New("bad request")
	}

	if len(resBody.Data.FlightOffers) < 1 {
		return types.FlightOffer{}, errors.New("not found")
	}

	return resBody.Data.FlightOffers[0], nil
}
