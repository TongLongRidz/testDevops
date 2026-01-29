package usecase

import (
	awardformdto "backend/internal/dto/award_form_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
)

type AwardUseCase interface {
	// ปรับปรุง: รับ userID เพื่อดึงข้อมูล student และ files
	SubmitAward(ctx context.Context, userID uint, input awardformdto.SubmitAwardRequest, files []models.AwardFileDirectory) error
	GetByKeyword(ctx context.Context, campusID int, keyword string, date string, studentYear int, page int, limit int) (*awardformdto.PaginatedAwardResponse, error)
	GetAwardsByStudentID(ctx context.Context, studentID int) ([]awardformdto.AwardFormResponse, error)
	IsDuplicate(studentID int, year int, semester int) (bool, error)
}

type awardUseCase struct {
	repo                *repository.AwardRepository
	studentService      StudentService
	academicYearService AcademicYearService
}

func NewAwardUseCase(r *repository.AwardRepository, ss StudentService, ays AcademicYearService) AwardUseCase {
	return &awardUseCase{
		repo:                r,
		studentService:      ss,
		academicYearService: ays,
	}
}

func (u *awardUseCase) SubmitAward(ctx context.Context, userID uint, input awardformdto.SubmitAwardRequest, files []models.AwardFileDirectory) error {
	// 1. ดึงข้อมูล Student จาก User ID
	student, err := u.studentService.GetStudentByUserID(ctx, userID)
	if err != nil || student == nil {
		return errors.New("student profile not found")
	}

	// 2. ดึงข้อมูล Academic Year ที่เปิดรับสมัคร (isCurrent = true และ isOpenRegistration = true)
	academicYear, err := u.academicYearService.GetLatestAbleRegister(ctx)
	if err != nil || academicYear == nil {
		return errors.New("no open registration period found")
	}

	// 3. เตรียม Model ตารางหลัก (Award_Form)
	form := models.AwardForm{
		StudentID:        int(student.StudentID),
		StudentFirstname: student.User.Firstname,
		StudentLastname:  student.User.Lastname,
		Email:            student.User.Email,
		StudentNumber:    student.StudentNumber,
		FacultyID:        int(student.FacultyID),
		DepartmentID:     int(student.DepartmentID),
		CampusID:         student.User.CampusID,
		AcademicYear:     academicYear.Year,
		Semester:         academicYear.Semester,
		AwardTypeID:      input.AwardTypeID,
		FormStatusID:     1, // สถานะเริ่มต้น: รอพิจารณา
		StudentYear:      input.StudentYear,
		AdvisorName:      input.AdvisorName,
		PhoneNumber:      input.PhoneNumber,
		Address:          input.Address,
		GPA:              input.GPA,
		DateOfBirth:      input.DateOfBirth,
	}

	var detail interface{}

	// 4. ตรวจสอบประเภทรางวัลตาม ID เพื่อเลือก Model ตารางลูก
	switch input.AwardTypeID {
	case 1: // Extracurricular Activity
		if input.Extracurricular == nil {
			return errors.New("extracurricular detail is required")
		}
		detail = &models.ExtracurricularActivity{
			QualificationType: input.Extracurricular.QualificationType,
			DateReceived:      input.Extracurricular.DateReceived,
			TeamName:          input.Extracurricular.TeamName,
			ProjectTitle:      input.Extracurricular.ProjectTitle,
			Prize:             input.Extracurricular.Prize,
			OrganizedBy:       input.Extracurricular.OrganizedBy,
			CompetitionLevel:  input.Extracurricular.CompetitionLevel,
			ActivityCategory:  input.Extracurricular.ActivityCategory,
			CompetitionName:   input.Extracurricular.CompetitionName,
		}
	case 2: // Good Behavior (ตัวอย่างถ้ามีข้อมูลต้องกรอกให้ใส่เพิ่มแบบ case 1)
		detail = &models.GoodBehavior{}
	case 3: // Creativity & Innovation
		if input.Creativity == nil {
			return errors.New("creativity detail is required")
		}
		detail = &models.CreativityInnovation{
			DateReceived:     input.Creativity.DateReceived,
			TeamName:         input.Creativity.TeamName,
			ProjectTitle:     input.Creativity.ProjectTitle,
			Prize:            input.Creativity.Prize,
			OrganizedBy:      input.Creativity.OrganizedBy,
			CompetitionLevel: input.Creativity.CompetitionLevel,
			ActivityCategory: input.Creativity.ActivityCategory,
			CompetitionName:  input.Creativity.CompetitionName,
		}
	default:
		return errors.New("invalid award type id")
	}

	// 5. เรียก Repository โดยส่งไฟล์ (Slice) เข้าไปด้วย
	// ถ้าไม่มีไฟล์ Slice นี้จะเป็นค่าว่าง ซึ่ง Repo ของเราจัดการได้
	return u.repo.CreateWithTransaction(ctx, &form, detail, files)
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
		FormID:           item.FormID,
		StudentID:        item.StudentID,
		StudentFirstname: item.StudentFirstname,
		StudentLastname:  item.StudentLastname,
		Email:            item.Email,
		StudentNumber:    item.StudentNumber,
		FacultyID:        item.FacultyID,
		DepartmentID:     item.DepartmentID,
		CampusID:         item.CampusID,
		AcademicYear:     item.AcademicYear,
		Semester:         item.Semester,
		FormStatusID:     item.FormStatusID,
		AwardTypeID:      item.AwardTypeID,
		AwardTypeName:    item.AwardType.AwardName,
		CreatedAt:        item.CreatedAt,
		LatestUpdate:     item.LatestUpdate,
		StudentYear:      item.StudentYear,
		AdvisorName:      item.AdvisorName,
		PhoneNumber:      item.PhoneNumber,
		Address:          item.Address,
		GPA:              item.GPA,
		DateOfBirth:      item.DateOfBirth,
		Files:            fileResponses, // ใส่ไฟล์ลงไปที่นี่
	}

	if item.Extracurricular != nil {
		res.Detail = item.Extracurricular
	} else if item.GoodBehavior != nil {
		res.Detail = item.GoodBehavior
	} else if item.Creativity != nil {
		res.Detail = item.Creativity
	}
	return res
}

func (u *awardUseCase) GetByKeyword(ctx context.Context, campusID int, keyword string, date string, studentYear int, page int, limit int) (*awardformdto.PaginatedAwardResponse, error) {
	// ตั้งค่า default สำหรับ pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	results, total, err := u.repo.GetByKeyword(ctx, campusID, keyword, date, studentYear, page, limit)
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

func (u *awardUseCase) IsDuplicate(studentID int, year int, semester int) (bool, error) {
	return u.repo.CheckDuplicate(studentID, year, semester)
}
