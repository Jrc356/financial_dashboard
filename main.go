package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/Jrc356/financial_dashboard/controllers"
	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func randomDollarAmount() float64 {
	val := rand.Float64() * float64(rand.Int63n(10000))
	ratio := math.Pow(10, 2)
	return math.Round(val*ratio) / ratio
}

func init() {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		"postgres",
		"postgres",
		"10.0.0.202",
		"5432",
		"postgres",
	)

	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Panicln(err)
	}

	db.Migrator().DropTable(&models.Account{})
	db.Migrator().DropTable(&models.AccountValue{})

	db.AutoMigrate(
		&models.Account{},
		&models.AccountValue{},
	)

	// TODO: Remove this test data
	accounts := []models.Account{
		{
			Name:     "Our Savings Account",
			Class:    models.Asset,
			Category: models.Cash,
		},
		{
			Name:     "Our Checking Account",
			Class:    models.Asset,
			Category: models.Cash,
		},
		{
			Name:      "My 401k",
			Class:     models.Asset,
			Category:  models.Retirement,
			TaxBucket: models.TaxDeferred,
		},
		{
			Name:      "SO 401k",
			Class:     models.Asset,
			Category:  models.Retirement,
			TaxBucket: models.TaxDeferred,
		},
		{
			Name:      "My IRA",
			Class:     models.Asset,
			Category:  models.Retirement,
			TaxBucket: models.Roth,
		},
		{
			Name:      "SO IRA",
			Class:     models.Asset,
			Category:  models.Retirement,
			TaxBucket: models.Roth,
		},
		{
			Name:      "House",
			Class:     models.Asset,
			Category:  models.Retirement,
			TaxBucket: models.Roth,
		},
		{
			Name:     "Student Loan",
			Class:    models.Liability,
			Category: models.Loan,
		},
		{
			Name:     "Mortgage",
			Class:    models.Liability,
			Category: models.Loan,
		},
		{
			Name:     "Auto Loan",
			Class:    models.Liability,
			Category: models.Loan,
		},
		{
			Name:     "Credit Card",
			Class:    models.Liability,
			Category: models.CreditCard,
		},
	}
	for _, account := range accounts {
		if err := models.CreateAccount(db, account); err != nil {
			log.Panic(err)
		}

		for i := 0; i < 10; i++ {
			value := models.AccountValue{
				AccountName: account.Name,
				Value:       randomDollarAmount(),
			}

			if err := models.CreateAccountValue(db, value); err != nil {
				log.Panic(err)
			}
		}
	}
}

func main() {
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./client/build", true)))
	router.Use(cors.Default())
	apiRouter := router.Group("/api")
	controllers.NewAccountController(db, apiRouter)
	controllers.NewFinanceController(db, apiRouter)
	router.Run()
}
