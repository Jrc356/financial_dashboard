package controllers

import (
	"net/http"

	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountController struct {
	DB *gorm.DB
}

func NewAccountController(db *gorm.DB, router *gin.RouterGroup) AccountController {
	accountController := AccountController{DB: db}

	accountRouter := router.Group("/accounts")
	{
		accountRouter.GET("", accountController.GetAccounts)
		accountRouter.POST("", accountController.CreateOrUpdateAccount)
		accountRouter.DELETE("", accountController.DeleteAccount)

		accountRouter.GET("/value", accountController.GetAccountValues)
		accountRouter.POST("/value", accountController.CreateAccountValue)
	}

	return accountController
}

func (controller *AccountController) CreateOrUpdateAccount(context *gin.Context) {
	var account models.Account

	name := context.Query("name")
	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unset parameter 'name' required."})
		return
	}

	if err := context.BindJSON(&account); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := models.AccountExists(controller.DB, name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		err := models.CreateAccount(controller.DB, account)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		account, err = models.UpdateAccount(controller.DB, name, account)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	context.JSON(http.StatusOK, account)
}

func (controller *AccountController) GetAccounts(context *gin.Context) {
	name := context.Query("name")
	var class models.AccountClass = models.AccountClass(context.Query("class"))

	if name != "" {
		account, err := models.GetAccountByNameWithValues(controller.DB, name)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, account)
		return
	}

	if class != "" {
		accounts, err := models.GetAllAccountsByClass(controller.DB, class)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, accounts)
		return
	}

	accounts, err := models.GetAllAccounts(controller.DB)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, accounts)
}

func (controller *AccountController) DeleteAccount(context *gin.Context) {
	name := context.Query("name")
	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unset parameter 'name' required."})
		return
	}

	exists, err := models.AccountExists(controller.DB, name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		context.JSON(http.StatusBadRequest, "Account does not exist")
		return
	} else {
		account, err := models.DeleteAccount(controller.DB, name)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, account)
	}
}

func (controller *AccountController) CreateAccountValue(context *gin.Context) {
	name := context.Query("name")
	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unset parameter 'name' required."})
		return
	}

	exists, err := models.AccountExists(controller.DB, name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		context.JSON(http.StatusBadRequest, "Account does not exist")
		return
	} else {
		var accountValue models.AccountValue
		if err := context.BindJSON(&accountValue); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		accountValue.AccountName = name
		err := models.CreateAccountValue(controller.DB, accountValue)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, accountValue)
	}
}

func (controller *AccountController) GetAccountValues(context *gin.Context) {
	name := context.Query("name")
	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unset parameter 'name' required."})
		return
	}

	exists, err := models.AccountExists(controller.DB, name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		context.JSON(http.StatusBadRequest, "Account does not exist")
		return
	} else {
		accountValues, err := models.GetAccountValues(controller.DB, name)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, accountValues)
	}
}

func (controller *AccountController) GetCurrentAccountValue(context *gin.Context) {
	name := context.Query("name")
	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Unset parameter 'name' required."})
		return
	}

	exists, err := models.AccountExists(controller.DB, name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		accountValue, err := models.GetLastAccountValue(controller.DB, name)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, accountValue)
	}
}
