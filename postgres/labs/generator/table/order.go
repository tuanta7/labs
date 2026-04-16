package table

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Order struct {
	ID        uint `gorm:"primary_key"`
	UserID    int
	ProductID int
	CreatedAt time.Time
}

func (o *Order) TableName() string {
	return "orders"
}

func FakeOrder(id, maxUserID uint) *Order {
	return &Order{
		ID:        id,
		UserID:    gofakeit.IntRange(1, int(maxUserID)),
		ProductID: gofakeit.IntRange(1, 100),
		CreatedAt: gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
	}
}

func FakeOrders(count, maxUserID uint) []*Order {
	orders := make([]*Order, count)
	for i := uint(0); i < count; i++ {
		orders[i] = FakeOrder(i+1, maxUserID)
	}
	return orders
}
