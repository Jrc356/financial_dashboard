package main

import (
	"fmt"
	"log"

	"github.com/Jrc356/financial_dashboard/controllers"
	"github.com/Jrc356/financial_dashboard/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

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
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}

	db.AutoMigrate(
		&models.Account{},
	)
}

func main() {
	router := gin.Default()
	accountsController := controllers.AccountController{DB: db}
	accountsRouter := router.Group("/account")
	accountsController.CreateRoutes(accountsRouter)
	router.Run()
}
