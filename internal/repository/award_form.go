package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AwardRepository struct {
	db *gorm.DB
}

type AwardSearchFilter struct {
	CampusID             int
	Keyword              string
	Date                 string
	StudentYear          int
	AwardType            string
	ExcludeVotedByUserID *uint
	FacultyID            *int
	DepartmentID         *int
	FormStatusID         *int
	SortBy               string
	SortOrder            string
	Page                 int
	Limit                int
}

type ApprovalLogSearchFilter struct {
	UserID    uint
	CampusID  int
	Keyword   string
	Date      string
	AwardType string
	Operation string
	SortBy    string
	SortOrder string
	Page      int
	Limit     int
}

type ApprovalLogHistoryRow struct {
	ApprovalLogID    uint      `gorm:"column:approval_log_id"`
	FormID           uint      `gorm:"column:form_id"`
	ReviewerUserID   uint      `gorm:"column:reviewer_user_id"`
	Operation        string    `gorm:"column:operation"`
	OperationDate    time.Time `gorm:"column:operation_date"`
	StudentFirstname string    `gorm:"column:student_firstname"`
	StudentLastname  string    `gorm:"column:student_lastname"`
	StudentNumber    string    `gorm:"column:student_number"`
	AcademicYear     int       `gorm:"column:academic_year"`
	AwardType        string    `gorm:"column:award_type"`
	CampusID         int       `gorm:"column:campus_id"`
}

func NewAwardRepository(db *gorm.DB) *AwardRepository {
	return &AwardRepository{db: db}
}

// ปรับปรุง: เพิ่มพารามิเตอร์ files เพื่อรองรับการบันทึกไฟล์แนบ (ถ้ามี)
func (r *AwardRepository) CreateWithTransaction(ctx context.Context, form *models.AwardForm, files []models.AwardFileDirectory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. บันทึกตารางหลัก (Award_Form)
		if err := tx.Create(form).Error; err != nil {
			return err
		}

		// 2. บันทึกตารางไฟล์แนบ (ถ้า len > 0 คือมีการแนบไฟล์มา)
		if len(files) > 0 {
			for i := range files {
				// ผูก ID ของไฟล์เข้ากับ FormID ที่เพิ่งสร้างใหม่
				files[i].FormID = form.FormID
				if err := tx.Create(&files[i]).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetByKeyword ค้นหาและกรองพร้อม pagination ตาม role scope
func (r *AwardRepository) GetByKeyword(ctx context.Context, filter AwardSearchFilter) ([]models.AwardForm, int64, error) {
	var list []models.AwardForm
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AwardForm{}).Where("campus_id = ?", filter.CampusID)

	// ค้นหาด้วย keyword (firstname, lastname, studentNumber, semester, year, award_type)
	if filter.Keyword != "" {
		query = query.Where(
			"student_firstname LIKE ? OR student_lastname LIKE ? OR student_number LIKE ? OR CONCAT(student_firstname, ' ', student_lastname) LIKE ?",
			"%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%",
		)
	}

	// กรองตามวันที่ (ถ้ามี)
	if filter.Date != "" {
		query = query.Where("DATE(created_at) = ?", filter.Date)
	}

	// กรองตามชั้นปี (ถ้ามี)
	if filter.StudentYear > 0 {
		query = query.Where("student_year = ?", filter.StudentYear)
	}

	// กรองตามประเภทรางวัล (ถ้ามี)
	if filter.AwardType != "" {
		query = query.Where("award_type = ?", filter.AwardType)
	}

	if filter.FacultyID != nil {
		query = query.Where("faculty_id = ?", *filter.FacultyID)
	}

	if filter.DepartmentID != nil {
		query = query.Where("department_id = ?", *filter.DepartmentID)
	}

	if filter.FormStatusID != nil {
		query = query.Where("form_status_id = ?", *filter.FormStatusID)
	}

	if filter.ExcludeVotedByUserID != nil {
		query = query.Where(
			`NOT EXISTS (
				SELECT 1
				FROM "Committee_Vote_Log" cvl
				WHERE cvl.form_id = "Award_Form".form_id
				  AND cvl.user_id = ?
			)`,
			*filter.ExcludeVotedByUserID,
		)
	}

	// นับจำนวนทั้งหมด
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// คำนวณ offset
	offset := (filter.Page - 1) * filter.Limit

	// ดึงข้อมูลพร้อม pagination และ preload
	orderClause := buildOrderClause(filter.SortBy, filter.SortOrder)

	err := query.
		Preload("AwardFiles").
		Order(orderClause).
		Limit(filter.Limit).
		Offset(offset).
		Find(&list).Error

	return list, total, err
}

func buildOrderClause(sortBy string, sortOrder string) string {
	order := strings.ToUpper(sortOrder)
	if order != "ASC" {
		order = "DESC"
	}

	switch sortBy {
	case "name":
		return fmt.Sprintf("student_firstname %s, student_lastname %s", order, order)
	case "studentNumber":
		return fmt.Sprintf("student_number %s", order)
	case "academicYear":
		return fmt.Sprintf("academic_year %s", order)
	case "awardType":
		return fmt.Sprintf("award_type %s", order)
	case "date":
		fallthrough
	default:
		return fmt.Sprintf("created_at %s", order)
	}
}

func (r *AwardRepository) GetApprovalHistoryByUserAndCampus(ctx context.Context, filter ApprovalLogSearchFilter) ([]ApprovalLogHistoryRow, int64, error) {
	var rows []ApprovalLogHistoryRow
	var total int64

	query := r.db.WithContext(ctx).
		Table(`"Award_Approval_Log" AS aal`).
		Joins(`JOIN "Award_Form" af ON af.form_id = aal.form_id`).
		Where("aal.user_id = ?", filter.UserID).
		Where("af.campus_id = ?", filter.CampusID)
		
	if filter.Keyword != "" {
		query = query.Where(
			"af.student_firstname LIKE ? OR af.student_lastname LIKE ? OR af.student_number LIKE ? OR CONCAT(af.student_firstname, ' ', af.student_lastname) LIKE ?",
			"%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%",
		)
	}

	if filter.Date != "" {
		query = query.Where("DATE(aal.approved_at) = ?", filter.Date)
	}

	if filter.AwardType != "" {
		query = query.Where("af.award_type = ?", filter.AwardType)
	}

	if filter.Operation != "" {
		query = query.Where("aal.approval_status = ?", filter.Operation)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit

	err := query.Select(`
		aal.approval_log_id AS approval_log_id,
		aal.form_id AS form_id,
		aal.user_id AS reviewer_user_id,
		aal.approval_status AS operation,
		aal.approved_at AS operation_date,
		af.student_firstname AS student_firstname,
		af.student_lastname AS student_lastname,
		af.student_number AS student_number,
		af.academic_year AS academic_year,
		af.award_type AS award_type,
		af.campus_id AS campus_id
	`).
		Order(buildApprovalLogOrderClause(filter.SortBy, filter.SortOrder)).
		Limit(filter.Limit).
		Offset(offset).
		Scan(&rows).Error

	return rows, total, err
}

func buildApprovalLogOrderClause(sortBy string, sortOrder string) string {
	order := strings.ToUpper(sortOrder)
	if order != "ASC" {
		order = "DESC"
	}

	switch sortBy {
	case "name":
		return fmt.Sprintf("af.student_firstname %s, af.student_lastname %s", order, order)
	case "studentNumber":
		return fmt.Sprintf("af.student_number %s", order)
	case "academicYear":
		return fmt.Sprintf("af.academic_year %s", order)
	case "awardType":
		return fmt.Sprintf("af.award_type %s", order)
	case "date":
		fallthrough
	default:
		return fmt.Sprintf("aal.approved_at %s", order)
	}
}

func (r *AwardRepository) GetHeadOfDepartmentScopeByUserID(ctx context.Context, userID uint) (int, int, error) {
	var scope struct {
		FacultyID    int `gorm:"column:faculty_id"`
		DepartmentID int `gorm:"column:department_id"`
	}

	err := r.db.WithContext(ctx).
		Table("Head_Of_Department").
		Select("faculty_id, department_id").
		Where("user_id = ?", userID).
		Take(&scope).Error
	if err != nil {
		return 0, 0, err
	}

	return scope.FacultyID, scope.DepartmentID, nil
}

func (r *AwardRepository) GetFacultyScopeByRoleAndUserID(ctx context.Context, roleID int, userID uint) (int, error) {
	tableName := ""
	switch roleID {
	case 3:
		tableName = "Associate_Dean"
	case 4:
		tableName = "Dean"
	default:
		return 0, nil
	}

	var scope struct {
		FacultyID int `gorm:"column:faculty_id"`
	}

	err := r.db.WithContext(ctx).
		Table(tableName).
		Select("faculty_id").
		Where("user_id = ?", userID).
		Take(&scope).Error
	if err != nil {
		return 0, err
	}

	return scope.FacultyID, nil
}

func (r *AwardRepository) GetByType(ctx context.Context, awardType string, campusID int) ([]models.AwardForm, error) {
	var list []models.AwardForm
	err := r.db.WithContext(ctx).
		Where("award_type = ? AND campus_id = ?", awardType, campusID).
		Preload("AwardFiles").
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *AwardRepository) GetByUserID(ctx context.Context, userID uint) ([]models.AwardForm, error) {
	var list []models.AwardForm
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("AwardFiles").
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *AwardRepository) GetByUserIDAndYear(ctx context.Context, userID uint, year int) ([]models.AwardForm, error) {
	var list []models.AwardForm
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND academic_year = ?", userID, year).
		Preload("AwardFiles").
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *AwardRepository) GetByUserIDAndSemester(ctx context.Context, userID uint, year int, semester int) ([]models.AwardForm, error) {
	var list []models.AwardForm
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND academic_year = ? AND semester = ?", userID, year, semester).
		Preload("AwardFiles").
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *AwardRepository) GetByUserIDWithYearSort(ctx context.Context, userID uint, years []int) ([]models.AwardForm, error) {
	var list []models.AwardForm
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if len(years) > 0 {
		query = query.Where("academic_year IN ?", years)
	}
	err := query.
		Preload("AwardFiles").
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *AwardRepository) GetByUserIDWithYearPaged(ctx context.Context, userID uint, years []int, page int, limit int) ([]models.AwardForm, int64, error) {
	var list []models.AwardForm
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AwardForm{}).Where("user_id = ?", userID)
	if len(years) > 0 {
		query = query.Where("academic_year IN ?", years)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.
		Preload("AwardFiles").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&list).Error

	return list, total, err
}

func (r *AwardRepository) GetByStudentID(ctx context.Context, studentID int) ([]models.AwardForm, error) {
	var list []models.AwardForm
	err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Preload("AwardFiles").
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *AwardRepository) GetByFormID(ctx context.Context, formID int) (*models.AwardForm, error) {
	var form models.AwardForm
	err := r.db.WithContext(ctx).
		Where("form_id = ?", formID).
		Preload("AwardFiles").
		First(&form).Error
	if err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *AwardRepository) CheckDuplicate(userID uint, year int, semester int) (bool, error) {
	var count int64
	// เช็คในตาราง AwardForm ว่ามีข้อมูลที่ user_id, academic_year, semester ตรงกันไหม
	err := r.db.Model(&models.AwardForm{}).
		Where("user_id = ? AND academic_year = ? AND semester = ?", userID, year, semester).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AwardRepository) UpdateAwardType(ctx context.Context, formID uint, awardType string) error {
	return r.db.WithContext(ctx).
		Model(&models.AwardForm{}).
		Where("form_id = ?", formID).
		Updates(map[string]interface{}{
			"award_type":    awardType,
			"latest_update": time.Now(),
		}).Error
}

func (r *AwardRepository) UpdateFormStatus(ctx context.Context, formID uint, formStatus int, rejectReason string) error {
	return r.db.WithContext(ctx).
		Model(&models.AwardForm{}).
		Where("form_id = ?", formID).
		Updates(map[string]interface{}{
			"form_status_id": formStatus,
			"reject_reason":  rejectReason,
			"latest_update":  time.Now(),
		}).Error
}

func (r *AwardRepository) CreateAwardApprovalLog(ctx context.Context, log *models.AwardApprovalLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AwardRepository) CreateAwardSignedLog(ctx context.Context, log *models.AwardSignedLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AwardRepository) SaveAwardTypeLog(ctx context.Context, log *models.AwardTypeLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AwardRepository) IsCommitteeNonChairman(ctx context.Context, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("\"Committee\" c").
		Joins("JOIN \"User\" u ON u.user_id = c.user_id").
		Where("c.user_id = ?", userID).
		Where("u.role_id = ?", 6).
		Where("c.is_chairman = ?", false).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AwardRepository) UpsertCommitteeVoteLog(ctx context.Context, formID uint, userID uint, operation string) error {
	var existing models.CommitteeVoteLog
	err := r.db.WithContext(ctx).
		Where("form_id = ? AND user_id = ?", formID, userID).
		First(&existing).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log := &models.CommitteeVoteLog{
				FormID:    formID,
				UserID:    userID,
				Operation: operation,
				VotedAt:   time.Now(),
			}
			return r.db.WithContext(ctx).Create(log).Error
		}
		return err
	}

	existing.Operation = operation
	existing.VotedAt = time.Now()
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *AwardRepository) CountNonChairmanCommittees(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Table("\"Committee\" c").
		Joins("JOIN \"User\" u ON u.user_id = c.user_id").
		Where("u.role_id = ?", 6).
		Where("c.is_chairman = ?", false).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *AwardRepository) CountCommitteeVotesByOperation(ctx context.Context, formID uint, operation string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.CommitteeVoteLog{}).
		Where("form_id = ? AND operation = ?", formID, operation).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *AwardRepository) GetApprovalLogsByUserID(ctx context.Context, userID uint) ([]models.AwardApprovalLog, error) {
	var logs []models.AwardApprovalLog
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("approved_at desc").
		Find(&logs).Error
	return logs, err
}

// GetApprovalLogByID fetches a single approval log by its ID
func (r *AwardRepository) GetApprovalLogByID(ctx context.Context, approvalLogID uint) (*models.AwardApprovalLog, error) {
	var log models.AwardApprovalLog
	err := r.db.WithContext(ctx).
		Where("approval_log_id = ?", approvalLogID).
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetAllAwardTypes ดึง award_type ทั้งหมดที่มีในตาราง (ไม่ซ้ำกัน)
func (r *AwardRepository) GetAllAwardTypes(ctx context.Context) ([]string, error) {
	var awardTypes []string
	err := r.db.WithContext(ctx).
		Model(&models.AwardForm{}).
		Distinct("award_type").
		Order("award_type ASC").
		Pluck("award_type", &awardTypes).Error

	if err != nil {
		return nil, err
	}

	return awardTypes, nil
}
