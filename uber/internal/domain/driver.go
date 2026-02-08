package domain

type Driver struct {
	ID          string  `json:"id" bson:"_id"`
	Name        string  `json:"name" bson:"name"`
	Phone       string  `json:"phone" bson:"phone"`
	Email       string  `json:"email" bson:"email"`
	RatingSum   float64 `json:"ratingSum" bson:"ratingSum"`
	RatingCount int     `json:"ratingCount" bson:"ratingCount"`
}

type Rating struct {
	TripID      string  `json:"tripId"`
	DriverID    string  `json:"driverId"`
	PassengerID string  `json:"passengerId"`
	Rating      float64 `json:"rating"`
	Comment     string  `json:"comment"`
}
