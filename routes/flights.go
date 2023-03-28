package routes

import (
	"encoding/json"

	"flightowl-api/database"

	"github.com/gin-gonic/gin"
)

func GetFlights(c *gin.Context) {
	type request struct {
		OriginCode      string `json:"originLocationCode"`
		DestinationCode string `json:"destinationLocationCode"`
		DepartureDate   string `json:"departureDate"`
		NumOfAdults     int    `json:"adults"`
	}

	body, err := c.GetRawData()
	if err != nil {
		c.Status(400)
		return
	}

	var req request
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.Status(400)
		return
	}

	offers, err := getFlightOffers(
		req.OriginCode, req.DestinationCode,
		req.DepartureDate, req.NumOfAdults)
	if err != nil {
		switch err.Error() {
		case "bad request":
			c.Status(400)
			return
		case "not found":
			c.Status(400)
			return
		default:
			panic(err)
		}
	}

	c.JSON(200, offers)
}

func SaveFlight(c *gin.Context) {
	id := c.MustGet("UserId").(int64)

	body, err := c.GetRawData()
	if err != nil {
		c.Status(400)
		return
	}

	err = database.InsertFlightOffer(id, string(body))
	if err != nil {
		panic("could not insert flight offer")
	}

	c.Status(201)
}

func GetSavedFlights(c *gin.Context) {
	id := c.MustGet("UserId").(int64)

	offers, err := database.SelectFlightOffers(id)
	if err != nil {
		c.Status(400)
		return
	}

	c.JSON(200, offers)
}

func CheckSavedFlight(c *gin.Context) {
	type request struct {
		OfferId int64 `json:"offerId"`
	}

	id := c.MustGet("UserId").(int64)

	body, err := c.GetRawData()
	if err != nil {
		c.Status(400)
		return
	}

	var req request
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.Status(400)
		return
	}

	previousOffer, err := database.SelectFlightOffer(req.OfferId, id)
	if err != nil {
		c.Status(400)
		return
	}

	currentOffer, err := getUpdatedFlightOffer(previousOffer)
	if err != nil {
		switch err.Error() {
		case "not found":
			c.Status(404)
		default:
			c.Status(400)
		}
		return
	}

	c.JSON(200, currentOffer)
}
