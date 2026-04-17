package main

import (
	"generator/table"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	PostgresDSN  = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	UserCount    = 1000
	ProductCount = 100
)

func main() {
	db, err := gorm.Open(postgres.Open(PostgresDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	_ = db.Migrator().DropTable(&table.Order{}, &table.User{}, &table.Product{})
	_ = db.AutoMigrate(&table.User{}, &table.Product{}, &table.Order{})

	db.CreateInBatches(table.FakeUsers(UserCount), 100)
	db.CreateInBatches(table.FakeProducts(ProductCount), 100)
	db.CreateInBatches(table.FakeOrders(100, UserCount, ProductCount), 100)
}
