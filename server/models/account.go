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

func (ac AccountClass) String() string {
	return string(ac)
}

func ParseAccountClass(s string) (ac AccountClass, err error) {
	categories := map[AccountClass]struct{}{
		Asset:     {},
		Liability: {},
	}
	cls := AccountClass(s)
	_, ok := categories[cls]
	if !ok {
		return ac, fmt.Errorf(`unknown or invalid account class: %s`, s)
	}
	return cls, nil
}

type AccountCategory string

const (
	Cash       AccountCategory = "cash"
	Retirement AccountCategory = "retirement"
	HSA        AccountCategory = "hsa"
	RealEstate AccountCategory = "real-estate"
	Loan       AccountCategory = "loan"
	CreditCard AccountCategory = "credit-card"
)

func (ac AccountCategory) String() string {
	return string(ac)
}

func ParseAccountCategory(s string) (ac AccountCategory, err error) {
	categories := map[AccountCategory]struct{}{
		Cash:       {},
		Retirement: {},
		HSA:        {},
		RealEstate: {},
		Loan:       {},
		CreditCard: {},
	}
	cat := AccountCategory(s)
	_, ok := categories[cat]
	if !ok {
		return ac, fmt.Errorf(`unknown or invalid account category: %s`, s)
	}
	return cat, nil
}

type TaxBucket string

const (
	TaxDeferred TaxBucket = "tax-deferred"
	Roth        TaxBucket = "roth"
	Taxable     TaxBucket = "taxable"
)

func (tb TaxBucket) String() string {
	return string(tb)
}

func ParseTaxBucket(s string) (tb TaxBucket, err error) {
	buckets := map[TaxBucket]struct{}{
		TaxDeferred: {},
		Roth:        {},
		Taxable:     {},
	}
	t := TaxBucket(s)
	_, ok := buckets[t]
	if !ok {
		return tb, fmt.Errorf(`unknown or invalid account category: %s`, s)
	}
	return t, nil
}

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
	if account.Class == "" {
		return fmt.Errorf("no account class provided")
	}

	if account.Category == "" {
		return fmt.Errorf("no account class provided")
	}

	_, err := ParseAccountCategory(account.Category.String())
	if err != nil {
		return err
	}

	_, err = ParseAccountClass(account.Class.String())
	if err != nil {
		return err
	}

	if account.Category == Retirement {
		if account.TaxBucket == "" {
			return fmt.Errorf("no tax bucket provided for retirement account: %s", account.TaxBucket)
		}
		_, err := ParseTaxBucket(account.TaxBucket.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func AccountExists(db *gorm.DB, name string) (bool, error) {
	count := int64(0)
	db.Model(&Account{}).Where("name = ?", name).Count(&count)
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

func GetAccountByNameWithValues(db *gorm.DB, accountName string) (Account, error) {
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

	var account Account
	result := db.Where("name = ?", accountName).First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	result = db.Model(&account).Omit("Values").Updates(&updates)
	return account, result.Error
}

func DeleteAccount(db *gorm.DB, accountName string) (Account, error) {
	var account Account
	result := db.Where("name = ?", accountName).First(&account)
	if result.Error != nil {
		return account, result.Error
	}

	result = db.Delete(&account)
	return account, result.Error
}
