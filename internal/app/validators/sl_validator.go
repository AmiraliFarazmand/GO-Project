package validators

import (
	"errors"
	"gorm.io/gorm"
	"proj/internal/app/models"
)

func checkUniquenessSL(sl models.SL, db *gorm.DB) error {
	var existingSl models.SL

	// Check if Code or Title already exists in the database
	if err := db.Where("code = ? OR title = ?", sl.Code, sl.Title).First(&existingSl).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No duplicate found
			return nil
		}
		// Query error
		return err
	}

	// Duplicate found
	return errors.New("duplicate code or title found for SL")
}


func ValidateSL(sl models.SL, db *gorm.DB) error {
	if sl.Code == "" {
		return errors.New("code cannot be empty")
	}
	if len(sl.Code) > 64 {
		return errors.New("code cannot exceed 64 characters")
	}
	if sl.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(sl.Title) > 64 {
		return errors.New("title cannot exceed 64 characters")
	}
	if err := checkUniquenessSL(sl,db); err!=nil{
		return err
	}
	return nil
}
