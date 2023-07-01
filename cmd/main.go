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

func createAssetControllers(router *gin.Engine) {
	assetsController := controllers.AssetController{DB: db}

	assetsRouter := router.Group("/asset")
	{
		assetsRouter.POST("", assetsController.CreateAsset)
		assetsRouter.GET("", assetsController.ListAssets)
		assetsRouter.GET("/values", assetsController.ListAllAssetValues)
	}

	assetRouter := assetsRouter.Group("/:assetName")
	{
		assetRouter.GET("", assetsController.GetAsset)
		assetRouter.PUT("", assetsController.UpdateAsset)

		assetRouter.GET("/value", assetsController.GetAssetValues)
		assetRouter.POST("/value", assetsController.CreateAssetValue)

		assetRouter.DELETE("", assetsController.DeleteAsset)
	}
}

func createLiabilityControllers(router *gin.Engine) {
	liabilitiesController := controllers.LiabilityController{DB: db}

	liabilitiesRouter := router.Group("/liability")
	{
		liabilitiesRouter.POST("", liabilitiesController.CreateLiability)
		liabilitiesRouter.GET("", liabilitiesController.ListLiabilities)
		liabilitiesRouter.GET("/:liabilityName", liabilitiesController.GetLiability)
		liabilitiesRouter.PUT("/:liabilityName", liabilitiesController.UpdateLiability)
		liabilitiesRouter.DELETE("/:liabilityName", liabilitiesController.DeleteLiability)
	}

	liabilityValueRouter := liabilitiesRouter.Group("/value")
	{
		liabilityValueRouter.GET("", liabilitiesController.ListAllLiabilityValues)
		liabilityValueRouter.GET("/:liability", liabilitiesController.GetLiabilityValues)
		liabilityValueRouter.POST("/:liability", liabilitiesController.CreateLiabilityValue)
	}

}

func main() {
	router := gin.Default()
	createAssetControllers(router)
	createLiabilityControllers(router)
	router.Run()
}
