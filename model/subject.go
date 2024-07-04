package model

type Subject struct {
	ID          uint         `gorm:"primaryKey"`
	SubjectName string       `json:"subject_name"`
	Credits     int          `json:"credits"`
	Enrollments []Enrollment `gorm:"foreignKey:SubjectID"` // One-to-Many relationship with Enrollment
}
