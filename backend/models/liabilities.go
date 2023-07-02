package models

import (
	"math"
	"time"

	"gorm.io/gorm"
)

type Liability struct {
	Name   string         `json:"name" gorm:"primaryKey" binding:"required"`
	Values LiabilityValue `gorm:"foreignKey:LiabilityName;references:Name"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func CreateLiability(db *gorm.DB, liability Liability) error {
	result := db.Create(&liability)
	return result.Error
}

func GetAllLiabilities(db *gorm.DB) ([]Liability, error) {
	var liabilities []Liability
	result := db.Find(&liabilities)
	return liabilities, result.Error
}

func GetLiability(db *gorm.DB, liabilityName string) (Liability, error) {
	var liability Liability
	result := db.Where("name = ?", liabilityName).First(&liability)
	return liability, result.Error
}

func UpdateLiability(db *gorm.DB, liabilityName string, updates Liability) (Liability, error) {
	liability, err := GetLiability(db, liabilityName)
	if err != nil {
		return liability, err
	}

	result := db.Model(&liability).Updates(&updates)
	return liability, result.Error
}

func DeleteLiability(db *gorm.DB, liabilityName string) (Liability, error) {
	liability, err := GetLiability(db, liabilityName)
	if err != nil {
		return liability, err
	}

	result := db.Delete(&liability)
	return liability, result.Error
}

type LiabilityValue struct {
	ID            uint
	LiabilityName string
	Value         float64 `json:"value"`
	CreatedAt     time.Time
}

func CreateLiabilityValue(db *gorm.DB, lv LiabilityValue) error {
	lv.Value = math.Round(lv.Value*100) / 100
	result := db.Create(&lv)
	return result.Error
}

func GetAllLiabilityValues(db *gorm.DB) ([]LiabilityValue, error) {
	var liabilityValues []LiabilityValue
	result := db.Order("created_at desc").Find(&liabilityValues)
	return liabilityValues, result.Error
}

func CalculateTotalLiabilityValue(db *gorm.DB) (float64, error) {
	liabilities, err := GetAllLiabilities(db)
	if err != nil {
		return 0, err
	}

	var totalLiabilities float64
	for _, liability := range liabilities {
		liabilityValue, err := GetLastLiabilityValue(db, liability.Name)
		if err != nil {
			return 0, err
		}
		totalLiabilities += liabilityValue
	}
	return totalLiabilities, nil
}

func GetLiabilityValues(db *gorm.DB, liabilityName string) ([]LiabilityValue, error) {
	var liabilityValues []LiabilityValue
	result := db.Order("created_at desc").Where("liability_name = ?", liabilityName).Find(&liabilityValues)
	return liabilityValues, result.Error
}

func GetLastLiabilityValue(db *gorm.DB, liabilityName string) (float64, error) {
	var liabilityValues LiabilityValue
	result := db.Order("created_at desc").Where("liability_name = ?", liabilityName).Find(&liabilityValues).Limit(1)
	return liabilityValues.Value, result.Error
}
