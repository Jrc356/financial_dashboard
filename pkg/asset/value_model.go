package asset

import "gorm.io/gorm"

type AssetValue struct {
	gorm.Model
	AssetName string  `json:"assetName"`
	Value     float32 `json:"value"`
}
