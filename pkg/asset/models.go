package asset

import (
	"time"

	"github.com/Jrc356/financial_dashboard/pkg/tax"
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
	TaxBucket tax.Bucket `json:"taxBucket"`
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
