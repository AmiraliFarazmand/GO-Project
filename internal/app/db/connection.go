package db

import (
	"errors"
	// "fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=ferz2020 dbname=goBootCamp port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New(">ERR db.Connect(). Failed to connect to database")
	}
	return db, nil
}
