package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"proj/internal/app/models"
	"proj/internal/app/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSLRoutes(router *gin.RouterGroup, db *gorm.DB) {
	router.POST("/", func(c *gin.Context) { createSL(c, db) })
	router.GET("/:id", func(c *gin.Context) { getSL(c, db) })
	router.PUT("/", func(c *gin.Context) { updateSL(c, db) })
	router.PATCH("/", func(c *gin.Context) { updateSL(c, db) })
	router.DELETE("/", func(c *gin.Context) { deleteSL(c, db) })
}

func createSL(c *gin.Context, db *gorm.DB) {
	var sl models.SL

	rawData, _ := c.GetRawData()
	// Rebind raw data into the SL struct
	json.Unmarshal(rawData, &sl)

	// Check if "hasDL" exists in the JSON
	var rawMap map[string]interface{}
	json.Unmarshal(rawData, &rawMap)
	c.ShouldBindJSON(&sl)
	
	if _, includeHasDL := rawMap["hasDL"]; !includeHasDL {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hasDL field was not defined"})
		return
	}
	
	sl.ID = 0
	sl.Version = 1

	if err := services.CreateSL(sl, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "SL created successfully"})
}

func getSL(c *gin.Context, db *gorm.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	sl, err := services.GetSL(id, db)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "record not found"})
		return
	}
	c.JSON(http.StatusOK, sl)
}

func updateSL(c *gin.Context, db *gorm.DB) {
	var sl models.SL
	if err := c.ShouldBindJSON(&sl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateSL(sl, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SL updated successfully"})
}

func deleteSL(c *gin.Context, db *gorm.DB) {
	var sl models.SL
	if err := c.ShouldBindJSON(&sl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.DeleteSL(sl.ID, sl.Version, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SL deleted successfully"})
}
