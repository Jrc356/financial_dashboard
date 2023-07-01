package controllers

import (
	"log"
	"net/http"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	if err := models.ValidateAsset(asset); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := controller.DB.Create(&asset)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, asset)
}

func (controller *AssetController) ListAssets(context *gin.Context) {
	var assets []models.Asset
	result := controller.DB.Find(&assets)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, assets)
}

func (controller *AssetController) GetAsset(context *gin.Context) {
	var asset models.Asset
	result := controller.DB.First(&asset, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, asset)
}

func (controller *AssetController) UpdateAsset(context *gin.Context) {
	var updatedAsset models.Asset
	if err := context.BindJSON(&updatedAsset); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.ValidateAsset(updatedAsset); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var asset models.Asset
	result := controller.DB.First(&asset, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	result = controller.DB.Model(&asset).Updates(&updatedAsset)
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, asset)
}

func (controller *AssetController) DeleteAsset(context *gin.Context) {
	var asset models.Asset
	result := controller.DB.First(&asset, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	result = controller.DB.Delete(&asset)
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, asset)
}

func (controller *AssetController) CreateAssetValue(context *gin.Context) {
	var assetValue models.AssetValue
	if err := context.BindJSON(&assetValue); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetValue.AssetName = context.Param("asset")
	result := controller.DB.Create(&assetValue)
	if result.Error != nil {
		// TODO: better handling
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, assetValue)
}

func (controller *AssetController) ListAllAssetValues(context *gin.Context) {
	var assetValues []models.AssetValue
	result := controller.DB.Find(&assetValues)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, assetValues)
}

func (controller *AssetController) GetAssetValues(context *gin.Context) {
	var assetValues []models.AssetValue
	result := controller.DB.Where("asset_name = ?", context.Param("asset")).Find(&assetValues)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, assetValues)
}
