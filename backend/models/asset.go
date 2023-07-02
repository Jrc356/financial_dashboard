package models

import (
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

type AssetType string

const (
	Savings    AssetType = "savings"
	Checking   AssetType = "checking"
	Retirement AssetType = "retirement"
	HSA        AssetType = "hsa"
	RealEstate AssetType = "real-estate"
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

func ValidateAsset(asset Asset) error {
	if asset.Name == "" {
		return fmt.Errorf("no asset name provided")
	}
	switch asset.Type {
	case Savings:
	case Checking:
	case Retirement:
	case RealEstate:
	case HSA:
	default:
		return fmt.Errorf("unknown or invalid asset type: %s", asset.Type)
	}

	if asset.Type == Retirement && asset.TaxBucket == "" {
		return fmt.Errorf("no tax bucket provided for retirement asset: %s", asset.TaxBucket)
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

type AssetValue struct {
	ID        uint
	AssetName string
	Value     float64 `json:"value"`
	CreatedAt time.Time
}

func CreateAssetValue(db *gorm.DB, av AssetValue) error {
	av.Value = math.Round(av.Value*100) / 100
	result := db.Create(&av)
	return result.Error
}

func GetAllAssetValues(db *gorm.DB) ([]AssetValue, error) {
	var assetValues []AssetValue
	result := db.Order("created_at desc").Find(&assetValues)
	return assetValues, result.Error
}

func CalculateTotalAssetValue(db *gorm.DB) (float64, error) {
	assets, err := GetAllAssets(db)
	if err != nil {
		return 0, err
	}

	var totalAssets float64
	for _, asset := range assets {
		assetValue, err := GetLastAssetValue(db, asset.Name)
		if err != nil {
			return 0, err
		}
		totalAssets += assetValue
	}
	return totalAssets, nil
}

func GetAssetValues(db *gorm.DB, assetName string) ([]AssetValue, error) {
	var assetValues []AssetValue
	result := db.Order("created_at desc").Where("asset_name = ?", assetName).Find(&assetValues)
	return assetValues, result.Error
}

func GetLastAssetValue(db *gorm.DB, assetName string) (float64, error) {
	var assetValue AssetValue
	result := db.Order("created_at desc").Where("asset_name = ?", assetName).Find(&assetValue).Limit(1)
	return assetValue.Value, result.Error
}
