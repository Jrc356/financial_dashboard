package controllers

import (
	"net/http"
	"time"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FinanceController struct {
	DB *gorm.DB
}

func NewFinanceController(db *gorm.DB, router *gin.RouterGroup) {
	financeController := FinanceController{DB: db}
	router.GET("/networth", financeController.GetNetWorthOverTime)
}

func (fc *FinanceController) GetNetWorthOverTime(context *gin.Context) {
	assets, err := models.GetAllAccountsByClassWithValues(fc.DB, models.Asset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	liabilities, err := models.GetAllAccountsByClassWithValues(fc.DB, models.Liability)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	values := make(map[time.Time]float64)

	round := 5 * time.Second
	for _, asset := range assets {
		for _, v := range asset.Values {
			entry := values[v.CreatedAt.Round(round)]
			entry += v.Value
			values[v.CreatedAt.Round(round)] = entry
		}
	}

	for _, liability := range liabilities {
		for _, v := range liability.Values {
			entry := values[v.CreatedAt.Round(round)]
			entry -= v.Value
			values[v.CreatedAt.Round(round)] = entry
		}
	}

	context.JSON(http.StatusOK, values)
}
