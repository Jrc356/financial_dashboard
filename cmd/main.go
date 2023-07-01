package main

import (
	"fmt"
	"log"

	"github.com/Jrc356/financial_dashboard/controllers"
	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

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
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Panicln(err)
	}

	db.AutoMigrate(
		&models.Asset{},
		&models.AssetValue{},
		&models.Liability{},
		&models.LiabilityValue{},
	)
}

func main() {
	router := gin.Default()
	controllers.NewAssetController(db, router)
	controllers.NewLiabilityController(db, router)
	router.Run()
}
