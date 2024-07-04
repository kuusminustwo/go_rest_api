package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;"`
	Student_ID    string    `json:"student_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	Firstname     string    `json:"firstname"`
	Lastname      string    `json:"lastname"`
	Register      string    `json:"register"`
	Major         string    `json:"major"`
	Enrolldate    string    `json:"enrolldate"`
	Undergraduate string    `json:"undergraduate"`
	GPA           float32   `json:"gpa"`
	TotalCredits  int       `json:"total_credits"`
}

type Students struct {
	Students []Student `json:"students"`
}

func (student *Student) BeforeCreate(tx *gorm.DB) (err error) {
	student.ID = uuid.New()
	return
}
