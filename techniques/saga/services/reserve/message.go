package main

type SagaCommand struct {
	SagaID    string  `json:"saga_id"`
	Step      string  `json:"step"`
	Action    string  `json:"action"`
	UserID    string  `json:"user_id"`
	StationID string  `json:"station_id"`
	Amount    float64 `json:"amount"`
}

type SagaResponse struct {
	SagaID  string `json:"saga_id"`
	Step    string `json:"step"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
