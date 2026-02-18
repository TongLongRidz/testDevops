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
	var req awardformdto.SubmitAwardRequest

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

	// 1. รับข้อมูลที่เป็น Text/JSON จาก Form
	// Debug: ดูค่าที่ได้รับทั้งหมด
	fmt.Println("=== DEBUG FORM VALUES ===")
	fmt.Printf("award_type: '%s'\n", c.FormValue("award_type"))
	fmt.Printf("student_year: '%s'\n", c.FormValue("student_year"))
	fmt.Printf("advisor_name: '%s'\n", c.FormValue("advisor_name"))

	awardType := c.FormValue("award_type")
	if awardType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "award_type is required",
		})
	}
	req.AwardType = awardType

	// ดึงข้อมูลทั่วไป (ที่ไม่ได้มาจาก token)
	studentYear, _ := strconv.Atoi(c.FormValue("student_year"))
	req.StudentYear = studentYear
	req.AdvisorName = c.FormValue("advisor_name")
	req.StudentPhoneNumber = c.FormValue("student_phone_number")
	req.StudentAddress = c.FormValue("student_address")
	gpa, _ := strconv.ParseFloat(c.FormValue("gpa"), 64)
	req.GPA = gpa
	if dobStr := c.FormValue("student_date_of_birth"); dobStr != "" {
		dob, _ := time.Parse("2006-01-02", dobStr)
		req.StudentDateOfBirth = dob
	}

	// Organization Information
	req.OrgName = c.FormValue("org_name")
	req.OrgType = c.FormValue("org_type")
	req.OrgLocation = c.FormValue("org_location")
	req.OrgPhoneNumber = c.FormValue("org_phone_number")

	// Form Detail
	req.FormDetail = c.FormValue("form_detail")

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

	// ดึงข้อมูล student จาก userID
	student, err := h.studentService.GetStudentByUserID(c.UserContext(), user.UserID)
	if err != nil || student == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Student profile not found",
		})
	}

	// ดึงการส่งฟอร์มของนักเรียนนี้ (sorted by created_at desc)
	results, err := h.useCase.GetAwardsByStudentID(c.UserContext(), int(student.StudentID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   results,
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
