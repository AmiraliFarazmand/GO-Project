package validators

import (
	"errors"
	"gorm.io/gorm"
	"proj/internal/app/models"
)
func checkUniquenessVoucher(vouch models.Voucher, db *gorm.DB) error {
	var existingVouch models.Voucher

	// Check if Code or Title already exists in the database
	if err := db.Where("number = ?", vouch.Number).First(&existingVouch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No duplicate found
			return nil
		}
		// Query error
		return err
	}
	return errors.New("duplicate number for Voucher")
}
func ValidateVoucher(voucher models.Voucher, db *gorm.DB) error {
	if voucher.Number == "" {
		return errors.New("number cannot be empty")
	}
	if len(voucher.Number) > 64 {
		return errors.New("number cannot exceed 64 characters")
	}
	if err:= checkUniquenessVoucher(voucher, db); err!=nil{
		return err
	}
	return nil
}