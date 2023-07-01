package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	LiabilityNameParam = "liabilityName"
)

type LiabilityController struct {
	DB *gorm.DB
}

func NewLiabilityController(db *gorm.DB, router *gin.Engine) {
	liabilitiesController := LiabilityController{DB: db}

	liabilitiesRouter := router.Group("/liability")
	{
		liabilitiesRouter.POST("", liabilitiesController.CreateLiability)
		liabilitiesRouter.GET("", liabilitiesController.ListLiabilities)
		liabilitiesRouter.GET("/values", liabilitiesController.ListAllLiabilityValues)
	}

	liabilityRouter := liabilitiesRouter.Group("/:" + LiabilityNameParam)
	{
		liabilityRouter.GET("", liabilitiesController.GetLiability)
		liabilityRouter.PUT("", liabilitiesController.UpdateLiability)
		liabilityRouter.DELETE("", liabilitiesController.DeleteLiability)

		liabilityRouter.GET("/value", liabilitiesController.GetLiabilityValues)
		liabilityRouter.POST("/value", liabilitiesController.CreateLiabilityValue)
	}
}

type LiabilityResponse struct {
	Name string
}

func LiabilityToLiabilityResponse(liability models.Liability) LiabilityResponse {
	return LiabilityResponse{
		Name: liability.Name,
	}
}

type LiabilityValueResponse struct {
	LiabilityName string
	Value         float64
	Date          string
}

func LiabilityValueToLiabilityValueResponse(av models.LiabilityValue) LiabilityValueResponse {
	return LiabilityValueResponse{
		LiabilityName: av.LiabilityName,
		Value:         av.Value,
		Date:          av.CreatedAt.Format("01-02-2006"),
	}
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
