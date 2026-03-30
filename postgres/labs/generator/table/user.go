package table

import "github.com/brianvoe/gofakeit/v7"

type User struct {
	ID   int `gorm:"primary_key"`
	Name string
}

func (u *User) TableName() string {
	return "users"
}

func FakeUser(id int) *User {
	return &User{
		ID:   id,
		Name: gofakeit.Name(),
	}
}

func FakeUsers(count int) []*User {
	users := make([]*User, count)
	for i := 0; i < count; i++ {
		users[i] = FakeUser(i + 1)
	}
	return users
}
