package validators

import (
	"errors"
	"proj/internal/app/models"

	"gorm.io/gorm"
)

func checkUniquenessDL(dl models.Dl, db *gorm.DB) error {
	var existingDL models.Dl

	// Check if Code or Title already exists in the database
	if err := db.Where("code = ? OR title = ?", dl.Code, dl.Title).First(&existingDL).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No duplicate found
			return nil
		}
		// Query error
		return err
	}

	// Duplicate found
	return errors.New("duplicate code or title found for DL")
}

func ValidateDL(dl models.Dl, db *gorm.DB) error {
	if dl.Code == "" {
		return errors.New("code cannot be empty")
	}
	if len(dl.Code) > 64 {
		return errors.New("code cannot exceed 64 characters")
	}
	if dl.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(dl.Title) > 64 {
		return errors.New("title cannot exceed 64 characters")
	}
	if err := checkUniquenessDL(dl, db); err != nil {
		return err
	}
	return nil
}
