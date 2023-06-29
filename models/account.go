package models

import (
	"fmt"

	"gorm.io/gorm"
)

type AccountType string

const (
	Savings    AccountType = "savings"
	Checking   AccountType = "checking"
	Retirement AccountType = "retirement"
	HSA        AccountType = "hsa"
)

func isValidAccount(at AccountType) bool {
	switch at {
	case Savings:
		return true
	case Checking:
		return true
	case Retirement:
		return true
	case HSA:
		return true
	default:
		return false
	}
}

type Account struct {
	gorm.Model
	Name string      `json:"name" gorm:"uniqueIndex"`
	Type AccountType `json:"type"`
	// TODO
}

func (a *Account) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("no account name provided")
	}
	if !isValidAccount(a.Type) {
		return fmt.Errorf("unknown or invalid account type: %s", a.Type)
	}
	return nil
}
