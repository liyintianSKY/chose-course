package models

import "time"

type Course struct {
	CourseID    int64  `gorm:"primaryKey;autoIncrement;" json:"course_id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	Teacher     string `gorm:"size:100;not null" json:"teacher"`
	Credits     int64  `gorm:"not null" json:"credits"`     // 课程学分
	MaxCapacity int64  `gorm:"not null" json:"maxCapacity"` // 课程的最大容量
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
