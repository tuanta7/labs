package table

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Order struct {
	ID        int `gorm:"primary_key"`
	UserID    int
	ProductID int
	CreatedAt time.Time
}

func (o *Order) TableName() string {
	return "orders"
}

func FakeOrder(id int) *Order {
	return &Order{
		ID:        id,
		UserID:    gofakeit.IntRange(1, 1000),
		ProductID: gofakeit.IntRange(1, 100),
		CreatedAt: gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
	}
}

func FakeOrders(count int) []*Order {
	orders := make([]*Order, count)
	for i := 0; i < count; i++ {
		orders[i] = FakeOrder(i + 1)
	}
	return orders
}
