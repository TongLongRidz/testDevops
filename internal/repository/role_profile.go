package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type RoleProfileRepository interface {
	CreateByRole(ctx context.Context, roleID int, userID uint, facultyID uint, departmentID uint) error
	GetHeadOfDepartmentByUserID(ctx context.Context, userID uint) (*models.HeadOfDepartment, error)
	GetAssociateDeanByUserID(ctx context.Context, userID uint) (*models.AssociateDean, error)
	GetDeanByUserID(ctx context.Context, userID uint) (*models.Dean, error)
	GetStudentDevelopmentByUserID(ctx context.Context, userID uint) (*models.StudentDevelopment, error)
	GetCommitteeByUserID(ctx context.Context, userID uint) (*models.Committee, error)
	GetChancellorByUserID(ctx context.Context, userID uint) (*models.Chancellor, error)
}

type roleProfileRepository struct {
	db *gorm.DB
}

func NewRoleProfileRepository(db *gorm.DB) RoleProfileRepository {
	return &roleProfileRepository{db: db}
}

func (r *roleProfileRepository) CreateByRole(ctx context.Context, roleID int, userID uint, facultyID uint, departmentID uint) error {
	switch roleID {
	case 2:
		profile := &models.HeadOfDepartment{
			UserID:       userID,
			FacultyID:    facultyID,
			DepartmentID: departmentID,
		}
		return r.db.WithContext(ctx).Create(profile).Error
	case 3:
		profile := &models.AssociateDean{
			UserID:    userID,
			FacultyID: facultyID,
		}
		return r.db.WithContext(ctx).Create(profile).Error
	case 4:
		profile := &models.Dean{
			UserID:    userID,
			FacultyID: facultyID,
		}
		return r.db.WithContext(ctx).Create(profile).Error
	case 5:
		profile := &models.StudentDevelopment{UserID: userID}
		return r.db.WithContext(ctx).Create(profile).Error
	case 6:
		profile := &models.Committee{UserID: userID, IsChairman: false}
		return r.db.WithContext(ctx).Create(profile).Error
	case 7:
		profile := &models.Chancellor{UserID: userID}
		return r.db.WithContext(ctx).Create(profile).Error
	case 8:
		profile := &models.Organization{
			UserID:                  userID,
			OrganizationName:        "",
			OrganizationType:        "",
			OrganizationLocation:    "",
			OrganizationPhoneNumber: "",
		}
		return r.db.WithContext(ctx).Create(profile).Error
	default:
		return nil
	}
}

func (r *roleProfileRepository) GetHeadOfDepartmentByUserID(ctx context.Context, userID uint) (*models.HeadOfDepartment, error) {
	var profile models.HeadOfDepartment
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *roleProfileRepository) GetAssociateDeanByUserID(ctx context.Context, userID uint) (*models.AssociateDean, error) {
	var profile models.AssociateDean
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *roleProfileRepository) GetDeanByUserID(ctx context.Context, userID uint) (*models.Dean, error) {
	var profile models.Dean
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *roleProfileRepository) GetStudentDevelopmentByUserID(ctx context.Context, userID uint) (*models.StudentDevelopment, error) {
	var profile models.StudentDevelopment
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *roleProfileRepository) GetCommitteeByUserID(ctx context.Context, userID uint) (*models.Committee, error) {
	var profile models.Committee
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *roleProfileRepository) GetChancellorByUserID(ctx context.Context, userID uint) (*models.Chancellor, error) {
	var profile models.Chancellor
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}
