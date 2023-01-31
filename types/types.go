package types

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type AccessToken struct {
	Value        string
	Duration     int64
	TimeReceived int64
}

type User struct {
	UserId     int64
	FirstName  string
	LastName   string
	Email      string
	Password   string
	Sex        string
	DateJoined string
	Admin      int64
}

type StoredOffer struct {
	OfferId     int64           `json:"offer_id"`
	DateSaved   string          `json:"date_saved"`
	FlightOffer FlightOfferData `json:"offer"`
	UserId      int64           `json:"user_id"`
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
	Aircraft     struct {
		Code string `json:"code"`
	} `json:"aircraft"`
	OperatedBy struct {
		CarrierCode string `json:"carrierCode"`
	} `json:"operating"`
	Duration    string `json:"duration"`
	Id          string `json:"id"`
	Stops       int    `json:"numberOfStops"`
	Blacklisted bool   `json:"blacklistedInEU"`
}

type Itinerary struct {
	Duration string    `json:"duration"`
	Segments []Segment `json:"segments"`
}

type OfferPrice struct {
	Currency string `json:"currency"`
	Total    string `json:"total"`
	Base     string `json:"base"`
	Fees     []struct {
		Amount string `json:"amount"`
		Type   string `json:"type"`
	} `json:"fees"`
	GrandTotal string `json:"grandTotal"`
}

type PriceOptions struct {
	FareType        []string `json:"fareType"`
	CheckedBagsOnly bool     `json:"includedCheckedBagsOnly"`
}

type SegmentDetail struct {
	SegmentId   string `json:"segmentId"`
	Cabin       string `json:"cabin"`
	FareBasis   string `json:"fareBasis"`
	BrandedFare string `json:"brandedFare"`
	Class       string `json:"class"`
	CheckedBags struct {
		Quantity int `json:"quantity"`
	} `json:"includedCheckedBags"`
}

type TravelerPrice struct {
	TravelerId   string `json:"travelerId"`
	FareOption   string `json:"fareOption"`
	TravelerType string `json:"travelerType"`
	Price        struct {
		Currency string `json:"currency"`
		Total    string `json:"total"`
		Base     string `json:"base"`
	} `json:"price"`
	DetailsBySegment []SegmentDetail `json:"fareDetailsBySegment"`
}

type Offer struct {
	Type                     string          `json:"type"`
	Id                       string          `json:"id"`
	Source                   string          `json:"source"`
	InstantTicketingRequired bool            `json:"instantTicketingRequired"`
	NonHomogeneous           bool            `json:"nonHomogeneous"`
	OneWay                   bool            `json:"oneWay"`
	LastTicketingDate        string          `json:"lastTicketingDate"`
	AvailableSeats           int             `json:"numberOfBookableSeats"`
	Itineraries              []Itinerary     `json:"itineraries"`
	Price                    OfferPrice      `json:"price"`
	PricingOptions           PriceOptions    `json:"pricingOptions"`
	ValidatingAirlineCodes   []string        `json:"validatingAirlineCodes"`
	TravelerPricings         []TravelerPrice `json:"travelerPricings"`
}

type FlightOffersResponse struct {
	Offers []Offer `json:"data"`
}

type FlightOfferData struct {
	Data Offer `json:"data"`
}

type FlightOfferPrice struct {
	Data struct {
		Type         string  `json:"type"`
		FlightOffers []Offer `json:"flightOffers"`
	} `json:"data"`
}
