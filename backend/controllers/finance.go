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

func (fc *FinanceController) CalculateNetWorth(context *gin.Context) {
	totalAssets, err := models.CalculateTotalAssetValue(fc.DB)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	totalLiabilities, err := models.CalculateTotalLiabilityValue(fc.DB)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	diff := totalAssets - totalLiabilities
	rounded2d := math.Round(diff*float64(100)) / float64(100)
	context.JSON(http.StatusOK, rounded2d) // fmt.Sprintf("%.2f", diff))
}
