package models

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountValue struct {
	ID          uint
	AccountName string          `json:"account_name" binding:"required"`
	Value       decimal.Decimal `json:"value" gorm:"type:decimal(19,2)"`
	CreatedAt   time.Time
}

func CreateAccountValue(db *gorm.DB, av AccountValue) error {
	if exists, err := AccountExists(db, av.AccountName); !exists {
		return fmt.Errorf(`account %s does not exist`, av.AccountName)
	} else if err != nil {
		return err
	}

	av.Value = av.Value.Round(2)
	result := db.Create(&av)
	return result.Error
}
