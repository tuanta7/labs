package faker

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func PastDate() time.Time {
	return gofakeit.DateRange(
		time.Now().AddDate(-3, 0, 0),
		time.Now().AddDate(-1, 0, 0),
	)
}
