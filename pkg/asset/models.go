package asset

import (
	"github.com/Jrc356/financial_dashboard/pkg/tax"
	"gorm.io/gorm"
)

// TODO: think this would be best converted into assets and add a tax bucket column
// see the assets table I have in sheets

type AssetType string

const (
	Savings    AssetType = "savings"
	Checking   AssetType = "checking"
	Retirement AssetType = "retirement"
	HSA        AssetType = "hsa"
)

type Asset struct {
	gorm.Model
	Name         string     `json:"name" gorm:"uniqueIndex" binding:"required"`
	Type         AssetType  `json:"type"`
	TaxBucket    tax.Bucket `json:"taxBucket"`
	CurrentValue float32    `json:"currentValue"`
}

type AssetValue struct {
	gorm.Model
	Asset Asset   `json:"asset" binding:"required"`
	Value float32 `json:"value" binding:"required"`
}
