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
	response := gin.H{}
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
	if err := services.CreateVoucher(voucher, tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var currVoucher models.Voucher
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
			response[tempkey] = err.Error()
		}
	}
	// check number of items
	if check, err := validators.CheckItemsNumber(currVoucher.ID, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to error on CheckItemsNumber"})
		return
	} else if !check {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to invalid number of items"})
		return
	}

	// Check if voucherr is balanced
	if isBalanced, err := validators.CheckBalance(currVoucher.ID, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to an error on CheckBalance()"})
		return
	} else if !isBalanced {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to not being balanced"})
		return
	}
	//Commit transaction
	if err := tx.Commit().Error; err != nil {
		response["Failed at commiting transaction"] = "0"
	}

	if len(response) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Voucher created successfully"})
	}

	c.JSON(http.StatusInternalServerError,
		gin.H{"message": "Voucher created successfully with some errors",
			"errors": response})
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

	// Call service
	voucher, err := services.GetVoucher(id, db)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, voucher)
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
	response := gin.H{}
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

	//Transaction begind
	tx := db.Begin()
	// Call service
	if err := services.UpdateVoucher(voucher, tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, item := range inserted {
		if err := services.CreateVoucherItem(item, tx); err != nil {
			tempKey := fmt.Sprintf("insertItem(%v)", item)
			response[tempKey] = err.Error()
		}
	}
	for _, item := range updated {
		if err := services.UpdateVoucherItem(item, tx); err != nil {
			tempKey := fmt.Sprintf("updateItem(%v)", item)
			response[tempKey] = err.Error()
		}
	}
	for _, item := range request.Items.Deleted {
		if err := services.DeleteVoucherItem(item, tx); err != nil {
			tempKey := fmt.Sprintf("deleteItem(%v)", item)
			response[tempKey] = err.Error()
		}
	}
	// check number of items
	if check, err := validators.CheckItemsNumber(request.Voucher.ID, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to error on CheckItemsNumber"})
		return
	} else if !check {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to invalid number of items"})
		return
	}

	// Check if voucherr is balanced
	if isBalanced, err := validators.CheckBalance(request.Voucher.ID, tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to an error on CheckBalance()"})
		return
	} else if !isBalanced {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ROLLBACK occured due to not being balanced"})
		return
	}
	//Commit transaction
	if err := tx.Commit().Error; err != nil {
		response["Failed at commiting transaction"] = "0"
	}

	if len(response) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Voucher updated successfully"})
	}

	c.JSON(http.StatusInternalServerError,
		gin.H{"message": "Voucher updated successfully with some errors",
			"errors": response})
}
