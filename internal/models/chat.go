package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Title    string    `validate:"required,min=1,max=200"`
	Messages []Message `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}
