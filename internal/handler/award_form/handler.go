package awardform

import (
	awardformdto "backend/internal/dto/award_form_dto"
	"backend/internal/models"
	"backend/internal/usecase"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AwardHandler struct {
	useCase             usecase.AwardUseCase
	logUseCase          usecase.AwardFormLogUseCase
	studentService      usecase.StudentService
	academicYearService usecase.AcademicYearService
}

func NewAwardHandler(u usecase.AwardUseCase, s usecase.StudentService, ays usecase.AcademicYearService, l usecase.AwardFormLogUseCase) *AwardHandler {
	return &AwardHandler{useCase: u, logUseCase: l, studentService: s, academicYearService: ays}
}

func (h *AwardHandler) Submit(c *fiber.Ctx) error {
	// สร้าง uploads folder อัตโนมัติ
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create uploads directory",
		})
	}

	// ดึงข้อมูลผู้ใช้ที่ login อยู่จาก middleware
	currentUser := c.Locals("current_user")
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: User not found",
		})
	}
	user, ok := currentUser.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data",
		})
	}

	// 1. Check Role และรับข้อมูลตามแต่ละ Role
	fmt.Printf("=== DEBUG: User RoleID = %d ===\n", user.RoleID)

	var req awardformdto.SubmitAwardRequest

	// ===== ROLE: STUDENT (RoleID = 1) =====
	if user.RoleID == 1 {
		fmt.Println("🎓 Processing STUDENT submission...")

		// Student กรอก:
		awardType := c.FormValue("award_type")
		if awardType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "award_type is required",
			})
		}
		req.AwardType = awardType

		studentYear, err := strconv.Atoi(c.FormValue("student_year"))
		if err != nil || studentYear == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_year is required and must be a valid number",
			})
		}
		req.StudentYear = studentYear

		advisorName := c.FormValue("advisor_name")
		if advisorName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "advisor_name is required",
			})
		}
		req.AdvisorName = advisorName

		studentPhoneNumber := c.FormValue("student_phone_number")
		if studentPhoneNumber == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_phone_number is required",
			})
		}
		req.StudentPhoneNumber = studentPhoneNumber

		studentAddress := c.FormValue("student_address")
		if studentAddress == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_address is required",
			})
		}
		req.StudentAddress = studentAddress

		gpa, err := strconv.ParseFloat(c.FormValue("gpa"), 64)
		if err != nil || gpa < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "gpa is required and must be a valid number",
			})
		}
		req.GPA = gpa

		dobStr := c.FormValue("student_date_of_birth")
		if dobStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_date_of_birth is required",
			})
		}
		dob, err := time.Parse("2006-01-02", dobStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_date_of_birth format should be YYYY-MM-DD",
			})
		}
		req.StudentDateOfBirth = dob

		formDetail := c.FormValue("form_detail")
		if formDetail == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "form_detail is required",
			})
		}
		req.FormDetail = formDetail

		// ===== ROLE: ORGANIZATION (RoleID = 9) =====
	} else if user.RoleID == 9 {
		fmt.Println("🏢 Processing ORGANIZATION submission...")

		// Organization กรอก:
		studentFirstname := c.FormValue("student_firstname")
		if studentFirstname == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_firstname is required",
			})
		}
		req.StudentFirstname = studentFirstname

		studentLastname := c.FormValue("student_lastname")
		if studentLastname == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_lastname is required",
			})
		}
		req.StudentLastname = studentLastname

		studentEmail := c.FormValue("student_email")
		if studentEmail == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_email is required",
			})
		}
		req.StudentEmail = studentEmail

		studentNumber := c.FormValue("student_number")
		if studentNumber == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_number is required",
			})
		}
		req.StudentNumber = studentNumber

		facultyID, err := strconv.Atoi(c.FormValue("faculty_id"))
		if err != nil || facultyID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "faculty_id is required and must be a valid number",
			})
		}
		req.FacultyID = facultyID

		departmentID, err := strconv.Atoi(c.FormValue("department_id"))
		if err != nil || departmentID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "department_id is required and must be a valid number",
			})
		}
		req.DepartmentID = departmentID

		awardType := c.FormValue("award_type")
		if awardType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "award_type is required",
			})
		}
		req.AwardType = awardType

		studentYear, err := strconv.Atoi(c.FormValue("student_year"))
		if err != nil || studentYear == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_year is required and must be a valid number",
			})
		}
		req.StudentYear = studentYear

		advisorName := c.FormValue("advisor_name")
		if advisorName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "advisor_name is required",
			})
		}
		req.AdvisorName = advisorName

		studentPhoneNumber := c.FormValue("student_phone_number")
		if studentPhoneNumber == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_phone_number is required",
			})
		}
		req.StudentPhoneNumber = studentPhoneNumber

		studentAddress := c.FormValue("student_address")
		if studentAddress == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_address is required",
			})
		}
		req.StudentAddress = studentAddress

		gpa, err := strconv.ParseFloat(c.FormValue("gpa"), 64)
		if err != nil || gpa < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "gpa is required and must be a valid number",
			})
		}
		req.GPA = gpa

		dobStr := c.FormValue("student_date_of_birth")
		if dobStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_date_of_birth is required",
			})
		}
		dob, err := time.Parse("2006-01-02", dobStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "student_date_of_birth format should be YYYY-MM-DD",
			})
		}
		req.StudentDateOfBirth = dob

		formDetail := c.FormValue("form_detail")
		if formDetail == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "form_detail is required",
			})
		}
		req.FormDetail = formDetail

	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Only Student (RoleID=1) and Organization (RoleID=9) can submit awards",
		})
	}

	// จัดการกับไฟล์แนบ (ถ้ามี)
	var awardFiles []models.AwardFileDirectory

	form, err := c.MultipartForm()
	if err == nil {
		// Debug: แสดง field names ทั้งหมดที่มีในฟอร์ม
		fmt.Println("🔍 Form fields ที่มีในฟอร์ม:")
		for fieldName := range form.File {
			fmt.Printf("  - %s: %d files\n", fieldName, len(form.File[fieldName]))
		}

		files := form.File["files"]

		fmt.Printf("📁 จำนวนไฟล์ที่ได้รับ (field 'files'): %d\n", len(files))

		// --- STEP 1: VALIDATION LOOP ---
		// เช็คไฟล์ทั้งหมดก่อนว่ามีอันไหนไม่ valid ไหม
		allowedExtensions := map[string]bool{".pdf": true}
		maxTotalSize := int64(10 * 1024 * 1024) // 10 MB
		var totalSize int64

		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if !allowedExtensions[ext] {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": fmt.Sprintf("ไม่อนุญาตให้อัปโหลดไฟล์ประเภท %s (รองรับเฉพาะ PDF)", ext),
				})
			}
			totalSize += file.Size
		}

		// เช็คขนาดไฟล์รวม
		if totalSize > maxTotalSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("ขนาดไฟล์รวมเกิน 10 MB (ได้รับ %.2f MB)", float64(totalSize)/(1024*1024)),
			})
		}

		// --- STEP 2: PROCESSING & SAVING LOOP ---
		// ถ้าผ่านการเช็คด้านบนมาได้ แสดงว่าไฟล์ทุกไฟล์ valid พร้อมบันทึก
		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			subDir := "pdf" // รองรับเฉพาะ PDF

			targetDir := filepath.Join(uploadDir, subDir)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": "error", "message": "Failed to create directory",
				})
			}

			newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			savePath := filepath.Join(targetDir, newFileName)

			if err := c.SaveFile(file, savePath); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to save file: " + err.Error(),
				})
			}

			fmt.Printf("✅ บันทึกไฟล์สำเร็จ: %s (ขนาด: %d bytes)\n", savePath, file.Size)

			cleanPath := filepath.ToSlash(savePath)
			awardFiles = append(awardFiles, models.AwardFileDirectory{
				FilePath:   cleanPath,
				FileType:   strings.TrimPrefix(ext, "."),
				FileSize:   file.Size,
				UploadedAt: time.Now(),
			})
		}
	}

	// 3. ส่งข้อมูลไปยัง UseCase พร้อม userID
	if err := h.useCase.SubmitAward(c.UserContext(), user.UserID, req, awardFiles); err != nil {

		// --- ส่วนที่เพิ่มเข้ามา: ลบไฟล์ทิ้งถ้า DB บันทึกไม่สำเร็จ ---
		for _, f := range awardFiles {
			// f.FilePath เก็บค่าเช่น "uploads/pdf/xxx.pdf"
			if removeErr := os.Remove(f.FilePath); removeErr != nil {
				fmt.Printf("Failed to cleanup file %s: %v\n", f.FilePath, removeErr)
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "บันทึกข้อมูลไม่สำเร็จ (อาจมีการส่งข้อมูลในปีการศึกษานี้ไปแล้ว): " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Award form submitted successfully",
	})
}

// GetByKeyword ค้นหาและกรองตามเงื่อนไข พร้อม pagination
// Query params: keyword, date (YYYY-MM-DD), student_year, page (default: 1), limit (default: 10)
func (h *AwardHandler) GetByKeyword(c *fiber.Ctx) error {
	// ดึงข้อมูล user จาก middleware
	currentUser := c.Locals("current_user")
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: User not found",
		})
	}
	user, ok := currentUser.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data",
		})
	}

	// รับ query parameters
	var req awardformdto.SearchAwardRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid query parameters",
		})
	}

	// ค้นหาและกรองตามวิทยาเขตของ user
	results, err := h.useCase.GetByKeyword(
		c.UserContext(),
		user.CampusID,
		req.Keyword,
		req.Date,
		req.StudentYear,
		req.AwardType,
		req.Page,
		req.Limit,
		req.Arrangement,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"data":       results.Data,
		"pagination": results.Pagination,
	})
}

func (h *AwardHandler) GetMySubmissions(c *fiber.Ctx) error {
	// ดึงข้อมูล user จาก middleware
	currentUser := c.Locals("current_user")
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: User not found",
		})
	}
	user, ok := currentUser.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data",
		})
	}

	var err error
	var pagedResults *awardformdto.PaginatedAwardResponse
	page := 1
	limit := 4

	yearQuery := c.Query("year")
	if yearQuery == "" {
		yearQuery = c.Query("years")
	}

	pageQuery := c.Query("page")
	limitQuery := c.Query("limit")
	if pageQuery != "" {
		pageValue, convErr := strconv.Atoi(pageQuery)
		if convErr != nil || pageValue <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid page parameter",
			})
		}
		page = pageValue
	}
	if limitQuery != "" {
		limitValue, convErr := strconv.Atoi(limitQuery)
		if convErr != nil || limitValue <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid limit parameter",
			})
		}
		limit = limitValue
	}

	var yearList []int
	if yearQuery != "" {
		parts := strings.Split(yearQuery, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			yearValue, convErr := strconv.Atoi(part)
			if convErr != nil || yearValue <= 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Invalid year parameter",
				})
			}
			yearList = append(yearList, yearValue)
		}
		if len(yearList) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "year is required",
			})
		}
	}

	// เช็ค Role และดึงข้อมูลตาม Role
	switch user.RoleID {
	case 1: // Student
		// ดึงการส่งฟอร์มของนักเรียนนี้ (sorted by created_at desc)
		if len(yearList) > 0 {
			pagedResults, err = h.useCase.GetAwardsByUserIDPaged(c.UserContext(), user.UserID, yearList, page, limit)
		} else {
			pagedResults, err = h.useCase.GetAwardsByUserIDPaged(c.UserContext(), user.UserID, nil, page, limit)
		}

	case 9: // Organization
		// ดึงการส่งฟอร์มของ organization นี้ (sorted by created_at desc)
		if len(yearList) > 0 {
			pagedResults, err = h.useCase.GetAwardsByUserIDPaged(c.UserContext(), user.UserID, yearList, page, limit)
		} else {
			pagedResults, err = h.useCase.GetAwardsByUserIDPaged(c.UserContext(), user.UserID, nil, page, limit)
		}

	default:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Only Student and Organization can view submissions",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"data":       pagedResults.Data,
		"pagination": pagedResults.Pagination,
	})
}

func (h *AwardHandler) GetMyCurrentSemesterSubmissions(c *fiber.Ctx) error {
	// ดึงข้อมูล user จาก middleware
	currentUser := c.Locals("current_user")
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: User not found",
		})
	}
	user, ok := currentUser.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data",
		})
	}

	// เช็ค Role
	if user.RoleID != 1 && user.RoleID != 9 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Only Student and Organization can view submissions",
		})
	}

	// ดึง Academic Year ปัจจุบันที่เปิดใช้งาน (isActive = true)
	currentSemester, err := h.academicYearService.GetCurrentSemester(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get current semester: " + err.Error(),
		})
	}

	if currentSemester == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No active semester found",
		})
	}

	// ดึงการส่งฟอร์มของ user ในภาคเรียนปัจจุบัน
	results, err := h.useCase.GetAwardsByUserIDAndSemester(c.UserContext(), user.UserID, int(currentSemester.Year), int(currentSemester.Semester))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if results == nil {
		results = []awardformdto.AwardFormResponse{}
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   results,
		"meta": fiber.Map{
			"academic_year": currentSemester.Year,
			"semester":      currentSemester.Semester,
		},
	})
}

func (h *AwardHandler) CreateLog(c *fiber.Ctx) error {
	currentUser := c.Locals("current_user")
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: User not found",
		})
	}
	user, ok := currentUser.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data",
		})
	}

	formID, err := strconv.Atoi(c.Params("formId"))
	if err != nil || formID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid formId",
		})
	}

	var req awardformdto.CreateAwardFormLogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	if req.FieldName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "field_name is required",
		})
	}
	req.FormID = uint(formID)

	log, err := h.logUseCase.CreateLog(c.UserContext(), user.UserID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data": awardformdto.AwardFormLogResponse{
			LogID:     log.LogID,
			FormID:    log.FormID,
			FieldName: log.FieldName,
			OldValue:  log.OldValue,
			NewValue:  log.NewValue,
			ChangedBy: log.ChangedBy,
			CreatedAt: log.CreatedAt,
		},
	})
}

func (h *AwardHandler) GetLogsByFormID(c *fiber.Ctx) error {
	formID, err := strconv.Atoi(c.Params("formId"))
	if err != nil || formID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid formId",
		})
	}

	logs, err := h.logUseCase.GetLogsByFormID(c.UserContext(), uint(formID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	response := make([]awardformdto.AwardFormLogResponse, 0, len(logs))
	for _, log := range logs {
		response = append(response, awardformdto.AwardFormLogResponse{
			LogID:     log.LogID,
			FormID:    log.FormID,
			FieldName: log.FieldName,
			OldValue:  log.OldValue,
			NewValue:  log.NewValue,
			ChangedBy: log.ChangedBy,
			CreatedAt: log.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func (h *AwardHandler) GetByFormID(c *fiber.Ctx) error {
	formID, err := strconv.Atoi(c.Params("formId"))
	if err != nil || formID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid formId",
		})
	}

	form, err := h.useCase.GetByFormID(c.UserContext(), formID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   form,
	})
}

func (h *AwardHandler) UpdateAwardType(c *fiber.Ctx) error {
	currentUser := c.Locals("current_user")
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: User not found",
		})
	}
	user, ok := currentUser.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data",
		})
	}

	formID, err := strconv.Atoi(c.Params("formId"))
	if err != nil || formID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid formId",
		})
	}

	var req awardformdto.UpdateAwardTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	if req.AwardType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "award_type is required",
		})
	}

	if err := h.useCase.UpdateAwardType(c.UserContext(), uint(formID), req.AwardType, user.UserID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "award_type updated",
	})
}

func (h *AwardHandler) UpdateFormStatus(c *fiber.Ctx) error {
	currentUser := c.Locals("current_user")
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: User not found",
		})
	}
	user, ok := currentUser.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data",
		})
	}

	formID, err := strconv.Atoi(c.Params("formId"))
	if err != nil || formID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid formId",
		})
	}

	var req awardformdto.UpdateFormStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	if req.FormStatusID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "form_status is required",
		})
	}

	if err := h.useCase.UpdateFormStatus(c.UserContext(), uint(formID), req.FormStatusID, user.UserID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "form_status updated",
	})
}

// GetAllAwardTypes - ดึง award_type ทั้งหมดที่มีในระบบ
func (h *AwardHandler) GetAllAwardTypes(c *fiber.Ctx) error {
	awardTypes, err := h.useCase.GetAllAwardTypes(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch award types",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   awardTypes,
	})
}
