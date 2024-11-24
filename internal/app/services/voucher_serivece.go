package services

import (
	"errors"
	"fmt"
	"proj/internal/app/models"
	"proj/internal/app/validators"

	"gorm.io/gorm"
)

func CreateVoucher(voucher models.Voucher, db *gorm.DB) (models.Voucher, error) {
	emptyVoucher := models.Voucher{}
	if err := validators.CheckUniquenessVoucher(voucher, db); err != nil {
		return emptyVoucher, errors.Join(fmt.Errorf(">ERR CreateVoucher(%v), faild at checking uniqueness",
			voucher), err)
	}
	if err := validators.ValidateVoucher(voucher, db); err != nil {
		return emptyVoucher, errors.Join(fmt.Errorf(">ERR CreateVoucher(%v), faild at validating",
			voucher), err)
	}

	if err := db.Create(&voucher).Error; err != nil {
		return emptyVoucher, errors.Join(fmt.Errorf(">ERR CreateVoucher(%v), faild at creating instance",
			voucher), err)
	}

	return voucher, nil
}

func GetVoucher(id int, db *gorm.DB) (models.Voucher, error) {
	var voucher models.Voucher
	if err := db.First(&voucher, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Voucher{},
				fmt.Errorf(">ERR GetVoucher(%d), voucher don't exist", id)
		}
		return models.Voucher{},
			errors.Join(fmt.Errorf(">ERR GetVoucher(%d), Uncommon error", id), err)
	}
	return voucher, nil
}

func UpdateVoucher(voucher models.Voucher, db *gorm.DB) error {
	if err := validators.ValidateVoucher(voucher, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateVoucher(%v),failed at validating", voucher), err)
	}

	var currentVoucher models.Voucher
	if err := db.First(&currentVoucher, voucher.ID).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateVoucher(%v), voucher not found to update", voucher), err)
	}

	if currentVoucher.Version != voucher.Version {
		return fmt.Errorf(">ERR UpdateVoucher(%v), versions don't match", voucher)
	}

	voucher.Version += 0.01
	if err := db.Save(&voucher).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateVoucher(%v), failed saving updated VOucher", voucher), err)
	}

	return nil
}

func DeleteVoucher(id int, version float64, db *gorm.DB) error {
	var v models.Voucher
	if err := db.First(&v, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(">ERR DeleteSL(%d), SL don't exist", id)
		}
	}

	var vItems []models.VoucherItem
	if err := db.Where("voucher_id = ?", id).Find(&vItems).Error; err != nil {
		return err
	}
	for _, item := range vItems {
		DeleteVoucherItem(item.ID, db)
	}

	if v.Version != version {
		return errors.New("SL.version is not what it should be")
	}

	if err := db.Delete(&models.Voucher{}, id).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateVoucher(%d), failed to delete voucher", id), err)
	}

	return nil
}
