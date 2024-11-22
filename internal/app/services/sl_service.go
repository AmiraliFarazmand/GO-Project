package services

import (
	"errors"
	"fmt"
	"proj/internal/app/models"
	"proj/internal/app/validators"
	"gorm.io/gorm"
)

func CreateSL(sl models.SL, db *gorm.DB) error {
	if err := validators.ValidateSL(sl, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR CreateSL(%v), faild at validating", sl), err)
	}
	
	if err := validators.CheckUniquenessSL(sl, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR CreateSL(%v), faild at checking uniqueness", sl), err)
	}

	if err := db.Create(&sl).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR CreateSL(%v), faild creating SL instance", sl), err)
	}

	return nil
}


func GetSL(id int, db *gorm.DB) (models.SL, error) {
	var sl models.SL
	if err := db.First(&sl, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SL{},
				fmt.Errorf(">ERR GetSL(%d), SL don't exist", id)
		}
		return models.SL{},
			errors.Join(fmt.Errorf(">ERR GetSL(%d), Uncommon error", id), err)
	}
	return sl, nil	
}


func UpdateSL(sl models.SL, db *gorm.DB) error {
	if err := validators.ValidateSL(sl, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateSL(%v), failed at validating", sl), err)
	}
	var currentSL models.SL
	if err := db.First(&currentSL, sl.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Join(fmt.Errorf(">ERR UpdateSL(%v),SL instance not found", sl), err)
		}
		return err
	}

	if currentSL.Version != sl.Version {
		return errors.New(">ERR UpdateSL(%v), versions are not same")
	}

	if err := validators.CheckUniquenessSL(sl, db); err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateSL(%v), failed because of checking uniqueness", sl), err)
	}

	sl.Version += 0.01

	if err := db.Save(&sl).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR UpdateSL(%v), failed saving updated SL", sl), err)
	}

	return nil
}



func DeleteSL(id int, db *gorm.DB) error {
	if hasReferencesSL(id, db) {
		return fmt.Errorf(">ERR DeleteSL(%d), it has reference somewhere",id)
	}

	if err := db.Delete(&models.SL{}, id).Error; err != nil {
		return errors.Join(fmt.Errorf(">ERR DeleteSL(%d), Error on deleting instance",id),err)
	}

	return nil
}

func hasReferencesSL(id int, db *gorm.DB) bool {
	var count int64
	db.Model(&models.VoucherItem{}).Where("sl_id = ?", id).Count(&count)
	return count > 0
}
