package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Enrollment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	StudentID  uuid.UUID `gorm:"type:uuid"`
	SubjectID  uint      `json:"subject_id"`
	SemesterID uint      `json:"semester_id"`
	Grade      float64   `json:"grade"`

	Student  Student  `gorm:"foreignKey:StudentID"`  // Belongs To relationship with Student model
	Subject  Subject  `gorm:"foreignKey:SubjectID"`  // Belongs To relationship with Subject model
	Semester Semester `gorm:"foreignKey:SemesterID"` // Belongs To relationship with Semester model
}

// Function to get enrollments by student ID
func GetEnrollmentsByStudentID(db *gorm.DB, studentID uuid.UUID) ([]Enrollment, error) {
	var enrollments []Enrollment
	if err := db.Where("student_id = ?", studentID).Preload("Subject").Preload("Semester").Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}
