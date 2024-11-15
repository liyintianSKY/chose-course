package models

import (
	"gorm.io/gorm"
)

type Enrollment struct {
	gorm.Model
	StudentID int64  `gorm:"index;not null" json:"student_id"` // 外键，关联学生
	CourseID  int64  `gorm:"index;not null" json:"course_id"`  // 外键，关联课程
	Status    string `gorm:"size:20;not null" json:"status"`   // 选课状态（如：已选、待抢、已取消）
}
