package validators

import (
	"errors"
	"proj/internal/app/models"
	"gorm.io/gorm"
)

func ValidateVoucherItem(voucherItem models.VoucherItem, db *gorm.DB) error {
	var sl models.SL
	if err := db.First(&sl, voucherItem.SlID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("SL with the given ID does not exist")
		}
		return err // Return other database errors
	}

	if !sl.HasDl {
		// If HasDl is false, DlID must be nil
		if voucherItem.DlID != nil {
			return errors.New("SL does not support DL references, DlID must be nil")
		}
	} else {
		// If HasDl is true, DlID must reference an existing DL
		if voucherItem.DlID == nil {
			return errors.New("SL requires a valid DL reference, DlID cannot be nil")
		}
		var dl models.Dl
		if err := db.First(&dl, *voucherItem.DlID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("DL with the given ID does not exist")
			}
			return err // Return other database errors
		}
	}
	if err := validateCreditDebit(voucherItem.Credit, voucherItem.Debit); err != nil {
		return err
	}
	return nil
}

func validateCreditDebit(credit, debit float64) error {
	if debit < 0 || credit < 0 {
		return errors.New("credit or Debit cannot be negative numbers")
	}
	if debit > 0 && credit > 0 {
		return errors.New("both debit and credit cannot be greater than 0 at the same time")
	}
	if debit == 0 && credit == 0 {
		return errors.New("either debit or credit must be greater than 0")
	}
	return nil
}
