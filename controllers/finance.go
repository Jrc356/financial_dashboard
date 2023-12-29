package controllers

import (
	"math"
	"net/http"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FinanceController struct {
	DB *gorm.DB
}

func NewFinanceController(db *gorm.DB, router *gin.RouterGroup) {
	financeController := FinanceController{DB: db}
	router.GET("/networth", financeController.CalculateNetWorth)
}

func CalculateTotalValue(db *gorm.DB, class models.AccountClass) (float64, error) {
	assets, err := models.GetAllAccountsByClass(db, class)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, asset := range assets {
		assetValue, err := models.GetLastAccountValue(db, asset.Name)
		if err != nil {
			return 0, err
		}
		total += assetValue.Value
	}
	return total, nil
}

func (fc *FinanceController) CalculateNetWorth(context *gin.Context) {
	totalAssets, err := CalculateTotalValue(fc.DB, models.Asset)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	totalLiabilities, err := CalculateTotalValue(fc.DB, models.Liability)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	diff := totalAssets - totalLiabilities
	rounded2d := math.Round(diff*float64(100)) / float64(100)
	context.JSON(http.StatusOK, rounded2d)
}
