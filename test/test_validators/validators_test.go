package test_validators

import (
	"proj/internal/app/models"
	"proj/internal/app/validators"
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

func TestCheckUniquenessDL(t *testing.T) {
	db := setupTestDB()

	// Seed the database with a DL
	db.Create(&models.DL{Code: "DL001", Title: "Test DL", Version: 1.0})

	t.Run("Unique DL", func(t *testing.T) {
		dl := models.DL{Code: "DL002", Title: "Unique DL"}
		err := validators.CheckUniquenessDL(dl, db)
		assert.NoError(t, err, "Unique DL should not return an error")
	})

	t.Run("Duplicate DL", func(t *testing.T) {
		dl := models.DL{Code: "DL001", Title: "Test DL"}
		err := validators.CheckUniquenessDL(dl, db)
		assert.Error(t, err, "Duplicate DL should return an error")
	})
}

func TestValidateDL(t *testing.T) {
	db := setupTestDB()

	t.Run("Valid DL", func(t *testing.T) {
		dl := models.DL{Code: "DL001", Title: "Valid DL"}
		err := validators.ValidateDL(dl, db)
		assert.NoError(t, err, "Valid DL should not return an error")
	})

	t.Run("Invalid DL - Empty Code", func(t *testing.T) {
		dl := models.DL{Code: "", Title: "No Code"}
		err := validators.ValidateDL(dl, db)
		assert.Error(t, err, "DL with empty code should return an error")
	})

	t.Run("Invalid DL - Code Too Long", func(t *testing.T) {
		dl := models.DL{Code: "A very very very very very very very very very very long code "+
		"111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
		 Title: "Long Code"}
		err := validators.ValidateDL(dl, db)
		assert.Error(t, err, "DL with a code longer than 64 characters should return an error")
	})

	t.Run("Invalid DL - Empty Title", func(t *testing.T) {
		dl := models.DL{Code: "DL002", Title: ""}
		err := validators.ValidateDL(dl, db)
		assert.Error(t, err, "DL with empty title should return an error")
	})
}

func TestValidateVoucherItem(t *testing.T) {
	db := setupTestDB()

	// Seed the database
	db.Create(&models.SL{ID: 1, Code: "SL001", Title: "SL Test", HasDL: true})
	db.Create(&models.DL{ID: 1, Code: "DL001", Title: "DL Test"})
	db.Create(&models.SL{ID: 2, Code: "SL002", Title: "SL No DL", HasDL: false})

	t.Run("Valid VoucherItem", func(t *testing.T) {
		vi := models.VoucherItem{SLID: 1, DLID: new(int), Debit: 100, Credit: 0}
		*vi.DLID = 1
		err := validators.ValidateVoucherItem(vi, db)
		assert.NoError(t, err, "Valid VoucherItem should not return an error")
	})

	t.Run("Invalid VoucherItem - Nonexistent SL", func(t *testing.T) {
		vi := models.VoucherItem{SLID: 999, DLID: nil, Debit: 100, Credit: 0}
		err := validators.ValidateVoucherItem(vi, db)
		assert.Error(t, err, "VoucherItem with nonexistent SL should return an error")
	})

	t.Run("Invalid VoucherItem - Nonexistent DL", func(t *testing.T) {
		vi := models.VoucherItem{SLID: 1, DLID: new(int), Debit: 100, Credit: 0}
		*vi.DLID = 999
		err := validators.ValidateVoucherItem(vi, db)
		assert.Error(t, err, "VoucherItem with nonexistent DL should return an error")
	})

	t.Run("Invalid VoucherItem - Credit and Debit Both Positive", func(t *testing.T) {
		vi := models.VoucherItem{SLID: 1, DLID: new(int), Debit: 100, Credit: 100}
		*vi.DLID = 1
		err := validators.ValidateVoucherItem(vi, db)
		assert.Error(t, err, "VoucherItem with both credit and debit positive should return an error")
	})

	t.Run("Invalid VoucherItem - Credit and Debit Both Zero", func(t *testing.T) {
		vi := models.VoucherItem{SLID: 1, DLID: nil, Debit: 0, Credit: 0}
		err := validators.ValidateVoucherItem(vi, db)
		assert.Error(t, err, "VoucherItem with both credit and debit zero should return an error")
	})
}
