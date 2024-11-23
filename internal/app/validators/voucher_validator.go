package validators

import (
	"errors"
	"fmt"
	"proj/internal/app/models"

	"gorm.io/gorm"
)

func CheckUniquenessVoucher(vouch models.Voucher, db *gorm.DB) error {
	var existingVouch models.Voucher

	if err := db.Where("id<>? AND number = ?", vouch.ID, vouch.Number).First(&existingVouch).Error; err != nil {
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

func CheckBalance(vId int, db *gorm.DB) (bool, error) {
	var vItems []models.VoucherItem
	if err := db.Where("voucher_id = ?", vId).Find(&vItems).Error; err != nil {
		return false, err
	}
	var creditSum, debitSum float64
	for _, item := range vItems {
		creditSum += item.Credit
		debitSum += item.Debit
	}
	return creditSum == debitSum, nil
}

func CheckItemsNumber(vId int, db *gorm.DB) (bool, error) {
	var vItems []models.VoucherItem
	if err := db.Where("voucher_id = ?", vId).Find(&vItems).Error; err != nil {
		return false, err
	}
	return len(vItems) < 500, nil
}
