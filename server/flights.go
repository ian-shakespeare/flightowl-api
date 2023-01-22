package server

import (
	"encoding/json"
	"io"
	"net/http"
)

type OffersRequest struct {
	OriginCode      string `json:"originLocationCode"`
	DestinationCode string `json:"destinationLocationCode"`
	DepartureDate   string `json:"departureDate"`
	NumOfAdults     int    `json:"adults"`
}

func getFlights(w http.ResponseWriter, r *http.Request) {
	_, err := loadSession(r)
	if err != nil {
		handleUnauthorized(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w)
		return
	}

	var offerFields OffersRequest
	err = json.Unmarshal(body, &offerFields)
	if err != nil {
		handleBadRequest(w)
		return
	}

	offers, err := getFlightOffers(
		offerFields.OriginCode, offerFields.DestinationCode,
		offerFields.DepartureDate, offerFields.NumOfAdults)
	if err != nil {
		switch err.Error() {
		case "bad request":
			handleBadRequest(w)
			return
		case "not found":
			handleNotFound(w)
			return
		default:
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offers)
	handleOK(w)
}
