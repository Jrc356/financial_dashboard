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
		"localhost",
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

	db.Migrator().DropTable(&models.Asset{})
	db.Migrator().DropTable(&models.AssetValue{})
	db.Migrator().DropTable(&models.Liability{})
	db.Migrator().DropTable(&models.LiabilityValue{})

	db.AutoMigrate(
		&models.Asset{},
		&models.AssetValue{},
		&models.Liability{},
		&models.LiabilityValue{},
	)

	// TODO: Remove this test data
	assets := []models.Asset{
		{
			Name: "Our Savings Account",
			Type: models.Savings,
		},
		{
			Name: "Our Checking Account",
			Type: models.Checking,
		},
		{
			Name:      "My 401k",
			Type:      models.Retirement,
			TaxBucket: models.TaxDeferred,
		},
		{
			Name:      "SO 401k",
			Type:      models.Retirement,
			TaxBucket: models.TaxDeferred,
		},
		{
			Name:      "My IRA",
			Type:      models.Retirement,
			TaxBucket: models.Roth,
		},
		{
			Name:      "SO IRA",
			Type:      models.Retirement,
			TaxBucket: models.Roth,
		},
		{
			Name:      "House",
			Type:      models.Retirement,
			TaxBucket: models.Roth,
		},
	}
	for _, asset := range assets {
		if err := models.CreateAsset(db, asset); err != nil {
			log.Panic(err)
		}

		for i := 0; i < 10; i++ {
			value := models.AssetValue{
				AssetName: asset.Name,
				Value:     randomDollarAmount(),
			}

			if err := models.CreateAssetValue(db, value); err != nil {
				log.Panic(err)
			}
		}
	}

	liabilities := []string{
		"Student Loan",
		"Mortgage",
		"Auto Loan",
		"Credit Card",
	}
	for _, liabilityName := range liabilities {
		liability := models.Liability{
			Name: liabilityName,
		}
		if err := models.CreateLiability(db, liability); err != nil {
			log.Panic(err)
		}

		for i := 0; i < 10; i++ {
			value := models.LiabilityValue{
				LiabilityName: liabilityName,
				Value:         randomDollarAmount(),
			}
			if err := models.CreateLiabilityValue(db, value); err != nil {
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
	controllers.NewAssetController(db, apiRouter)
	controllers.NewLiabilityController(db, apiRouter)
	controllers.NewFinanceController(db, apiRouter)
	router.Run()
}
