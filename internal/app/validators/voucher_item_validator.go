package validators

import (
	"errors"
	"proj/internal/app/models"
	"gorm.io/gorm"
)

func ValidateVoucherItem(voucherItem models.VoucherItem, db *gorm.DB) error {
	var sl models.SL
	if err := db.First(&sl, voucherItem.SLID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(">ERR ValidateVoucherItem(vi, dn). SLID does not exist")
		}
		return errors.Join(errors.New(">ERR ValidateVoucherItem(vi, dn). orther ORM errors"),err)
	}

	if !sl.HasDL {
		if voucherItem.DLID != nil {
			return errors.New(">ERR ValidateVoucherItem(vi, dn), SL does not support DL references, DLID must be nil")
		}
	} else {
		// If HasDL is true, DLID must reference an existing DL
		if voucherItem.DLID == nil {
			return errors.New(">ERR ValidateVoucherItem(vi, dn), SL requires a valid DL reference, DLID cannot be nil")
		}
		var dl models.DL
		if err := db.First(&dl, *voucherItem.DLID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New(">ERR ValidateVoucherItem(vi, dn), DL with the given ID does not exist")
			}
			return err 
		}
	}
	if err := validateCreditDebit(voucherItem.Credit, voucherItem.Debit); err != nil {
		return errors.Join(errors.New(">ERR ValidateVoucherItem(vi, dn)"),err)
	}
	return nil
}

func validateCreditDebit(credit, debit float64) error {
	if debit < 0 || credit < 0 {
		return errors.New(">ERR validateCreditDebit(c,d), credit or Debit cannot be negative numbers")
	}
	if debit > 0 && credit > 0 {
		return errors.New(">ERR validateCreditDebit(c,d), both debit and credit cannot be greater than 0 at the same time")
	}
	if debit == 0 && credit == 0 {
		return errors.New(">ERR validateCreditDebit(c,d), either debit or credit must be greater than 0")
	}
	return nil
}
