package studentdto

type CreateStudentRequest struct {
	StudentNumber string `json:"student_number" binding:"required,min=1,max=20"`
	FacultyID     uint   `json:"faculty_id" binding:"required"`
	DepartmentID  uint   `json:"department_id" binding:"required"`
}

type StudentResponse struct {
	StudentID     uint   `json:"student_id"`
	UserID        uint   `json:"user_id"`
	StudentNumber string `json:"student_number"`
	FacultyID     uint   `json:"faculty_id"`
	DepartmentID  uint   `json:"department_id"`
}

type StudentWithUserResponse struct {
	StudentID     uint   `json:"student_id"`
	UserID        uint   `json:"user_id"`
	StudentNumber string `json:"student_number"`
	FacultyID     uint   `json:"faculty_id"`
	DepartmentID  uint   `json:"department_id"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Email         string `json:"email"`
	ImagePath     string `json:"image_path"`
}
