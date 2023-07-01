package models

import (
	"time"

	"gorm.io/gorm"
)

type AssetValue struct {
	ID        uint
	AssetName string  `json:"assetName"`
	Value     float32 `json:"value"`
	CreatedAt time.Time
}

type AssetValueResponse struct {
	AssetName string
	Value     float32
	Date      time.Time //string
}

func AssetValueToAssetValueResponse(av AssetValue) AssetValueResponse {
	return AssetValueResponse{
		AssetName: av.AssetName,
		Value:     av.Value,
		Date:      av.CreatedAt, // av.CreatedAt.Format("01-02-2006"),
	}
}

func CreateAssetValue(db *gorm.DB, av AssetValue) (AssetValue, error) {
	result := db.Create(&av)
	return av, result.Error
}

func GetAllAssetValues(db *gorm.DB) ([]AssetValue, error) {
	var assetValues []AssetValue
	result := db.Order("created_at desc").Find(&assetValues)
	return assetValues, result.Error
}

func GetAssetValues(db *gorm.DB, assetName string) ([]AssetValue, error) {
	var assetValues []AssetValue
	result := db.Order("created_at desc").Where("asset_name = ?", assetName).Find(&assetValues)
	return assetValues, result.Error
}
