package academicyear

import "time"

type CreateAcademicYear struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type UpdateAcademicYear struct {
	Year      int       `json:"year" binding:"required"`
	Semester  int       `json:"semester" binding:"required,min=1,max=2"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type AcademicYearResponse struct {
	AcademicYearID uint      `json:"academic_year_id"`
	Year           int       `json:"year"`
	Semester       int       `json:"semester"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}
