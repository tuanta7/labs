package table

import (
	"generator/faker"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Order struct {
	ID        uint `gorm:"primary_key"`
	UserID    int
	ProductID int
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	User      User      `gorm:"foreignkey:UserID"`
	Product   Product   `gorm:"foreignkey:ProductID"`
}

func (o *Order) TableName() string {
	return "orders"
}

func FakeOrder(id, maxUserID, maxProductID uint) *Order {
	return &Order{
		ID:        id,
		UserID:    gofakeit.IntRange(1, int(maxUserID)),
		ProductID: gofakeit.IntRange(1, int(maxProductID)),
		CreatedAt: faker.PastDate(),
		UpdatedAt: time.Now(),
	}
}

func FakeOrders(count, maxUserID, maxProductID uint) []*Order {
	orders := make([]*Order, count)
	for i := uint(0); i < count; i++ {
		orders[i] = FakeOrder(i+1, maxUserID, maxProductID)
	}
	return orders
}
