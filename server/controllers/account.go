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

		accountRouter.POST("/value", accountController.CreateAccountValue)
	}

	return accountController
}

func (controller *AccountController) CreateOrUpdateAccount(context *gin.Context) {
	var account models.Account

	if err := context.BindJSON(&account); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.ValidateAccount(account)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := models.AccountExists(controller.DB, account.Name)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		err := models.CreateAccount(controller.DB, account)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			context.JSON(http.StatusCreated, account)
		}
	} else {
		var err error
		account, err = models.UpdateAccount(controller.DB, account.Name, account)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			context.JSON(http.StatusOK, account)
		}
	}
}

func (controller *AccountController) GetAccounts(context *gin.Context) {
	name := context.Query("name")
	var class models.AccountClass = models.AccountClass(context.Query("class"))

	if name != "" {
		account, err := models.GetAccountByNameWithValues(controller.DB, name)
		if err == gorm.ErrRecordNotFound {
			context.AbortWithStatusJSON(http.StatusNotFound, account)
			return
		}
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, account)
		return
	}

	if class != "" {
		_, err := models.ParseAccountClass(class.String())
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		accounts, err := models.GetAllAccountsByClassWithValues(controller.DB, class)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, accounts)
		return
	}

	accounts, err := models.GetAllAccountsWithValues(controller.DB)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, accounts)
}

func (controller *AccountController) DeleteAccount(context *gin.Context) {
	name := context.Query("name")
	if name == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unset parameter 'name' required."})
		return
	}

	exists, err := models.AccountExists(controller.DB, name)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, "account does not exist")
		return
	}
	if !exists {
		context.AbortWithStatusJSON(http.StatusNotFound, "account does not exist")
		return
	} else {
		account, err := models.DeleteAccount(controller.DB, name)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, account)
	}
}

func (controller *AccountController) CreateAccountValue(context *gin.Context) {
	var accountValue models.AccountValue

	if err := context.BindJSON(&accountValue); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if accountValue.Value.IsZero() {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": `"value" must be > 0`})
		return
	}

	exists, err := models.AccountExists(controller.DB, accountValue.AccountName)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !exists {
		context.AbortWithStatusJSON(http.StatusBadRequest, "Account does not exist")
		return
	} else {
		err := models.CreateAccountValue(controller.DB, accountValue)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, accountValue)
	}
}
