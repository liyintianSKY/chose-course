package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	StudentID   int64  `gorm:"primaryKey;not null" json:"student_id"`
	Password    string `gorm:"size:255;not null" json:"password"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Grade       int64  `gorm:"index;not null" json:"grade"`
	Class       int64  `gorm:"index;" json:"class"`
	PhoneNumber string `gorm:"size:20;not null" json:"phone_number"`
}
