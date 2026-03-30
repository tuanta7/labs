package main

import (
	"generator/table"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	PostgresDSN = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
)

func main() {
	db, err := gorm.Open(postgres.Open(PostgresDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(&table.User{})
	db.CreateInBatches(table.FakeUsers(1000), 1000)
}
