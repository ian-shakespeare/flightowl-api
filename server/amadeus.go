package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// const apiKey = ""
// const apiSecret = ""
const baseURL = "https://test.api.amadeus.com"

var accessToken = AccessToken{"", 0, 0}

func isTokenExpired(token *AccessToken) bool {
	return time.Now().Unix() > token.duration+token.timeReceived
}

func getAccessToken() error {
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

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var rawResBody AuthResponse
	err = json.Unmarshal(resBody, &rawResBody)
	if err != nil {
		return err
	}
	accessToken.token = rawResBody.AccessToken
	accessToken.duration = rawResBody.ExpiresIn
	accessToken.timeReceived = time.Now().Unix()

	return nil
}

func GetFlightOffers(originCode string, destinationCode string, departureDate string, numOfAdults int) error {
	if isTokenExpired(&accessToken) {
		getAccessToken()
	}

	resourseURL := "/v2/shopping/flight-offers"
	queryParameters := fmt.Sprintf("?originLocationCode=%s&destinationLocationCode%s&departureDate=%s&adults=%d&currencyCode=USD")
	requestURL := baseURL + resourseURL + queryParameters

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", accessToken.token)
}
