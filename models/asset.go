package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AssetType string

const (
	Savings    AssetType = "savings"
	Checking   AssetType = "checking"
	Retirement AssetType = "retirement"
	HSA        AssetType = "hsa"
)

type Asset struct {
	Name      string     `json:"name" gorm:"primaryKey" binding:"required"`
	Type      AssetType  `json:"type"`
	TaxBucket TaxBucket  `json:"taxBucket"`
	Values    AssetValue `gorm:"foreignKey:AssetName;references:Name"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type AssetValue struct {
	gorm.Model
	AssetName string  `json:"assetName"`
	Value     float32 `json:"value"`
}

func ValidateAsset(a Asset) error {
	if a.Name == "" {
		return fmt.Errorf("no asset name provided")
	}
	switch a.Type {
	case Savings:
	case Checking:
	case Retirement:
	case HSA:
	default:
		return fmt.Errorf("unknown or invalid asset type: %s", a.Type)
	}

	if a.TaxBucket != "" {
		switch a.TaxBucket {
		case TaxDeferred:
		case Taxable:
		case Roth:
		default:
			return fmt.Errorf("unknown or invalid asset taxBucket: %s", a.TaxBucket)
		}
	}

	return nil
}
