package usecase

import (
	awardformdto "backend/internal/dto/award_form_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"strings"
	"time"
)

type AwardUseCase interface {
	// ปรับปรุง: รับ userID เพื่อดึงข้อมูล student และ files
	SubmitAward(ctx context.Context, userID uint, input awardformdto.SubmitAwardRequest, files []models.AwardFileDirectory) error
	GetByKeyword(ctx context.Context, userID uint, roleID int, campusID int, keyword string, date string, studentYear int, awardType string, sortBy string, sortOrder string, page int, limit int) (*awardformdto.PaginatedAwardResponse, error)
	GetAwardsByUserID(ctx context.Context, userID uint) ([]awardformdto.AwardFormResponse, error)
	GetAwardsByStudentID(ctx context.Context, studentID int) ([]awardformdto.AwardFormResponse, error)
	GetAwardsByUserIDAndYear(ctx context.Context, userID uint, year int) ([]awardformdto.AwardFormResponse, error)
	GetAwardsByUserIDAndSemester(ctx context.Context, userID uint, year int, semester int) ([]awardformdto.AwardFormResponse, error)
	GetAwardsByUserIDWithYearSort(ctx context.Context, userID uint, years []int) ([]awardformdto.AwardFormResponse, error)
	GetAwardsByUserIDPaged(ctx context.Context, userID uint, years []int, page int, limit int) (*awardformdto.PaginatedAwardResponse, error)
	GetByFormID(ctx context.Context, formID int) (*awardformdto.AwardFormResponse, error)
	IsDuplicate(userID uint, year int, semester int) (bool, error)
	UpdateAwardType(ctx context.Context, formID uint, awardType string, changedBy uint) error
	UpdateFormStatus(ctx context.Context, formID uint, formStatus int, rejectReason string, changedBy uint) error
	UpdateFormStatusWithLog(ctx context.Context, formID uint, formStatus int, rejectReason string, changedBy uint) error
	CommitteeVote(ctx context.Context, formID uint, operation string, votedBy uint) (*awardformdto.CommitteeVoteResult, error)
	GetApprovalLogsByUserID(ctx context.Context, userID uint) ([]models.AwardApprovalLog, error)
	GetApprovalHistory(ctx context.Context, userID uint, campusID int, keyword string, date string, awardType string, operation string, sortBy string, sortOrder string, page int, limit int) (*awardformdto.PaginatedApprovalLogResponse, error)
	GetAllAwardTypes(ctx context.Context) ([]string, error)
	GetApprovalLogDetail(ctx context.Context, approvalLogID uint) (*models.AwardApprovalLog, error)
}

type awardUseCase struct {
	repo                *repository.AwardRepository
	studentService      StudentService
	organizationService OrganizationService
	academicYearService AcademicYearService
}

func NewAwardUseCase(r *repository.AwardRepository, ss StudentService, os OrganizationService, ays AcademicYearService) AwardUseCase {
	return &awardUseCase{
		repo:                r,
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

	// 2. เตรียม Model ตารางหลัก (Award_Form) - พื้นฐาน
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
		FormDetail:         input.FormDetail,
	}

	// 3. เช็ค Role และดึงข้อมูลตาม Role
	student, studentErr := u.studentService.GetStudentByUserID(ctx, userID)
	org, orgErr := u.organizationService.GetByUserID(ctx, userID)

	// ===== ROLE: STUDENT (RoleID = 1) =====
	if studentErr == nil && student != nil {
		// ใช้ค่าที่ auto-fill จาก token ใน handler
		form.StudentFirstname = input.StudentFirstname
		form.StudentLastname = input.StudentLastname
		form.StudentEmail = input.StudentEmail
		form.StudentNumber = student.StudentNumber
		form.FacultyID = int(student.FacultyID)
		form.DepartmentID = int(student.DepartmentID)
		form.CampusID = student.User.CampusID
		// Organization info ไม่ต้องกรอก
		form.OrgName = ""
		form.OrgType = ""
		form.OrgLocation = ""
		form.OrgPhoneNumber = ""

		// ===== ROLE: ORGANIZATION (RoleID = 8) =====
	} else if orgErr == nil && org != nil {
		// ดึงจาก Organization Service:
		form.OrgName = org.OrganizationName
		form.OrgType = org.OrganizationType
		form.OrgLocation = org.OrganizationLocation
		form.OrgPhoneNumber = org.OrganizationPhoneNumber
		form.CampusID = org.User.CampusID

		// ใช้ข้อมูลที่ Organization กรอก:
		form.StudentFirstname = input.StudentFirstname
		form.StudentLastname = input.StudentLastname
		form.StudentEmail = input.StudentEmail
		form.StudentNumber = input.StudentNumber
		form.FacultyID = input.FacultyID
		form.DepartmentID = input.DepartmentID

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

func (u *awardUseCase) GetByKeyword(ctx context.Context, userID uint, roleID int, campusID int, keyword string, date string, studentYear int, awardType string, sortBy string, sortOrder string, page int, limit int) (*awardformdto.PaginatedAwardResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5
	}
	if limit > 5 {
		limit = 5
	}

	normalizedSortBy := normalizeSortBy(sortBy)
	normalizedSortOrder := normalizeSortOrder(sortOrder)

	filter := repository.AwardSearchFilter{
		CampusID:    campusID,
		Keyword:     strings.TrimSpace(keyword),
		Date:        strings.TrimSpace(date),
		StudentYear: studentYear,
		AwardType:   strings.TrimSpace(awardType),
		SortBy:      normalizedSortBy,
		SortOrder:   normalizedSortOrder,
		Page:        page,
		Limit:       limit,
	}

	if formStatusID, hasRequiredStatus := requiredFormStatusByRole(roleID); hasRequiredStatus {
		filter.FormStatusID = &formStatusID
	}

	if roleID == 6 {
		filter.ExcludeVotedByUserID = &userID
	}

	if needsDepartmentScopeByRole(roleID) {
		facultyID, departmentID, err := u.repo.GetHeadOfDepartmentScopeByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if departmentID == 0 {
			return nil, errors.New("unable to resolve head of department scope")
		}
		filter.DepartmentID = &departmentID
		_ = facultyID
	}

	if needsFacultyScopeByRole(roleID) {
		facultyID, err := u.repo.GetFacultyScopeByRoleAndUserID(ctx, roleID, userID)
		if err != nil {
			return nil, err
		}
		if facultyID == 0 {
			return nil, errors.New("unable to resolve faculty scope")
		}
		filter.FacultyID = &facultyID
	}

	results, total, err := u.repo.GetByKeyword(ctx, filter)
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

func normalizeSortBy(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "name":
		return "name"
	case "studentnumber":
		return "studentNumber"
	case "academicyear":
		return "academicYear"
	case "awardtype":
		return "awardType"
	case "date", "":
		return "date"
	default:
		return "date"
	}
}

func normalizeSortOrder(sortOrder string) string {
	switch strings.ToLower(strings.TrimSpace(sortOrder)) {
	case "asc":
		return "asc"
	case "desc", "des", "":
		return "desc"
	default:
		return "desc"
	}
}

func normalizeOperation(operation string) string {
	switch strings.TrimSpace(operation) {
	case "", "ทั้งหมด":
		return ""
	case "อนุมัติ", "approve", "approved":
		return "approve"
	case "ไม่อนุมัติ", "ปฏิเสธ", "reject", "rejected", "ตีกลับ", "return", "returned":
		return "reject"
	default:
		return strings.TrimSpace(operation)
	}
}

func requiredFormStatusByRole(roleID int) (int, bool) {
	switch roleID {
	case 2:
		return 1, true
	case 3:
		return 2, true
	case 4:
		return 4, true
	case 5:
		return 6, true
	case 6:
		return 8, true
	case 7:
		return 12, true
	default:
		return 0, false
	}
}

func needsDepartmentScopeByRole(roleID int) bool {
	return roleID == 2
}

func needsFacultyScopeByRole(roleID int) bool {
	return roleID == 3 || roleID == 4
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

func (u *awardUseCase) GetAwardsByUserIDAndYear(ctx context.Context, userID uint, year int) ([]awardformdto.AwardFormResponse, error) {
	results, err := u.repo.GetByUserIDAndYear(ctx, userID, year)
	if err != nil {
		return nil, err
	}
	var response []awardformdto.AwardFormResponse
	for _, item := range results {
		response = append(response, mapToAwardResponse(item))
	}
	return response, nil
}

func (u *awardUseCase) GetAwardsByUserIDAndSemester(ctx context.Context, userID uint, year int, semester int) ([]awardformdto.AwardFormResponse, error) {
	results, err := u.repo.GetByUserIDAndSemester(ctx, userID, year, semester)
	if err != nil {
		return nil, err
	}
	var response []awardformdto.AwardFormResponse
	for _, item := range results {
		response = append(response, mapToAwardResponse(item))
	}
	return response, nil
}

func (u *awardUseCase) GetAwardsByUserIDWithYearSort(ctx context.Context, userID uint, years []int) ([]awardformdto.AwardFormResponse, error) {
	results, err := u.repo.GetByUserIDWithYearSort(ctx, userID, years)
	if err != nil {
		return nil, err
	}
	var response []awardformdto.AwardFormResponse
	for _, item := range results {
		response = append(response, mapToAwardResponse(item))
	}
	return response, nil
}

func (u *awardUseCase) GetAwardsByUserIDPaged(ctx context.Context, userID uint, years []int, page int, limit int) (*awardformdto.PaginatedAwardResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	results, total, err := u.repo.GetByUserIDWithYearPaged(ctx, userID, years, page, limit)
	if err != nil {
		return nil, err
	}

	var response []awardformdto.AwardFormResponse
	for _, item := range results {
		response = append(response, mapToAwardResponse(item))
	}

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

	if err := u.repo.UpdateAwardType(ctx, formID, awardType); err != nil {
		return err
	}
	return nil
}

func (u *awardUseCase) UpdateFormStatus(ctx context.Context, formID uint, formStatus int, rejectReason string, changedBy uint) error {
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

	trimmedRejectReason := strings.TrimSpace(rejectReason)
	if isRejectOrReturnStatus(formStatus) && trimmedRejectReason == "" {
		return errors.New("reject_reason is required for reject or return status")
	}
	if !isRejectOrReturnStatus(formStatus) {
		trimmedRejectReason = ""
	}

	if err := u.repo.UpdateFormStatus(ctx, formID, formStatus, trimmedRejectReason); err != nil {
		return err
	}

	if approvalStatus, shouldLogApproval := mapFormStatusToApprovalStatus(formStatus); shouldLogApproval {
		log := &models.AwardApprovalLog{
			FormID:         formID,
			UserID:         changedBy,
			ApprovalStatus: approvalStatus,
			RejectReason:   trimmedRejectReason,
			ApprovedAt:     time.Now(),
		}
		if err := u.repo.CreateAwardApprovalLog(ctx, log); err != nil {
			return err
		}
	}

	return nil
}

// UpdateFormStatusWithLog - สำหรับ role 5 (Student Development) บันทึกประวัติการอนุมัติ/ตีกลับ
func (u *awardUseCase) UpdateFormStatusWithLog(ctx context.Context, formID uint, formStatus int, rejectReason string, changedBy uint) error {
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

	trimmedRejectReason := strings.TrimSpace(rejectReason)

	// สำหรับการอนุมัติ (status = 2) ไม่ต้องมี reject_reason
	// สำหรับการตีกลับ (status = 3) ต้องมี reject_reason
	if formStatus == 3 && trimmedRejectReason == "" {
		return errors.New("reject_reason is required for rejection")
	}

	// 1. อัปเดตสถานะฟอร์ม
	if err := u.repo.UpdateFormStatus(ctx, formID, formStatus, trimmedRejectReason); err != nil {
		return err
	}

	// 2. Map formStatus เป็น status string
	statusMap := map[int]string{
		2: "approve",
		3: "reject",
	}

	newStatus := statusMap[formStatus]

	// 3. กำหนด logType
	logType := "approval"
	if formStatus == 3 {
		logType = "rejection"
	}

	// 4. บันทึก award type log
	typeLog := &models.AwardTypeLog{
		FormID:    formID,
		UserID:    changedBy,
		LogType:   logType,
		Status:    newStatus,
		ChangedAt: time.Now(),
	}

	// สำหรับการตีกลับ บันทึก old value (award type เดิม) และ reason
	if formStatus == 3 {
		typeLog.OldValue = form.AwardType
		typeLog.RejectReason = trimmedRejectReason
	}

	if err := u.repo.SaveAwardTypeLog(ctx, typeLog); err != nil {
		return err
	}

	if shouldLogSigned := shouldCreateSignedLog(formStatus); shouldLogSigned {
		signedLog := &models.AwardSignedLog{
			FormID:   formID,
			UserID:   changedBy,
			SignedAt: time.Now(),
		}
		if err := u.repo.CreateAwardSignedLog(ctx, signedLog); err != nil {
			return err
		}
	}

	return nil
}

func (u *awardUseCase) CommitteeVote(ctx context.Context, formID uint, operation string, votedBy uint) (*awardformdto.CommitteeVoteResult, error) {
	form, err := u.repo.GetByFormID(ctx, int(formID))
	if err != nil {
		return nil, err
	}
	if form == nil {
		return nil, errors.New("form not found")
	}

	isEligible, err := u.repo.IsCommitteeNonChairman(ctx, votedBy)
	if err != nil {
		return nil, err
	}
	if !isEligible {
		return nil, errors.New("only committee members with is_chairman=false can vote")
	}

	normalized := strings.ToLower(strings.TrimSpace(operation))
	if normalized == "approve" {
		normalized = "approve"
	}
	if normalized == "approved" {
		normalized = "approve"
	}
	if normalized == "rejected" {
		normalized = "reject"
	}
	if normalized != "approve" && normalized != "reject" {
		return nil, errors.New("operation must be approve or reject")
	}

	if err := u.repo.UpsertCommitteeVoteLog(ctx, formID, votedBy, normalized); err != nil {
		return nil, err
	}

	totalCommittee, err := u.repo.CountNonChairmanCommittees(ctx)
	if err != nil {
		return nil, err
	}
	if totalCommittee == 0 {
		return nil, errors.New("no committee members available for voting")
	}

	approvedCount, err := u.repo.CountCommitteeVotesByOperation(ctx, formID, "approve")
	if err != nil {
		return nil, err
	}
	rejectCount, err := u.repo.CountCommitteeVotesByOperation(ctx, formID, "reject")
	if err != nil {
		return nil, err
	}

	half := totalCommittee / 2
	hasApproveMajority := approvedCount > half
	hasRejectMajority := rejectCount > half
	hasMajority := hasApproveMajority || hasRejectMajority
	currentFormStatusID := form.FormStatusID
	if hasApproveMajority {
		if err := u.repo.UpdateFormStatus(ctx, formID, 10, ""); err != nil {
			return nil, err
		}
		currentFormStatusID = 10
	} else if hasRejectMajority {
		if err := u.repo.UpdateFormStatus(ctx, formID, 11, "คณะกรรมการไม่เห็นชอบ"); err != nil {
			return nil, err
		}
		currentFormStatusID = 11
	}

	return &awardformdto.CommitteeVoteResult{
		Operation:      normalized,
		ApproveCount:   approvedCount,
		RejectCount:    rejectCount,
		TotalVoters:    totalCommittee,
		VotedCount:     approvedCount + rejectCount,
		HasMajority:    hasMajority,
		MajorityTarget: (totalCommittee / 2) + 1,
		FormStatusID:   currentFormStatusID,
	}, nil
}

func mapFormStatusToApprovalStatus(formStatus int) (string, bool) {
	switch formStatus {
	case 2, 4, 6, 8, 10:
		return "approve", true
	case 3, 5, 7, 11:
		return "reject", true
	case 9:
		return "reject", true
	default:
		return "", false
	}
}

func shouldCreateSignedLog(formStatus int) bool {
	switch formStatus {
	case 12, 13:
		return true
	default:
		return false
	}
}

func isRejectOrReturnStatus(formStatus int) bool {
	switch formStatus {
	case 3, 5, 7, 9, 11:
		return true
	default:
		return false
	}
}

func (u *awardUseCase) GetAllAwardTypes(ctx context.Context) ([]string, error) {
	return u.repo.GetAllAwardTypes(ctx)
}

func (u *awardUseCase) GetApprovalHistory(ctx context.Context, userID uint, campusID int, keyword string, date string, awardType string, operation string, sortBy string, sortOrder string, page int, limit int) (*awardformdto.PaginatedApprovalLogResponse, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	if campusID == 0 {
		return nil, errors.New("invalid campus id")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5
	}
	if limit > 5 {
		limit = 5
	}

	filter := repository.ApprovalLogSearchFilter{
		UserID:    userID,
		CampusID:  campusID,
		Keyword:   strings.TrimSpace(keyword),
		Date:      strings.TrimSpace(date),
		AwardType: strings.TrimSpace(awardType),
		Operation: normalizeOperation(operation),
		SortBy:    normalizeSortBy(sortBy),
		SortOrder: normalizeSortOrder(sortOrder),
		Page:      page,
		Limit:     limit,
	}

	rows, total, err := u.repo.GetApprovalHistoryByUserAndCampus(ctx, filter)
	if err != nil {
		return nil, err
	}

	responseData := make([]awardformdto.ApprovalLogHistoryResponse, 0, len(rows))
	for _, row := range rows {
		responseData = append(responseData, awardformdto.ApprovalLogHistoryResponse{
			ApprovalLogID:    row.ApprovalLogID,
			FormID:           row.FormID,
			ReviewerUserID:   row.ReviewerUserID,
			Operation:        row.Operation,
			OperationDate:    row.OperationDate,
			StudentFirstname: row.StudentFirstname,
			StudentLastname:  row.StudentLastname,
			StudentNumber:    row.StudentNumber,
			AcademicYear:     row.AcademicYear,
			AwardType:        row.AwardType,
			CampusID:         row.CampusID,
		})
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &awardformdto.PaginatedApprovalLogResponse{
		Data: responseData,
		Pagination: awardformdto.PaginationMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalItems:  total,
			Limit:       limit,
		},
	}, nil
}

func (u *awardUseCase) GetApprovalLogsByUserID(ctx context.Context, userID uint) ([]models.AwardApprovalLog, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}
	return u.repo.GetApprovalLogsByUserID(ctx, userID)
}

// GetApprovalLogDetail implements AwardUseCase for fetching approval log detail by ID
func (u *awardUseCase) GetApprovalLogDetail(ctx context.Context, approvalLogID uint) (*models.AwardApprovalLog, error) {
	if approvalLogID == 0 {
		return nil, errors.New("invalid approval log id")
	}
	return u.repo.GetApprovalLogByID(ctx, approvalLogID)
}
