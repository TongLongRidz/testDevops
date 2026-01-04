package academicyear

import "time"

type CreateAcademicYearDTO struct {
    Year      int       `json:"year"`
    Semester  int       `json:"semester"`
    StartDate time.Time `json:"start_date"`
    EndDate   time.Time `json:"end_date"`
}

type UpdateAcademicYearDTO struct {
    Year      int       `json:"year"`
    Semester  int       `json:"semester"`
    StartDate time.Time `json:"start_date"`
    EndDate   time.Time `json:"end_date"`
}

type AcademicYearResponseDTO struct {
    AcademicYearID uint      `json:"academic_year_id"`
    Year           int       `json:"year"`
    Semester       int       `json:"semester"`
    StartDate      time.Time `json:"start_date"`
    EndDate        time.Time `json:"end_date"`
}