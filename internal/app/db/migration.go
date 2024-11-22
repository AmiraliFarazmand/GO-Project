package db

import (
	"log"
	// "proj/internal/app/models"
	"proj/internal/app/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	if err := db.AutoMigrate(
		&models.DL{},
		&models.SL{},
		&models.Voucher{},
		&models.VoucherItem{},
	); err != nil {
		log.Fatalf(">ERR db.RunMigraitons(). Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully!")
}
