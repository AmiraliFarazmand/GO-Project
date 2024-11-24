package test_services

import (
	// "log"
	"errors"
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

	db.Create(&models.DL{ID: 1, Code: "DL001", Title: "Deletable DL", Version: 1.0})

	t.Run("Delete DL Successfully", func(t *testing.T) {
		err := services.DeleteDL(1, 1.0, db)
		assert.NoError(t, err, "Deleting an existing DL should succeed")

		var dl models.DL
		result := db.First(&dl, 1)
		assert.Error(t, result.Error, "Deleted DL should no longer exist")
		assert.True(t, errors.Is(result.Error, gorm.ErrRecordNotFound), "Error should be record not found")
	})

	t.Run("Delete DL - Not Found", func(t *testing.T) {
		err := services.DeleteDL(999, 1.0, db)
		assert.Error(t, err, "Deleting a non-existent DL should fail")
		assert.Contains(t, err.Error(), "not found", "Error message should indicate DL not found")
	})

	t.Run("Delete DL - Has References", func(t *testing.T) {
		db.Create(&models.DL{ID: 2, Code: "DL002", Title: "Referenced DL", Version: 1.0})
		tmp := 2
		db.Create(&models.VoucherItem{DLID: &tmp, SLID: 1, Debit: 100, Credit: 0})

		err := services.DeleteDL(2, 1.0, db)
		assert.Error(t, err, "Deleting a DL with references should fail")
		assert.Contains(t, err.Error(), "referenced elsewhere", "Error message should indicate references")
	})
	t.Run("Delete DL - Versions not match", func(t *testing.T) {
		db.Create(&models.DL{ID: 3, Code: "DL003", Title: "Simple DL", Version: 1.3})

		err := services.DeleteDL(3, 1.4, db)
		assert.Error(t, err, "Deleting a DL with a version mismatch should fail")

		var dl models.DL
		db.First(&dl, 3)
		assert.Equal(t, 1.3, dl.Version, "Version of the DL should remain unchanged")
	})
}
