package validators

import (
	"errors"
	"fmt"
	"proj/internal/app/models"
	"gorm.io/gorm"
)

func CheckUniquenessSL(sl models.SL, db *gorm.DB) error {
	var existingSL models.SL

	if err := db.Where("id<>? AND (code = ? OR title = ?)",sl.ID, sl.Code, sl.Title).First(&existingSL).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return errors.Join(fmt.Errorf(">ERR CheckUniquenessSL(sl)"), err)
	}
	return errors.New(">ERR CheckUniquenessSL(sl), duplicate code or title found for SL")
}

func ValidateSL(sl models.SL, db *gorm.DB) error {
	if sl.Code == "" {
		return errors.New(">ERR ValidateSL(sl), code cannot be empty")
	}
	if len(sl.Code) > 64 {
		return errors.New(">ERR ValidateSL(sl), code cannot exceed 64 characters")
	}
	if sl.Title == "" {
		return errors.New(">ERR ValidateSL(sl), title cannot be empty")
	}
	if len(sl.Title) > 64 {
		return errors.New(">ERR ValidateSL(sl), title cannot exceed 64 characters")
	}
	return nil
}
