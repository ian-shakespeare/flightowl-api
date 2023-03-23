package server

import (
	"encoding/json"
	"io"
	"net/http"

	"flightowl-api/database"
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

func saveFlight(w http.ResponseWriter, r *http.Request) {
	id, err := loadSession(r)
	if err != nil {
		handleUnauthorized(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w)
		return
	}

	err = database.InsertFlightOffer(id, string(body))
	if err != nil {
		panic("could not insert flight offer")
	}

	handleCreated(w)
}

func getSavedFlights(w http.ResponseWriter, r *http.Request) {
	id, err := loadSession(r)
	if err != nil {
		handleUnauthorized(w)
		return
	}

	offers, err := database.SelectFlightOffers(id)
	if err != nil {
		handleBadRequest(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offers)

	handleOK(w)
}

func checkSavedFlight(w http.ResponseWriter, r *http.Request) {
	id, err := loadSession(r)
	if err != nil {
		handleUnauthorized(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleBadRequest(w)
		return
	}

	var bodyJSON struct {
		OfferId int64 `json:"offerId"`
	}
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		handleBadRequest(w)
		return
	}

	previousOffer, err := database.SelectFlightOffer(bodyJSON.OfferId, id)
	if err != nil {
		handleBadRequest(w)
		return
	}

	currentOffer, err := getUpdatedFlightOffer(previousOffer)
	if err != nil {
		switch err.Error() {
		case "not found":
			handleNotFound(w)
		default:
			handleBadRequest(w)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentOffer)

	handleOK(w)
}
