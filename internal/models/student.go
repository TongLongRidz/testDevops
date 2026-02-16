package models

type Student struct {
	StudentID     uint       `gorm:"primaryKey;column:student_id" json:"student_id"`
	UserID        uint       `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 Student
	User          User       `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
	StudentNumber string     `gorm:"type:varchar(10);column:student_number" json:"student_number"`
	FacultyID     uint       `gorm:"column:faculty_id" json:"faculty_id"`
	Faculty       Faculty    `gorm:"foreignKey:FacultyID"`
	DepartmentID  uint       `gorm:"column:department_id" json:"department_id"`
	Department    Department `gorm:"foreignKey:DepartmentID"`
}

func (Student) TableName() string {
	return "Student"
}
