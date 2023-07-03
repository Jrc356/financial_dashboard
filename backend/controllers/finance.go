package controllers

import (
	"math"
	"net/http"

	"github.com/Jrc356/financial_dashboard/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FinanceController struct {
	DB *gorm.DB
}

func NewFinanceController(db *gorm.DB, router *gin.Engine) {
	financeController := FinanceController{DB: db}
	router.GET("/networth", financeController.CalculateNetWorth)
}

func CalculateTotalAssetValue(db *gorm.DB) (float64, error) {
	assets, err := models.GetAllAssets(db)
	if err != nil {
		return 0, err
	}

	var totalAssets float64
	for _, asset := range assets {
		assetValue, err := models.GetLastAssetValue(db, asset.Name)
		if err != nil {
			return 0, err
		}
		totalAssets += assetValue.Value
	}
	return totalAssets, nil
}

func CalculateTotalLiabilityValue(db *gorm.DB) (float64, error) {
	liabilities, err := models.GetAllLiabilities(db)
	if err != nil {
		return 0, err
	}

	var totalLiabilities float64
	for _, liability := range liabilities {
		liabilityValue, err := models.GetLastLiabilityValue(db, liability.Name)
		if err != nil {
			return 0, err
		}
		totalLiabilities += liabilityValue
	}
	return totalLiabilities, nil
}

func (fc *FinanceController) CalculateNetWorth(context *gin.Context) {
	totalAssets, err := CalculateTotalAssetValue(fc.DB)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	totalLiabilities, err := CalculateTotalLiabilityValue(fc.DB)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	diff := totalAssets - totalLiabilities
	rounded2d := math.Round(diff*float64(100)) / float64(100)
	context.JSON(http.StatusOK, rounded2d)
}
