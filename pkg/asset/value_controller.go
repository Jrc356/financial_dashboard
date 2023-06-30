package asset

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetValueController struct {
	DB              *gorm.DB
	AssetController AssetController
}

func (c *AssetValueController) CreateRoutes(rg *gin.RouterGroup) {
	rg.POST("/new", c.CreateAssetValue)
	rg.GET("", c.ListAssetValues)
	rg.GET("/:asset", c.GetAssetValue)
}

func (controller *AssetValueController) CreateAssetValue(context *gin.Context) {
	var assetValue AssetValue
	if err := context.BindJSON(&assetValue); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var asset Asset
	result := controller.DB.First(&assetValue.Asset)
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	assetValue.Asset = asset
	result = controller.DB.Create(&assetValue)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, assetValue)
}

func (controller *AssetValueController) ListAssetValues(context *gin.Context) {
	var assetValues []AssetValue
	result := controller.DB.Find(&assetValues)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, assetValues)
}

func (controller *AssetValueController) GetAssetValue(context *gin.Context) {
	assetValue := AssetValue{Asset: Asset{Name: context.Param("asset")}}
	result := controller.DB.Find(&assetValue)
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, assetValue)
}
