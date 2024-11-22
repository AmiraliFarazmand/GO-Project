package main

import (
	// "fmt"
	"fmt"
	"log"
	"proj/internal/app/db"
	"proj/internal/app/models"
	"proj/internal/app/services"
	// "proj/internal/app/validators"
)

func main() {
	// Initialize database connection
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Run migrations
	db.RunMigrations(database)
	dlID := 12
	// Create a VoucherItem
	voucherItem := models.VoucherItem{
		VoucherID: 5,
		SLID:      1, // You may need to adjust these based on your actual data
		DLID:      &dlID,
		Debit:     100.0,
		Credit:    0.0,
	}

	// Test creating a VoucherItem
	err = services.CreateVoucherItem(voucherItem, database)
	if err != nil {
		log.Printf("Error creating VoucherItem: %v", err)
	} else {
		log.Println("VoucherItem created successfully!")
	}

	// // Test retrieving the VoucherItem by ID (e.g., ID = 1)
	// retrievedVoucherItem, err := services.GetVoucherItem(10, database)
	// if err != nil {
	// 	log.Printf("Error retrieving VoucherItem: %v", err)
	// } else {
	// 	log.Printf("Retrieved VoucherItem: %+v", retrievedVoucherItem)
	// }

	// // Test updating the VoucherItem
	// voucherItemToUpdate := models.VoucherItem{
	// 	ID:        11, // Assuming the ID is 1; adjust if needed
	// 	VoucherID: 6,
	// 	SLID:      1,
	// 	DLID:      &dlID,
	// 	Debit:     120.0,
	// 	Credit:    4.0,
	// }

	// err = services.UpdateVoucherItem(voucherItemToUpdate, database)
	// if err != nil {
	// 	log.Printf("Error updating VoucherItem: %v", err)
	// } else {
	// 	log.Println("VoucherItem updated successfully!")
	// }

	// // Test retrieving the updated VoucherItem
	// updatedVoucherItem, err := services.GetVoucherItem(12, database)
	// if err != nil {
	// 	log.Printf("Error retrieving updated VoucherItem: %v", err)
	// } else {
	// 	log.Printf("Updated VoucherItem: %+v", updatedVoucherItem)
	// }

	// // Test deleting the VoucherItem
	// err = services.DeleteVoucherItem(12, database)
	// if err != nil {
	// 	log.Printf("Error deleting VoucherItem: %v", err)
	// } else {
	// 	log.Println("VoucherItem deleted successfully!")
	// }
	sl, err:= services.GetSL(1151,database)
	if err==nil{
		fmt.Printf("%v",sl)
	}else{
		fmt.Println(err)
	}
	log.Println("mirese be inja")
}
