package main

import (
	"log"
	"proj/internal/app/db"
	"proj/internal/app/models"
	"proj/internal/app/validators"
	// "gorm.io/gorm"
	// "gorm.io/gorm"
)

func main() {
	// Initialize database connection
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Run migrations
	db.RunMigrations(database)
	dl := 1
	voucherItem := models.VoucherItem{
		SlID:      1,
		// DlID:      nil, // or &someDlID
		DlID:      &dl, // or &someDlID
		VoucherID: 1,
		Debit:     1020.0,
		Credit:    0,
	}
	
	if err := validators.ValidateVoucherItem(voucherItem, database); err != nil {
		log.Println("Validation error:", err)
	} else if err := database.Create(&voucherItem).Error; err != nil {
		log.Printf("Failed to create VI: %v\n", err)
	} else {
		log.Println("VI record created successfully!")
	}

	log.Printf("mirese be inja")
}