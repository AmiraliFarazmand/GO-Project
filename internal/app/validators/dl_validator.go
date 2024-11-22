package validators

import (
	"errors"
	"fmt"
	"proj/internal/app/models"
	"gorm.io/gorm"
)

func CheckUniquenessDL(dl models.DL, db *gorm.DB) error {
	var existingDL models.DL

	if err := db.Where("id <>? AND (code = ? OR title = ?) ",dl.ID, dl.Code, dl.Title).First(&existingDL).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return errors.Join(fmt.Errorf(">ERR CheckUniquenessDL(dl)" ), err)
	}

	return fmt.Errorf(">ERR CheckUniquenessDL(dl), Duplicated items inside DL")
}


func ValidateDL(dl models.DL, db *gorm.DB) error {
	if dl.Code == "" {
		return errors.New(">ERR ValidateDL(dl), code cannot be empty")
	}
	if len(dl.Code) > 64 {
		return errors.New(">ERR ValidateDL(dl), code cannot exceed 64 characters")
	}
	if dl.Title == "" {
		return errors.New(">ERR ValidateDL(dl), title cannot be empty")
	}
	if len(dl.Title) > 64 {
		return errors.New(">ERR ValidateDL(dl), title cannot exceed 64 characters")
	}
	return nil
}
