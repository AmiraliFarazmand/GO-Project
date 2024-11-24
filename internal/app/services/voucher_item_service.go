package services

import (
	"errors"
	"fmt"
	"proj/internal/app/models"
	"proj/internal/app/validators"

	"gorm.io/gorm"
)

func CreateVoucherItem(voucherItem models.VoucherItem, db *gorm.DB) error {
	if err := validators.ValidateVoucherItem(voucherItem, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR CreateVoucherItem(%v), faild at validating", voucherItem), err)
	}

	if err := db.Create(&voucherItem).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR CreateVoucherItem(%v),failed at creating item", voucherItem), err)
	}

	return nil
}

func GetVoucherItem(id int, db *gorm.DB) (models.VoucherItem, error) {
	var voucherItem models.VoucherItem
	// if err := db.First(&voucherItem, id).Error; err != nil {
	if err := db.Preload("DL").Preload("Sl").Preload("Voucher").First(&voucherItem, id).Error; err != nil {
		return models.VoucherItem{}, errors.Join(fmt.Errorf(">ERR GetVoucherItem(%d)", id), err)
	}
	return voucherItem, nil
}

func UpdateVoucherItem(voucherItem models.VoucherItem, vID int, db *gorm.DB) error {
	var currentVoucherItem models.VoucherItem
	if err := db.First(&currentVoucherItem, voucherItem.ID).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateVoucherItem(%v), VI notfound", voucherItem), err)
	}
	if currentVoucherItem.VoucherID != vID {
		return errors.New("voucher item belongs to another voucher")
	}
	if err := validators.ValidateVoucherItem(voucherItem, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateVoucherItem(%v), failed at validation", voucherItem), err)
	}

	if err := db.Save(&voucherItem).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateVoucherItem(%v),failed to save updated VI", voucherItem), err)
	}

	return nil
}

func DeleteVoucherItem(id int, vID int, db *gorm.DB) error {
	var voucherItem models.VoucherItem
	if err := db.First(&voucherItem, id).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR DeleteVoucherItem(%v), VI not found", id), err)
	}

	if voucherItem.VoucherID != vID {
		return errors.New("voucher item belongs to another voucher")
	}
	if err := db.Delete(&voucherItem).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR DeleteVoucherItem(%v), failed to delete VI", id), err)
	}

	return nil
}
