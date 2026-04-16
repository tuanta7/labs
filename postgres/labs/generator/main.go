package main

import (
	"generator/table"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	PostgresDSN = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	UserCount   = 1000
)

func main() {
	db, err := gorm.Open(postgres.Open(PostgresDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	_ = db.Migrator().DropTable(&table.User{}, &table.Order{})
	_ = db.AutoMigrate(&table.User{}, &table.Order{})

	db.CreateInBatches(table.FakeUsers(UserCount), 100)
	db.CreateInBatches(table.FakeOrders(100, UserCount), 100)
}
