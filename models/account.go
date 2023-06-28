package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
	// TODO
}
