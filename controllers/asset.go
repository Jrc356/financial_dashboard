package controllers

import (
	"net/http"
	"time"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	paramAssetName   = "assetName"
	routerGroupAsset = "/asset"
)

type AssetController struct {
	DB *gorm.DB
}

func NewAssetController(db *gorm.DB, router *gin.RouterGroup) AssetController {
	assetController := AssetController{DB: db}

	assetsRouter := router.Group(routerGroupAsset)
	{
		assetsRouter.POST("", assetController.CreateAsset)
		assetsRouter.GET("", assetController.ListAssets)
		assetsRouter.GET("/values", assetController.ListAssetValues)
	}

	assetRouter := assetsRouter.Group("/:" + paramAssetName)
	{
		assetRouter.GET("", assetController.GetAsset)
		assetRouter.PUT("", assetController.UpdateAsset)
		assetRouter.DELETE("", assetController.DeleteAsset)

		assetRouter.GET("/value", assetController.GetAssetValues)
		assetRouter.POST("/value", assetController.CreateAssetValue)
	}
	return assetController
}

type assetResponse struct {
	Name      string
	Type      models.AssetType
	TaxBucket models.TaxBucket
}

func assetToAssetResponse(asset models.Asset) assetResponse {
	return assetResponse{
		Name:      asset.Name,
		Type:      asset.Type,
		TaxBucket: asset.TaxBucket,
	}
}

type assetValueResponse struct {
	Value float64
	Date  time.Time
}

func assetValueToAssetValueResponse(av models.AssetValue) assetValueResponse {
	return assetValueResponse{
		Value: av.Value,
		Date:  av.CreatedAt,
	}
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
	context.JSON(http.StatusOK, assetToAssetResponse(asset))
}

func (controller *AssetController) ListAssets(context *gin.Context) {
	assets, err := models.GetAllAssets(controller.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := []assetResponse{}
	for _, asset := range assets {
		responses = append(responses, assetToAssetResponse(asset))
	}
	context.JSON(http.StatusOK, responses)
}

func (controller *AssetController) GetAsset(context *gin.Context) {
	asset, err := models.GetAsset(controller.DB, context.Param(paramAssetName))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, assetToAssetResponse(asset))
}

func (controller *AssetController) UpdateAsset(context *gin.Context) {
	var updatedAsset models.Asset
	if err := context.BindJSON(&updatedAsset); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset, err := models.UpdateAsset(controller.DB, context.Param(paramAssetName), updatedAsset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, assetToAssetResponse(asset))
}

func (controller *AssetController) DeleteAsset(context *gin.Context) {
	asset, err := models.DeleteAsset(controller.DB, context.Param(paramAssetName))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, assetToAssetResponse(asset))
}

func (controller *AssetController) CreateAssetValue(context *gin.Context) {
	var assetValue models.AssetValue
	if err := context.BindJSON(&assetValue); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assetValue.AssetName = context.Param("assetName")
	err := models.CreateAssetValue(controller.DB, assetValue)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, assetValueToAssetValueResponse(assetValue))
}

func (controller *AssetController) ListAssetValues(context *gin.Context) {
	assetValues, err := models.GetAllAssetValues(controller.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := []assetValueResponse{}
	for _, assetValue := range assetValues {
		response = append(response, assetValueToAssetValueResponse(assetValue))
	}
	context.JSON(http.StatusOK, response)
}

func (controller *AssetController) GetAssetValues(context *gin.Context) {
	assetValues, err := models.GetAssetValues(controller.DB, context.Param(paramAssetName))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := []assetValueResponse{}
	for _, assetValue := range assetValues {
		response = append(response, assetValueToAssetValueResponse(assetValue))
	}
	context.JSON(http.StatusOK, response)
}

func (controller *AssetController) GetCurrentAssetValue(context *gin.Context) {
	assetValue, err := models.GetLastAssetValue(controller.DB, context.Param(paramAssetName))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, assetValueToAssetValueResponse(assetValue))
}
