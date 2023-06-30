package asset

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetController struct {
	DB *gorm.DB
}

func (c *AssetController) CreateRoutes(rg *gin.RouterGroup) {
	rg.POST("", c.CreateAsset)
	rg.GET("", c.ListAssets)
	rg.GET("/:id", c.GetAsset)
	rg.PUT("/:id", c.UpdateAsset)
	rg.DELETE("/:id", c.DeleteAsset)
}

func (controller *AssetController) CreateAsset(context *gin.Context) {
	var asset Asset
	if err := context.BindJSON(&asset); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateAsset(asset); err != nil {
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
	var assets []Asset
	result := controller.DB.Find(&assets)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, assets)
}

func (controller *AssetController) GetAsset(context *gin.Context) {
	var asset Asset
	result := controller.DB.First(&asset, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, asset)
}

func (controller *AssetController) UpdateAsset(context *gin.Context) {
	var updatedAsset Asset
	if err := context.BindJSON(&updatedAsset); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateAsset(updatedAsset); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var asset Asset
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
	var asset Asset
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

func ValidateAsset(a Asset) error {
	if a.Name == "" {
		return fmt.Errorf("no asset name provided")
	}
	switch a.Type {
	case Savings:
	case Checking:
	case Retirement:
	case HSA:
	default:
		return fmt.Errorf("unknown or invalid asset type: %s", a.Type)
	}
	return nil
}