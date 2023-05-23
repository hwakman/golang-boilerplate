package models

import (
	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	Name string `gorm:"not null"`
	File string
}
