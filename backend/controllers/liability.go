package controllers

import (
	"fmt"
	"net/http"

	"github.com/Jrc356/financial_dashboard/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	paramLiabilityName   = "liabilityName"
	routerGroupLiability = "/liability"
)

type LiabilityController struct {
	DB *gorm.DB
}

func NewLiabilityController(db *gorm.DB, router *gin.Engine) {
	liabilitiesController := LiabilityController{DB: db}

	liabilitiesRouter := router.Group(routerGroupLiability)
	{
		liabilitiesRouter.POST("", liabilitiesController.CreateLiability)
		liabilitiesRouter.GET("", liabilitiesController.ListLiabilities)
		liabilitiesRouter.GET("/values", liabilitiesController.ListAllLiabilityValues)
	}

	liabilityRouter := liabilitiesRouter.Group("/:" + paramLiabilityName)
	{
		liabilityRouter.GET("", liabilitiesController.GetLiability)
		liabilityRouter.PUT("", liabilitiesController.UpdateLiability)
		liabilityRouter.DELETE("", liabilitiesController.DeleteLiability)

		liabilityRouter.GET("/value", liabilitiesController.GetLiabilityValues)
		liabilityRouter.POST("/value", liabilitiesController.CreateLiabilityValue)
	}
}

type liabilityResponse struct {
	Name string
}

func liabilityToLiabilityResponse(liability models.Liability) liabilityResponse {
	return liabilityResponse{
		Name: liability.Name,
	}
}

type liabilityValueResponse struct {
	LiabilityName string
	Value         float64
	Date          string
}

func liabilityValueToLiabilityValueResponse(av models.LiabilityValue) liabilityValueResponse {
	return liabilityValueResponse{
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

	err := models.CreateLiability(controller.DB, liability)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, liabilityToLiabilityResponse(liability))
}

func (controller *LiabilityController) ListLiabilities(context *gin.Context) {
	liabilities, err := models.GetAllLiabilities(controller.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := []liabilityResponse{}
	for _, liability := range liabilities {
		response = append(response, liabilityToLiabilityResponse(liability))
	}
	context.JSON(http.StatusOK, response)
}

func (controller *LiabilityController) GetLiability(context *gin.Context) {
	liability, err := models.GetLiability(controller.DB, context.Param(paramLiabilityName))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, liabilityToLiabilityResponse(liability))
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

	liability, err := models.UpdateLiability(controller.DB, context.Param(paramLiabilityName), updatedLiability)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, liabilityToLiabilityResponse(liability))
}

func (controller *LiabilityController) DeleteLiability(context *gin.Context) {
	liability, err := models.DeleteLiability(controller.DB, context.Param(paramLiabilityName))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, liabilityToLiabilityResponse(liability))
}

func (controller *LiabilityController) CreateLiabilityValue(context *gin.Context) {
	var liabilityValue models.LiabilityValue
	if err := context.BindJSON(&liabilityValue); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	liabilityValue.LiabilityName = context.Param("liability")
	err := models.CreateLiabilityValue(controller.DB, liabilityValue)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, liabilityValueToLiabilityValueResponse(liabilityValue))
}

func (controller *LiabilityController) ListAllLiabilityValues(context *gin.Context) {
	liabilityValues, err := models.GetAllLiabilityValues(controller.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := []liabilityValueResponse{}
	for _, liabilityValue := range liabilityValues {
		response = append(response, liabilityValueToLiabilityValueResponse(liabilityValue))
	}
	context.JSON(http.StatusOK, response)
}

func (controller *LiabilityController) GetLiabilityValues(context *gin.Context) {
	liabilityValues, err := models.GetLiabilityValues(controller.DB, context.Param(paramLiabilityName))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := []liabilityValueResponse{}
	for _, liabilityValue := range liabilityValues {
		response = append(response, liabilityValueToLiabilityValueResponse(liabilityValue))
	}
	context.JSON(http.StatusOK, response)
}

func ValidateLiability(l models.Liability) error {
	if l.Name == "" {
		return fmt.Errorf("no liability name provided")
	}

	return nil
}
