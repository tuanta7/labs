package table

import (
	"generator/faker"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Product struct {
	ID          uint      `gorm:"primary_key"`
	Name        string    `gorm:"size:255;unique;not null"`
	Description string    `gorm:"type:text"`
	Currency    string    `gorm:"size:255;not null"`
	Price       float64   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (Product) TableName() string {
	return "products"
}

func FakeProduct(id uint) *Product {
	fakeProduct := gofakeit.Product()

	return &Product{
		ID:          id,
		Name:        fakeProduct.Name,
		Description: fakeProduct.Description,
		Currency:    gofakeit.Currency().Short,
		Price:       fakeProduct.Price,
		CreatedAt:   faker.PastDate(),
		UpdatedAt:   time.Now(),
	}
}

func FakeProducts(count uint) []*Product {
	products := make([]*Product, count)
	for i := uint(0); i < count; i++ {
		products[i] = FakeProduct(i + 1)
	}

	return products
}
