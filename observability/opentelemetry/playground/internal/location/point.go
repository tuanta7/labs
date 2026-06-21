package location

type Driver struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
}

type Location struct {
	DriverID    string  `json:"driverId"`
	IsAvailable bool    `json:"isAvailable"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
