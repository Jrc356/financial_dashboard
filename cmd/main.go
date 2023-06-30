package main

import (
	"fmt"
	"log"

	"github.com/Jrc356/financial_dashboard/controllers"
	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		"postgres",
		"postgres",
		"10.0.0.202",
		"5432",
		"postgres",
	)

	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}

	db.AutoMigrate(
		&models.Asset{},
	)
}

func main() {
	router := gin.Default()
	assetsController := controllers.AssetController{DB: db}
	assetsRouter := router.Group("/asset")
	assetsController.CreateRoutes(assetsRouter)
	router.Run()
}
