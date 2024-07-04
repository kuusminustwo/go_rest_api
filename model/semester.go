package model

type Semester struct {
	ID           uint         `gorm:"primaryKey"`
	SemesterName string       `json:"semester_name"`
	Enrollments  []Enrollment `gorm:"foreignKey:SemesterID"` // One-to-Many relationship with Enrollment
}
