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

func CalculateCurrentTotalValue(db *gorm.DB, class models.AccountClass) (float64, error) {
	accounts, err := models.GetAllAccountsByClass(db, class)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, asset := range accounts {
		total += asset.Values[0].Value
	}
	return total, nil
}

func (fc *FinanceController) CalculateNetWorth(context *gin.Context) {
	totalAssets, err := CalculateCurrentTotalValue(fc.DB, models.Asset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	totalLiabilities, err := CalculateCurrentTotalValue(fc.DB, models.Liability)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	diff := totalAssets - totalLiabilities
	rounded2d := math.Round(diff*float64(100)) / float64(100)
	context.JSON(http.StatusOK, rounded2d)
}
