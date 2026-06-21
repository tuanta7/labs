package table

import (
	"generator/faker"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type User struct {
	ID        uint `gorm:"primary_key"`
	Name      string
	Status    string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (u *User) TableName() string {
	return "users"
}

func FakeUser(id uint) *User {
	return &User{
		ID:        id,
		Name:      gofakeit.Name(),
		Status:    gofakeit.RandomString([]string{"active", "inactive", "pending"}),
		CreatedAt: faker.PastDate(),
		UpdatedAt: time.Now(),
	}
}

func FakeUsers(count uint) []*User {
	users := make([]*User, count)
	for i := uint(0); i < count; i++ {
		users[i] = FakeUser(i + 1)
	}
	return users
}
