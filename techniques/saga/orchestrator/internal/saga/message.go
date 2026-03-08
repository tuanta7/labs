package saga

type ActionType string

const (
	ActionExecute  ActionType = "EXECUTE"
	ActionRollback ActionType = "ROLLBACK"
)

type StepType string

const (
	StepReserve  StepType = "RESERVE"
	StepPayment  StepType = "PAYMENT"
	StepCharging StepType = "CHARGING"
)

type Command struct {
	SagaID    string     `json:"saga_id"`
	Step      StepType   `json:"step"`
	Action    ActionType `json:"action"`
	UserID    string     `json:"user_id"`
	StationID string     `json:"station_id"`
	Amount    float64    `json:"amount"`
}

type Response struct {
	SagaID  string   `json:"saga_id"`
	Step    StepType `json:"step"`
	Success bool     `json:"success"`
	Error   string   `json:"error,omitempty"`
}
