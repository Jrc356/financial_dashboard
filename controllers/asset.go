package controllers

import (
	"net/http"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	AssetNameParam = "assetName"
)

type AssetController struct {
	DB *gorm.DB
}

func (controller *AssetController) CreateAsset(context *gin.Context) {
	var asset models.Asset
	if err := context.BindJSON(&asset); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.CreateAsset(controller.DB, asset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, models.AssetToAssetResponse(asset))
}

func (controller *AssetController) ListAssets(context *gin.Context) {
	assets, err := models.GetAllAssets(controller.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := []models.AssetResponse{}
	for _, asset := range assets {
		responses = append(responses, models.AssetToAssetResponse(asset))
	}
	context.JSON(http.StatusOK, responses)
}

func (controller *AssetController) GetAsset(context *gin.Context) {
	asset, err := models.GetAsset(controller.DB, context.Param(AssetNameParam))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, models.AssetToAssetResponse(asset))
}

func (controller *AssetController) UpdateAsset(context *gin.Context) {
	var updatedAsset models.Asset
	if err := context.BindJSON(&updatedAsset); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset, err := models.UpdateAsset(controller.DB, context.Param(AssetNameParam), updatedAsset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, models.AssetToAssetResponse(asset))
}

func (controller *AssetController) DeleteAsset(context *gin.Context) {
	asset, err := models.DeleteAsset(controller.DB, context.Param(AssetNameParam))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, models.AssetToAssetResponse(asset))
}

func (controller *AssetController) CreateAssetValue(context *gin.Context) {
	var assetValue models.AssetValue
	if err := context.BindJSON(&assetValue); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetValue.AssetName = context.Param("assetName")
	assetValue, err := models.CreateAssetValue(controller.DB, assetValue)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, models.AssetValueToAssetValueResponse(assetValue))
}

func (controller *AssetController) ListAllAssetValues(context *gin.Context) {
	assetValues, err := models.GetAllAssetValues(controller.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := []models.AssetValueResponse{}
	for _, assetValue := range assetValues {
		responses = append(responses, models.AssetValueToAssetValueResponse(assetValue))
	}
	context.JSON(http.StatusOK, responses)
}

func (controller *AssetController) GetAssetValues(context *gin.Context) {
	assetValues, err := models.GetAssetValues(controller.DB, context.Param(AssetNameParam))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	responses := []models.AssetValueResponse{}
	for _, assetValue := range assetValues {
		responses = append(responses, models.AssetValueToAssetValueResponse(assetValue))
	}
	context.JSON(http.StatusOK, responses)
}
