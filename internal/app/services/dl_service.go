package services

import (
	"errors"
	"fmt"
	"proj/internal/app/models"
	"proj/internal/app/validators"
	"gorm.io/gorm"
)

func CreateDL(dl models.DL, db *gorm.DB) error {
	if err := validators.CheckUniquenessDL(dl, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR CreateDL(%v), faild at checking uniqueness",dl),err)
	}
	if err := validators.ValidateDL(dl, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR CreateDL(%v), failed at validation",dl),err)
	}

	if err := db.Create(&dl).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR CreateDL(%v), failed creating instance",dl),err)
	}

	return nil
}

func GetDL(id int, db *gorm.DB) (models.DL, error) {
	var dl models.DL
	if err := db.First(&dl, id).Error; err != nil {
		return models.DL{}, errors.Join(fmt.Errorf(">ERR GetDL(%v), not found", dl), err)
	}
	return dl, nil
}

func UpdateDL(dl models.DL, db *gorm.DB) error {
	if err := validators.ValidateDL(dl, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateDl(%v), failed at validating", dl), err)
	}
	var currentDL models.DL
	if err := db.First(&currentDL, dl.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(">ERR UpdateDl(%v), DL not found", dl)
		}
		return err
	}

	// check Version
	if currentDL.Version != dl.Version {
		return fmt.Errorf(">ERR UpdateDl(%v), versions are not  same", dl)
	}

	if err := validators.CheckUniquenessDL(dl, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateDl(%v), failed due to uniqueness violation", dl), err)
	}
	dl.Version += 0.01

	if err := db.Save(&dl).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateDl(%v),failed saving updated DL", dl), err)
	}

	return nil
}

func DeleteDL(id int, db *gorm.DB) error {
	if hasReferences(id, db) {
		return fmt.Errorf(">ERR DeleteDL(%d, referenced elsewhere", id)
	}

	if err := db.Delete(&models.DL{}, id).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR DeleteDL(%d, referenced elsewhere", id), err)
	}

	return nil
}

func hasReferences(id int, db *gorm.DB) bool {
	var count int64
	db.Model(&models.VoucherItem{}).Where("dl_id = ?", id).Count(&count)
	return count > 0
}
