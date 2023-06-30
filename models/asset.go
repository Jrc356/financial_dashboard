package models

import (
	"fmt"

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

func isValidAsset(at AssetType) bool {
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

type Asset struct {
	gorm.Model
	Name string    `json:"name" gorm:"uniqueIndex"`
	Type AssetType `json:"type"`
	// TODO
}

func (a *Asset) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("no asset name provided")
	}
	if !isValidAsset(a.Type) {
		return fmt.Errorf("unknown or invalid asset type: %s", a.Type)
	}
	return nil
}
