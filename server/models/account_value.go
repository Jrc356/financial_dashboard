package models

import (
	"math"
	"time"

	"gorm.io/gorm"
)

type AccountValue struct {
	ID          uint
	AccountName string
	Value       float64 `json:"value"`
	CreatedAt   time.Time
}

func CreateAccountValue(db *gorm.DB, av AccountValue) error {
	av.Value = math.Round(av.Value*100) / 100
	result := db.Create(&av)
	return result.Error
}

func GetAllAccountValues(db *gorm.DB) ([]AccountValue, error) {
	var accountValues []AccountValue
	result := db.Order("created_at desc").Find(&accountValues)
	return accountValues, result.Error
}

func GetAccountValues(db *gorm.DB, accountName string) ([]AccountValue, error) {
	var accountValues []AccountValue
	result := db.Order("created_at desc").Where("account_name = ?", accountName).Find(&accountValues)
	return accountValues, result.Error
}

func GetLastAccountValue(db *gorm.DB, accountName string) (AccountValue, error) {
	var accountValue AccountValue
	result := db.Order("created_at desc").Where("account_name = ?", accountName).Find(&accountValue).Limit(1)
	return accountValue, result.Error
}
