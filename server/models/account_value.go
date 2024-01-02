package models

import (
	"fmt"
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
	if exists := AccountExists(db, av.AccountName); !exists {
		return fmt.Errorf(`account %s does not exist`, av.AccountName)
	}
	av.Value = math.Round(av.Value*100) / 100
	result := db.Create(&av)
	return result.Error
}
