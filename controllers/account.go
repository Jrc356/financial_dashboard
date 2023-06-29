package controllers

import (
	"log"
	"net/http"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountController struct {
	DB *gorm.DB
}

func (c *AccountController) CreateRoutes(rg *gin.RouterGroup) {
	rg.POST("", c.CreateAccount)
	rg.GET("", c.ListAccounts)
	rg.GET("/:id", c.GetAccount)
	rg.PUT("/:id", c.UpdateAccount)
	rg.DELETE("/:id", c.DeleteAccount)
}

func (controller *AccountController) CreateAccount(context *gin.Context) {
	var account models.Account
	if err := context.BindJSON(&account); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := account.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := controller.DB.Create(&account)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, account)
}

func (controller *AccountController) ListAccounts(context *gin.Context) {
	var accounts []models.Account
	result := controller.DB.Find(&accounts)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, accounts)
}

func (controller *AccountController) GetAccount(context *gin.Context) {
	var account models.Account
	result := controller.DB.First(&account, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, account)
}

func (controller *AccountController) UpdateAccount(context *gin.Context) {
	var updatedAccount models.Account
	if err := context.BindJSON(&updatedAccount); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := updatedAccount.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account models.Account
	result := controller.DB.First(&account, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	result = controller.DB.Model(&account).Updates(&updatedAccount)
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, account)
}

func (controller *AccountController) DeleteAccount(context *gin.Context) {
	var account models.Account
	result := controller.DB.First(&account, context.Param("id"))
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	result = controller.DB.Delete(&account)
	if result.Error != nil {
		log.Println(result.Error)
		context.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	context.JSON(http.StatusOK, account)
}
