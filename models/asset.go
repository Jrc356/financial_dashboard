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

type AssetResponse struct {
	Name      string
	Type      AssetType
	TaxBucket TaxBucket
}

func AssetToAssetResponse(asset Asset) AssetResponse {
	return AssetResponse{
		Name:      asset.Name,
		Type:      asset.Type,
		TaxBucket: asset.TaxBucket,
	}
}

func ValidateAsset(asset Asset) error {
	if asset.Name == "" {
		return fmt.Errorf("no asset name provided")
	}
	switch asset.Type {
	case Savings:
	case Checking:
	case Retirement:
	case HSA:
	default:
		return fmt.Errorf("unknown or invalid asset type: %s", asset.Type)
	}

	if asset.TaxBucket != "" {
		switch asset.TaxBucket {
		case TaxDeferred:
		case Taxable:
		case Roth:
		default:
			return fmt.Errorf("unknown or invalid asset taxBucket: %s", asset.TaxBucket)
		}
	}
	return nil
}

func CreateAsset(db *gorm.DB, asset Asset) error {
	err := ValidateAsset(asset)
	if err != nil {
		return err
	}
	result := db.Create(&asset)
	return result.Error
}

func GetAllAssets(db *gorm.DB) ([]Asset, error) {
	var assets []Asset
	result := db.Find(&assets)
	return assets, result.Error
}

func GetAsset(db *gorm.DB, assetName string) (Asset, error) {
	var asset Asset
	result := db.Where("name = ?", assetName).First(&asset)
	return asset, result.Error
}

func UpdateAsset(db *gorm.DB, assetName string, updates Asset) (Asset, error) {
	if err := ValidateAsset(updates); err != nil {
		return Asset{}, err
	}

	asset, err := GetAsset(db, assetName)
	if err != nil {
		return asset, err
	}

	result := db.Model(&asset).Updates(&updates)
	return asset, result.Error
}

func DeleteAsset(db *gorm.DB, assetName string) (Asset, error) {
	asset, err := GetAsset(db, assetName)
	if err != nil {
		return asset, err
	}

	result := db.Delete(&asset)
	return asset, result.Error
}
