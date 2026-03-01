package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChatID uint   `gorm:"not null;index"`
	Text   string `validate:"required,min=1,max=5000"`
}
