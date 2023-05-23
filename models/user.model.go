package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"unique"`
	Age   int
    ImageURL string
}
