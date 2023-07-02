package controllers

import (
	"net/http"

	"github.com/Jrc356/financial_dashboard/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	AssetNameParam = "assetName"
)

type AssetController struct {
	DB *gorm.DB
}

func NewAssetController(db *gorm.DB, router *gin.Engine) {
	assetsController := AssetController{DB: db}

	assetsRouter := router.Group("/asset")
	{
		assetsRouter.POST("", assetsController.CreateAsset)
		assetsRouter.GET("", assetsController.ListAssets)
		assetsRouter.GET("/values", assetsController.ListAllAssetValues)
	}

	assetRouter := assetsRouter.Group("/:" + AssetNameParam)
	{
		assetRouter.GET("", assetsController.GetAsset)
		assetRouter.PUT("", assetsController.UpdateAsset)
		assetRouter.DELETE("", assetsController.DeleteAsset)

		assetRouter.GET("/value", assetsController.GetAssetValues)
		assetRouter.POST("/value", assetsController.CreateAssetValue)

	}
}

type AssetResponse struct {
	Name      string
	Type      models.AssetType
	TaxBucket models.TaxBucket
}

func AssetToAssetResponse(asset models.Asset) AssetResponse {
	return AssetResponse{
		Name:      asset.Name,
		Type:      asset.Type,
		TaxBucket: asset.TaxBucket,
	}
}

type AssetValueResponse struct {
	AssetName string
	Value     float64
	Date      string
}

func AssetValueToAssetValueResponse(av models.AssetValue) AssetValueResponse {
	return AssetValueResponse{
		AssetName: av.AssetName,
		Value:     av.Value,
		Date:      av.CreatedAt.Format("01-02-2006"),
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
	context.JSON(http.StatusOK, AssetToAssetResponse(asset))
}

func (controller *AssetController) ListAssets(context *gin.Context) {
	assets, err := models.GetAllAssets(controller.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := []AssetResponse{}
	for _, asset := range assets {
		responses = append(responses, AssetToAssetResponse(asset))
	}
	context.JSON(http.StatusOK, responses)
}

func (controller *AssetController) GetAsset(context *gin.Context) {
	asset, err := models.GetAsset(controller.DB, context.Param(AssetNameParam))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	asset, err := models.UpdateAsset(controller.DB, context.Param(AssetNameParam), updatedAsset)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, AssetToAssetResponse(asset))
}

func (controller *AssetController) DeleteAsset(context *gin.Context) {
	asset, err := models.DeleteAsset(controller.DB, context.Param(AssetNameParam))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, AssetToAssetResponse(asset))
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
	context.JSON(http.StatusOK, AssetValueToAssetValueResponse(assetValue))
}

func (controller *AssetController) ListAllAssetValues(context *gin.Context) {
	assetValues, err := models.GetAllAssetValues(controller.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := []AssetValueResponse{}
	for _, assetValue := range assetValues {
		responses = append(responses, AssetValueToAssetValueResponse(assetValue))
	}
	context.JSON(http.StatusOK, responses)
}

func (controller *AssetController) GetAssetValues(context *gin.Context) {
	assetValues, err := models.GetAssetValues(controller.DB, context.Param(AssetNameParam))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	responses := []AssetValueResponse{}
	for _, assetValue := range assetValues {
		responses = append(responses, AssetValueToAssetValueResponse(assetValue))
	}
	context.JSON(http.StatusOK, responses)
}
