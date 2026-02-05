package domain

type Driver struct {
	ID             string  `json:"id" bson:"_id"`
	Name           string  `json:"name" bson:"name"`
	Phone          string  `json:"phone" bson:"phone"`
	Email          string  `json:"email" bson:"email"`
	Rating         float64 `json:"rating" bson:"rating"`
	Score          float64 `json:"score" bson:"score"`
	CompletedTrips int     `json:"completedTrips" bson:"completedTrips"`
}

type Rating struct {
	TripID      string  `json:"tripId"`
	DriverID    string  `json:"driverId"`
	PassengerID string  `json:"passengerId"`
	Rating      float64 `json:"rating"`
	Comment     string  `json:"comment"`
}

type Behavior struct {
	PreferredSpeedBand []int
	MicroRouteBiases   map[string]any
	StopPattern        map[string]any
}
