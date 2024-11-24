package routes

import (
	"fmt"
	"net/http"
	"proj/internal/app/models"
	"proj/internal/app/services"
	"proj/internal/app/validators"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterVoucherRoutes(router *gin.RouterGroup, db *gorm.DB) {
	router.POST("/", func(c *gin.Context) { createVoucher(c, db) })
	router.PUT("/", func(c *gin.Context) { updateVoucher(c, db) })
	router.DELETE("/", func(c *gin.Context) { deleteVoucher(c, db) })
	router.GET("/:id", func(c *gin.Context) { getVoucher(c, db) })
}

func createVoucher(c *gin.Context, db *gorm.DB) {
	//Transaction begin
	tx := db.Begin()
	var request struct {
		Voucher struct {
			Number string `json:"number"`
		} `json:"voucher"`
		Items []struct {
			SLID   int     `json:"sl_id"`
			DLID   *int    `json:"dl_id"`
			Debit  float64 `json:"debit"`
			Credit float64 `json:"credit"`
		} `json:"items"`
	}
	errorsOnResponse := gin.H{}
	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Prepare Voucher model
	voucher := models.Voucher{
		Number:  request.Voucher.Number,
		Version: 1.0,
	}

	// Call service
	var currVoucher models.Voucher
	var err error
	if currVoucher, err = services.CreateVoucher(voucher, tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Last(&currVoucher)
	var voucherItems []models.VoucherItem
	for _, item := range request.Items {
		voucherItems = append(voucherItems, models.VoucherItem{
			VoucherID: currVoucher.ID,
			SLID:      item.SLID,
			DLID:      item.DLID,
			Debit:     item.Debit,
			Credit:    item.Credit,
		})
	}

	for _, item := range voucherItems {
		if err := services.CreateVoucherItem(item, tx); err != nil {
			tempkey := fmt.Sprintf("insertItem(%v)", item)
			errorsOnResponse[tempkey] = err.Error()
		}
	}

	// Check if voucherr is balanced and have <500 items
	if finalCheckRes, err := finalCheck(currVoucher.ID, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to an error on finalCheck()",
			"errors": errorsOnResponse})
		return
	} else if !finalCheckRes {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to not being balanced" +
			" or having more than 500 items",
			"errors": errorsOnResponse})
		return
	}

	//Commit transaction
	if err := tx.Commit().Error; err != nil {
		errorsOnResponse["Failed at commiting transaction"] = "0"
	}

	if len(errorsOnResponse) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Voucher created successfully"})
		return
	}

	c.JSON(http.StatusInternalServerError,
		gin.H{"message": "Voucher created successfully with some errors",
			"errors": errorsOnResponse})
}

func deleteVoucher(c *gin.Context, db *gorm.DB) {
	var request struct {
		ID      int     `json:"id"`
		Version float64 `json:"version"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Call service
	if err := services.DeleteVoucher(request.ID, request.Version, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voucher deleted successfully"})
}

func getVoucher(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Call service for voucher
	voucher, err := services.GetVoucher(id, db)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var vItems []models.VoucherItem
	var responseVItems []models.VoucherItem
	if err := db.Where("voucher_id = ?", id).Find(&vItems).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	for _, item := range vItems {
		tempVI, _ := services.GetVoucherItem(item.ID, db)
		responseVItems = append(responseVItems, tempVI)
	}
	c.JSON(http.StatusOK, gin.H{"voucher": voucher,
		"items": responseVItems})

}

func updateVoucher(c *gin.Context, db *gorm.DB) {
	var request struct {
		Voucher struct {
			ID      int     `json:"id"`
			Number  string  `json:"number"`
			Version float64 `json:"version"`
		} `json:"voucher"`
		Items struct {
			Inserted []struct {
				SLID   int     `json:"sl_id"`
				DLID   *int    `json:"dl_id"`
				Debit  float64 `json:"debit"`
				Credit float64 `json:"credit"`
			} `json:"inserted"`
			Updated []struct {
				ID     int     `json:"id"`
				SLID   int     `json:"sl_id"`
				DLID   *int    `json:"dl_id"`
				Debit  float64 `json:"debit"`
				Credit float64 `json:"credit"`
			} `json:"updated"`
			Deleted []int `json:"deleted"`
		} `json:"items"`
	}
	errorsOnResponse := gin.H{}
	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Prepare Voucher model
	voucher := models.Voucher{
		ID:      request.Voucher.ID,
		Number:  request.Voucher.Number,
		Version: request.Voucher.Version,
	}

	// Prepare items
	inserted := []models.VoucherItem{}
	for _, item := range request.Items.Inserted {
		inserted = append(inserted, models.VoucherItem{
			VoucherID: request.Voucher.ID,
			SLID:      item.SLID,
			DLID:      item.DLID,
			Debit:     item.Debit,
			Credit:    item.Credit,
		})
	}

	updated := []models.VoucherItem{}
	for _, item := range request.Items.Updated {
		updated = append(updated, models.VoucherItem{
			ID:        item.ID,
			VoucherID: request.Voucher.ID,
			SLID:      item.SLID,
			DLID:      item.DLID,
			Debit:     item.Debit,
			Credit:    item.Credit,
		})
	}

	//Transaction begins
	tx := db.Begin()
	// Call service for voucher
	if err := services.UpdateVoucher(voucher, tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// call service for items
	for _, item := range inserted {
		if err := services.CreateVoucherItem(item, tx); err != nil {
			tempKey := fmt.Sprintf("insertItem(%v)", item)
			errorsOnResponse[tempKey] = err.Error()
		}
	}
	for _, item := range updated {
		if err := services.UpdateVoucherItem(item, tx); err != nil {
			tempKey := fmt.Sprintf("updateItem(%v)", item)
			errorsOnResponse[tempKey] = err.Error()
		}
	}
	for _, item := range request.Items.Deleted {
		if err := services.DeleteVoucherItem(item, tx); err != nil {
			tempKey := fmt.Sprintf("deleteItem(%v)", item)
			errorsOnResponse[tempKey] = err.Error()
		}
	}

	// Check if voucherr is balanced and have <500 items
	if finalCheckRes, err := finalCheck(request.Voucher.ID, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to an error on finalCheck()",
			"errors": errorsOnResponse})
		return
	} else if !finalCheckRes {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to not being balanced" +
			" or having more than 500 items",
			"errors": errorsOnResponse})
		return
	}
	//Commit transaction
	if err := tx.Commit().Error; err != nil {
		errorsOnResponse["Failed at commiting transaction"] = "0"
	}

	if len(errorsOnResponse) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Voucher updated successfully"})
		return
	}

	c.JSON(http.StatusInternalServerError,
		gin.H{"message": "Voucher updated successfully with some errors",
			"errors": errorsOnResponse})
}

func finalCheck(vId int, tx *gorm.DB) (bool, error) {
	var isBalanced, check bool
	var err error
	if isBalanced, err = validators.CheckBalance(vId, tx); err != nil {
		return false, err
	}
	if check, err = validators.CheckItemsNumber(vId, tx); err != nil {
		return false, err
	}
	return isBalanced && check, err
}
