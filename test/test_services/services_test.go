package test_services

import (
	// "log"
	"proj/internal/app/models"
	"proj/internal/app/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.DL{}, &models.SL{}, &models.Voucher{}, &models.VoucherItem{})
	return db
}


func TestGetDL(t *testing.T) {
	db := setupTestDB()

	// Seed the database
	db.Create(&models.DL{ID: 1, Code: "DL001", Title: "Test DL"})

	t.Run("Get DL Successfully", func(t *testing.T) {
		dl, err := services.GetDL(1, db)
		assert.NoError(t, err, "Fetching an existing DL should succeed")
		assert.Equal(t, "DL001", dl.Code, "Fetched DL should have correct code")
	})

	t.Run("Get DL - Not Found", func(t *testing.T) {
		_, err := services.GetDL(999, db)
		assert.Error(t, err, "Fetching a non-existent DL should fail")
	})
}

func TestUpdateDL(t *testing.T) {
	db := setupTestDB()

	// Seed the database
	db.Create(&models.DL{ID: 1, Code: "DL001", Title: "Original DL", Version: 1.0})

	t.Run("Update DL Successfully", func(t *testing.T) {
		dl := models.DL{ID: 1, Code: "DL002", Title: "Updated DL", Version: 1.0}
		err := services.UpdateDL(dl, db)
		assert.NoError(t, err, "Updating a valid DL should succeed")

		var updatedDL models.DL
		db.First(&updatedDL, 1)
		assert.Equal(t, "DL002", updatedDL.Code, "Updated DL should have the correct code")
		assert.Equal(t, 1.01, updatedDL.Version, "Updated DL version should be incremented")
	})

	t.Run("Update DL - Version Mismatch", func(t *testing.T) {
		dl := models.DL{ID: 1, Code: "DL002", Title: "Invalid Update", Version: 2.0}
		err := services.UpdateDL(dl, db)
		assert.Error(t, err, "Updating a DL with version mismatch should fail")
	})

	t.Run("Update DL - Not Found", func(t *testing.T) {
		dl := models.DL{ID: 999, Code: "Nonexistent", Title: "Invalid DL"}
		err := services.UpdateDL(dl, db)
		assert.Error(t, err, "Updating a non-existent DL should fail")
	})
}

func TestDeleteDL(t *testing.T) {
	db := setupTestDB()

	// Seed the database
	db.Create(&models.DL{ID: 1, Code: "DL001", Title: "Deletable DL"})

	// t.Run("Delete DL Successfully", func(t *testing.T) {
	// 	err := services.DeleteDL(1, db)
	// 	assert.NoError(t, err, "Deleting an existing DL should succeed")

	// 	var dl models.DL
	// 	result := db.First(&dl, 1)
	// 	assert.Error(t, result.Error, "Deleted DL should no longer exist")
	// })

	// t.Run("Delete DL - Not Found", func(t *testing.T) {
	// 	err := services.DeleteDL(999, db)
	// 	log.Printf(">>>>>%v",err)
	// 	assert.Error(t, err, "Deleting a non-existent DL should fail")
	// })

	// t.Run("Delete DL - Has References", func(t *testing.T) {
	// 	dlIdTemp :=1
	// 	db.Create(&models.VoucherItem{DLID: &dlIdTemp})
	// 	err := services.DeleteDL(1, db)
	// 	assert.Error(t, err, "Deleting a DL with references should fail")
	// })
}
