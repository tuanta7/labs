package domain

type Trip struct {
	ID          string      `json:"id" `
	DriverID    string      `json:"driverId"`
	PassengerID string      `json:"passengerId"`
	Tracks      []*Location `json:"tracks"`
}
