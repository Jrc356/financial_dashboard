package models

import (
	"fmt"
	"math"
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
	Values    AccountValue    `gorm:"foreignKey:AccountName;references:Name"`

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
	if err := db.Where("name = ?", name).Count(&count).Error; err != nil {
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
	result := db.Find(&accounts)
	return accounts, result.Error
}

func GetAccountByName(db *gorm.DB, accountName string) (Account, error) {
	var account Account
	result := db.Where("name = ?", accountName).First(&account)
	return account, result.Error
}

func GetAllAccountsByClass(db *gorm.DB, class AccountClass) ([]Account, error) {
	var accounts []Account
	result := db.Where("class = ?", class).Find(&accounts)
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

	result := db.Model(&account).Updates(&updates)
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
