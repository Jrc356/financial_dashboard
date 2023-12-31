package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AccountClass string

const (
	Asset     AccountClass = "asset"
	Liability AccountClass = "liability"
)

type AccountCategory string

const (
	Cash       AccountCategory = "cash"
	Retirement AccountCategory = "retirement"
	HSA        AccountCategory = "hsa"
	RealEstate AccountCategory = "real-estate"
	Loan       AccountCategory = "loan"
	CreditCard AccountCategory = "credit-card"
)

type TaxBucket string

const (
	TaxDeferred TaxBucket = "tax-deferred"
	Roth        TaxBucket = "roth"
	Taxable     TaxBucket = "taxable"
)

type Account struct {
	Name      string          `json:"name" gorm:"primaryKey" binding:"required"`
	Class     AccountClass    `json:"class"`
	Category  AccountCategory `json:"category"`
	TaxBucket TaxBucket       `json:"taxBucket"`
	Values    []AccountValue  `json:"values" gorm:"foreignKey:AccountName;references:Name"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func ValidateAccount(account Account) error {
	if account.Name == "" {
		return fmt.Errorf("no account name provided")
	}
	switch account.Category {
	case Cash:
	case Retirement:
	case RealEstate:
	case HSA:
	case Loan:
	case CreditCard:
	default:
		return fmt.Errorf("unknown or invalid account category: %s", account.Category)
	}

	if account.Category == Retirement && account.TaxBucket == "" {
		return fmt.Errorf("no tax bucket provided for retirement account: %s", account.TaxBucket)
	}
	return nil
}

func AccountExists(db *gorm.DB, name string) (bool, error) {
	count := int64(0)
	if err := db.Model(&Account{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}

func CreateAccount(db *gorm.DB, account Account) error {
	if err := ValidateAccount(account); err != nil {
		return err
	}
	result := db.Create(&account)
	return result.Error
}

func GetAllAccounts(db *gorm.DB) ([]Account, error) {
	var accounts []Account
	result := db.Preload("Values", func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") }).Find(&accounts)
	return accounts, result.Error
}

func GetAccountByName(db *gorm.DB, accountName string) (Account, error) {
	var account Account
	result := db.Preload("Values", func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") }).Where("name = ?", accountName).First(&account)
	return account, result.Error
}

func GetAllAccountsByClass(db *gorm.DB, class AccountClass) ([]Account, error) {
	var accounts []Account
	result := db.Preload("Values", func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") }).Where("class = ?", class).Find(&accounts)
	return accounts, result.Error
}

func UpdateAccount(db *gorm.DB, accountName string, updates Account) (Account, error) {
	if err := ValidateAccount(updates); err != nil {
		return Account{}, err
	}

	account, err := GetAccountByName(db, accountName)
	if err != nil {
		return account, err
	}

	result := db.Preload("Values", func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") }).Model(&account).Updates(&updates)
	return account, result.Error
}

func DeleteAccount(db *gorm.DB, accountName string) (Account, error) {
	account, err := GetAccountByName(db, accountName)
	if err != nil {
		return account, err
	}

	result := db.Delete(&account)
	return account, result.Error
}