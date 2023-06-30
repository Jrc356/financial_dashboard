package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetValueController struct {
	DB              *gorm.DB
	AssetController AssetController
}

func (c *AssetValueController) CreateRoutes(rg *gin.RouterGroup) {
	rg.GET("", c.ListAllAssetValues)
	rg.GET("/:asset", c.GetAssetValues)
	rg.POST("/:asset", c.CreateAssetValue)
}

func (controller *AssetValueController) CreateAssetValue(context *gin.Context) {
	var assetValue AssetValue
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

func (controller *AssetValueController) ListAllAssetValues(context *gin.Context) {
	var assetValues []AssetValue
	result := controller.DB.Find(&assetValues)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, assetValues)
}

func (controller *AssetValueController) GetAssetValues(context *gin.Context) {
	var assetValues []AssetValue
	result := controller.DB.Where("asset_name = ?", context.Param("asset")).Find(&assetValues)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, assetValues)
}
