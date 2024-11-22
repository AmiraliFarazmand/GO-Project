package validators

import (
	"errors"
	"fmt"
	"proj/internal/app/models"
	"gorm.io/gorm"
)

func CheckUniquenessVoucher(vouch models.Voucher, db *gorm.DB) error {
	var existingVouch models.Voucher

	if err := db.Where("number = ?", vouch.Number).First(&existingVouch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return errors.Join(fmt.Errorf(">ERR CheckUniquenessVoucher(v), ORM level error"), err)
	}
	return errors.New(">ERR CheckUniquenessVoucher(v), duplicate number for Voucher")
}

func ValidateVoucher(voucher models.Voucher, db *gorm.DB) error {
	if voucher.Number == "" {
		return errors.New(">ERR ValidateVoucher(v), number cannot be empty")
	}
	if len(voucher.Number) > 64 {
		return errors.New(">ERR ValidateVoucher(v), number cannot exceed 64 characters")
	}

	return nil
}
