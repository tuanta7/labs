package main

type Location struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"-"`
	UserID    string  `json:"userId"`
	TripID    string  `json:"tripId"`
}
