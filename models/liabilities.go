package models

import (
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

type LiabilityValue struct {
	gorm.Model
	LiabilityName string  `json:"liabilityName"`
	Value         float32 `json:"value"`
}
