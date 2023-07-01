package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LiabilityController struct {
	DB *gorm.DB
}

func (c *LiabilityController) CreateRoutes(rg *gin.RouterGroup) {
	rg.POST("", c.CreateLiability)
	rg.GET("", c.ListLiabilities)
	rg.GET("/:id", c.GetLiability)
	rg.PUT("/:id", c.UpdateLiability)
	rg.DELETE("/:id", c.DeleteLiability)
}

func (controller *LiabilityController) CreateLiability(context *gin.Context) {
	var liability models.Liability
	if err := context.BindJSON(&liability); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateLiability(liability); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := controller.DB.Create(&liability)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, liability)
}

func (controller *LiabilityController) ListLiabilities(context *gin.Context) {
	var liabilities []models.Liability
	result := controller.DB.Find(&liabilities)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, liabilities)
}

func (controller *LiabilityController) GetLiability(context *gin.Context) {
	var liability models.Liability
	result := controller.DB.First(&liability, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, liability)
}

func (controller *LiabilityController) UpdateLiability(context *gin.Context) {
	var updatedLiability models.Liability
	if err := context.BindJSON(&updatedLiability); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateLiability(updatedLiability); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var liability models.Liability
	result := controller.DB.First(&liability, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	result = controller.DB.Model(&liability).Updates(&updatedLiability)
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, liability)
}

func (controller *LiabilityController) DeleteLiability(context *gin.Context) {
	var liability models.Liability
	result := controller.DB.First(&liability, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	result = controller.DB.Delete(&liability)
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, liability)
}

func (controller *LiabilityController) CreateLiabilityValue(context *gin.Context) {
	var liabilityValue models.LiabilityValue
	if err := context.BindJSON(&liabilityValue); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	liabilityValue.LiabilityName = context.Param("liability")
	result := controller.DB.Create(&liabilityValue)
	if result.Error != nil {
		// TODO: better handling
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, liabilityValue)
}

func (controller *LiabilityController) ListAllLiabilityValues(context *gin.Context) {
	var liabilityValues []models.LiabilityValue
	result := controller.DB.Find(&liabilityValues)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, liabilityValues)
}

func (controller *LiabilityController) GetLiabilityValues(context *gin.Context) {
	var liabilityValues []models.LiabilityValue
	result := controller.DB.Where("liability_name = ?", context.Param("liability")).Find(&liabilityValues)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, liabilityValues)
}

func ValidateLiability(l models.Liability) error {
	if l.Name == "" {
		return fmt.Errorf("no liability name provided")
	}

	return nil
}
