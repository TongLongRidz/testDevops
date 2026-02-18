package usecase

import (
	awardformdto "backend/internal/dto/award_form_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
)

type AwardUseCase interface {
	// ปรับปรุง: รับ userID เพื่อดึงข้อมูล student และ files
	SubmitAward(ctx context.Context, userID uint, input awardformdto.SubmitAwardRequest, files []models.AwardFileDirectory) error
	GetByKeyword(ctx context.Context, campusID int, keyword string, date string, studentYear int, awardType string, page int, limit int, arrangement string) (*awardformdto.PaginatedAwardResponse, error)
	GetAwardsByUserID(ctx context.Context, userID uint) ([]awardformdto.AwardFormResponse, error)
	GetAwardsByStudentID(ctx context.Context, studentID int) ([]awardformdto.AwardFormResponse, error)
	GetByFormID(ctx context.Context, formID int) (*awardformdto.AwardFormResponse, error)
	IsDuplicate(userID uint, year int, semester int) (bool, error)
	UpdateAwardType(ctx context.Context, formID uint, awardType string, changedBy uint) error
	UpdateFormStatus(ctx context.Context, formID uint, formStatus int, changedBy uint) error
}

type awardUseCase struct {
	repo                *repository.AwardRepository
	logUseCase          AwardFormLogUseCase
	studentService      StudentService
	organizationService OrganizationService
	academicYearService AcademicYearService
}

func NewAwardUseCase(r *repository.AwardRepository, ss StudentService, os OrganizationService, ays AcademicYearService, logUC AwardFormLogUseCase) AwardUseCase {
	return &awardUseCase{
		repo:                r,
		logUseCase:          logUC,
		studentService:      ss,
		organizationService: os,
		academicYearService: ays,
	}
}

func (u *awardUseCase) SubmitAward(ctx context.Context, userID uint, input awardformdto.SubmitAwardRequest, files []models.AwardFileDirectory) error {
	// 1. ดึงข้อมูล Academic Year ที่เปิดรับสมัคร
	academicYear, err := u.academicYearService.GetLatestAbleRegister(ctx)
	if err != nil || academicYear == nil {
		return errors.New("no open registration period found")
	}

	// 2. เตรียม Model ตารางหลัก (Award_Form)
	now := time.Now()
	form := models.AwardForm{
		UserID:             userID,
		AcademicYear:       academicYear.Year,
		Semester:           academicYear.Semester,
		AwardType:          input.AwardType,
		FormStatusID:       1,
		CreatedAt:          now,
		LatestUpdate:       now,
		StudentYear:        input.StudentYear,
		AdvisorName:        input.AdvisorName,
		StudentPhoneNumber: input.StudentPhoneNumber,
		StudentAddress:     input.StudentAddress,
		GPA:                input.GPA,
		StudentDateOfBirth: input.StudentDateOfBirth,
		OrgName:            input.OrgName,
		OrgType:            input.OrgType,
		OrgLocation:        input.OrgLocation,
		OrgPhoneNumber:     input.OrgPhoneNumber,
		FormDetail:         input.FormDetail,
	}

	// 3. เช็ค Role และดึงข้อมูลตาม Role
	// ต้องดึง User เพื่อเช็ค RoleID
	student, studentErr := u.studentService.GetStudentByUserID(ctx, userID)
	org, orgErr := u.organizationService.GetByUserID(ctx, userID)

	if studentErr == nil && student != nil {
		// Role: Student (RoleID = 1)
		form.StudentFirstname = student.User.Firstname
		form.StudentLastname = student.User.Lastname
		form.StudentEmail = student.User.Email
		form.StudentNumber = student.StudentNumber
		form.FacultyID = int(student.FacultyID)
		form.DepartmentID = int(student.DepartmentID)
		form.CampusID = student.User.CampusID
	} else if orgErr == nil && org != nil {
		// Role: Organization (RoleID = 9)
		form.OrgName = org.OrganizationName
		form.OrgType = org.OrganizationType
		form.OrgLocation = org.OrganizationLocation
		form.OrgPhoneNumber = org.OrganizationPhoneNumber
		// CampusID ดึงจาก User
		form.CampusID = org.User.CampusID
	} else {
		return errors.New("user must be either a student or an organization")
	}

	// เรียก Repository โดยส่งไฟล์ (Slice) เข้าไปด้วย
	return u.repo.CreateWithTransaction(ctx, &form, files)
}

func mapToAwardResponse(item models.AwardForm) awardformdto.AwardFormResponse {
	var fileResponses []awardformdto.FileResponse

	// วนลูปแปลงจาก Model ไฟล์ เป็น Response ไฟล์
	for _, f := range item.AwardFiles {
		fileResponses = append(fileResponses, awardformdto.FileResponse{
			FileDirID: f.FileDirID,
			FileType:  f.FileType,
			FileSize:  f.FileSize,
			FilePath:  f.FilePath,
		})
	}

	res := awardformdto.AwardFormResponse{
		FormID:             item.FormID,
		UserID:             item.UserID,
		StudentFirstname:   item.StudentFirstname,
		StudentLastname:    item.StudentLastname,
		StudentEmail:       item.StudentEmail,
		StudentNumber:      item.StudentNumber,
		FacultyID:          item.FacultyID,
		DepartmentID:       item.DepartmentID,
		CampusID:           item.CampusID,
		AcademicYear:       item.AcademicYear,
		Semester:           item.Semester,
		FormStatusID:       item.FormStatusID,
		AwardType:          item.AwardType,
		CreatedAt:          item.CreatedAt,
		LatestUpdate:       item.LatestUpdate,
		StudentYear:        item.StudentYear,
		AdvisorName:        item.AdvisorName,
		StudentPhoneNumber: item.StudentPhoneNumber,
		StudentAddress:     item.StudentAddress,
		GPA:                item.GPA,
		StudentDateOfBirth: item.StudentDateOfBirth,
		OrgName:            item.OrgName,
		OrgType:            item.OrgType,
		OrgLocation:        item.OrgLocation,
		OrgPhoneNumber:     item.OrgPhoneNumber,
		FormDetail:         item.FormDetail,
		RejectReason:       item.RejectReason,
		Files:              fileResponses,
	}

	return res
}

func (u *awardUseCase) GetByKeyword(ctx context.Context, campusID int, keyword string, date string, studentYear int, awardType string, page int, limit int, arrangement string) (*awardformdto.PaginatedAwardResponse, error) {
	// ตั้งค่า default สำหรับ pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	order := "desc"
	if strings.ToLower(arrangement) == "asc" {
		order = "asc"
	}

	results, total, err := u.repo.GetByKeyword(ctx, campusID, keyword, date, studentYear, awardType, page, limit, order)
	if err != nil {
		return nil, err
	}

	var response []awardformdto.AwardFormResponse
	for _, item := range results {
		response = append(response, mapToAwardResponse(item))
	}

	// คำนวณจำนวนหน้า
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &awardformdto.PaginatedAwardResponse{
		Data: response,
		Pagination: awardformdto.PaginationMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalItems:  total,
			Limit:       limit,
		},
	}, nil
}

func (u *awardUseCase) GetAwardsByStudentID(ctx context.Context, studentID int) ([]awardformdto.AwardFormResponse, error) {
	results, err := u.repo.GetByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	var response []awardformdto.AwardFormResponse
	for _, item := range results {
		response = append(response, mapToAwardResponse(item))
	}
	return response, nil
}

func (u *awardUseCase) GetByFormID(ctx context.Context, formID int) (*awardformdto.AwardFormResponse, error) {
	form, err := u.repo.GetByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form == nil {
		return nil, errors.New("form not found")
	}

	response := mapToAwardResponse(*form)
	return &response, nil
}

func (u *awardUseCase) GetAwardsByUserID(ctx context.Context, userID uint) ([]awardformdto.AwardFormResponse, error) {
	results, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var response []awardformdto.AwardFormResponse
	for _, item := range results {
		response = append(response, mapToAwardResponse(item))
	}
	return response, nil
}

func (u *awardUseCase) IsDuplicate(userID uint, year int, semester int) (bool, error) {
	return u.repo.CheckDuplicate(userID, year, semester)
}

func (u *awardUseCase) UpdateAwardType(ctx context.Context, formID uint, awardType string, changedBy uint) error {
	form, err := u.repo.GetByFormID(ctx, int(formID))
	if err != nil {
		return err
	}
	if form == nil {
		return errors.New("form not found")
	}
	if awardType == "" {
		return errors.New("award_type is required")
	}
	if form.AwardType == awardType {
		return nil
	}

	oldValue := form.AwardType
	if err := u.repo.UpdateAwardType(ctx, formID, awardType); err != nil {
		return err
	}

	_, err = u.logUseCase.CreateLog(ctx, changedBy, &awardformdto.CreateAwardFormLogRequest{
		FormID:    formID,
		FieldName: "award_type",
		OldValue:  oldValue,
		NewValue:  awardType,
	})
	return err
}

func (u *awardUseCase) UpdateFormStatus(ctx context.Context, formID uint, formStatus int, changedBy uint) error {
	form, err := u.repo.GetByFormID(ctx, int(formID))
	if err != nil {
		return err
	}
	if form == nil {
		return errors.New("form not found")
	}
	if formStatus == 0 {
		return errors.New("form_status is required")
	}
	if form.FormStatusID == formStatus {
		return nil
	}

	oldValue := form.FormStatusID
	if err := u.repo.UpdateFormStatus(ctx, formID, formStatus); err != nil {
		return err
	}

	_, err = u.logUseCase.CreateLog(ctx, changedBy, &awardformdto.CreateAwardFormLogRequest{
		FormID:    formID,
		FieldName: "form_status",
		OldValue:  strconv.Itoa(oldValue),
		NewValue:  strconv.Itoa(formStatus),
	})
	return err
}
