package main

import (
	"log"
	"proj/internal/routes"
	"proj/internal/app/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var database *gorm.DB // Assume initialized as described previously

func main() {
	var err error
	// Initialize database connection
	database, err = db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Run migrations
	db.RunMigrations(database)	
	router := gin.Default()
	dlRoutes := router.Group("/dl")
	routes.RegisterDLRoutes(dlRoutes, database)
	slRoutes := router.Group("/sl")
	routes.RegisterSLRoutes(slRoutes, database)
	vRoutes := router.Group("/voucher")
	routes.RegisterVoucherRoutes(vRoutes, database)

	router.Run(":8080")
}
