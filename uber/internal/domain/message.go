package domain

type DriverLocation struct {
	TripID   string   `json:"tripId"`
	DriverID string   `json:"driverId"`
	Location Location `json:"location"`
}

type TripConfirmation struct {
	PassengerID string   `json:"passengerId"`
	Location    Location `json:"location"`
}
