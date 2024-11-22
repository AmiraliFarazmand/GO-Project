package routes

import (
	"net/http"
	"proj/internal/app/models"
	"proj/internal/app/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDLRoutes(router *gin.RouterGroup, db *gorm.DB) {
	router.POST("/", func(c *gin.Context) { createDL(c, db) })
	router.GET("/", func(c *gin.Context) {c.JSON(200, gin.H{
		"message":"hello",
	}) })
	router.GET("/:id", func(c *gin.Context) { getDL(c, db) })
	router.PUT("/", func(c *gin.Context) { updateDL(c, db) })
	router.DELETE("/", func(c *gin.Context) { deleteDL(c, db) })
}

func createDL(c *gin.Context, db *gorm.DB) {
	var dl models.DL
	if err := c.ShouldBindJSON(&dl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateDL(dl, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "DL created successfully"})
}
func getDL(c *gin.Context, db *gorm.DB) {
    id,_ := strconv.Atoi(c.Param("id"))
    dl, err := services.GetDL(id, db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, dl)
}

func updateDL(c *gin.Context, db *gorm.DB) {
	var dl models.DL
	if err := c.ShouldBindJSON(&dl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateDL(dl, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DL updated successfully"})
}

func deleteDL(c *gin.Context, db *gorm.DB) {
	var dl struct {
		ID      int     `json:"id"`
		Version float64 `json:"version"`
	}
	if err := c.ShouldBindJSON(&dl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.DeleteDL(dl.ID, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DL deleted successfully"})
}
