package usecase

import (
	userdto "backend/internal/dto/user_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"strings"
	"time"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	UpdateUserByID(ctx context.Context, userID uint, req *userdto.EditUserRequest) (*models.User, error)
	ChangeCommitteeRoleByID(ctx context.Context, userID uint, isChairman bool) (*models.User, error)
	GetAllUsersByCampus(ctx context.Context, campusID int, page int, limit int) ([]models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// GetUserByID ดึงข้อมูล user ตาม ID
func (s *userService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUserByID แก้ไขข้อมูล user (ยกเว้น password)
func (s *userService) UpdateUserByID(ctx context.Context, userID uint, req *userdto.EditUserRequest) (*models.User, error) {
	updates := map[string]interface{}{
		"latest_update": time.Now(),
	}

	if req.Firstname != nil {
		updates["firstname"] = strings.TrimSpace(*req.Firstname)
	}
	if req.Prefix != nil {
		updates["prefix"] = strings.TrimSpace(*req.Prefix)
	}
	if req.Lastname != nil {
		updates["lastname"] = strings.TrimSpace(*req.Lastname)
	}
	if req.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*req.Email))
		if email == "" {
			return nil, errors.New("email cannot be empty")
		}
		// ตรวจสอบ email ซ้ำ
		if existing, err := s.repo.GetUserByEmail(email); err == nil && existing != nil && existing.UserID != userID {
			return nil, errors.New("email already in use")
		}
		updates["email"] = email
	}
	if req.ImagePath != nil {
		updates["image_path"] = strings.TrimSpace(*req.ImagePath)
	}
	if req.Provider != nil {
		updates["provider"] = strings.TrimSpace(*req.Provider)
	}
	if req.RoleID != nil {
		updates["role_id"] = *req.RoleID
	}
	if req.CampusID != nil {
		updates["campus_id"] = *req.CampusID
	}
	if req.IsFirstLogin != nil {
		updates["is_first_login"] = *req.IsFirstLogin
	}

	updated, err := s.repo.UpdateUserFields(ctx, userID, updates)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *userService) ChangeCommitteeRoleByID(ctx context.Context, userID uint, isChairman bool) (*models.User, error) {
	if err := s.repo.SetCommitteeChairman(ctx, userID, isChairman); err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(userID)
}

// GetAllUsersByCampus ดึงข้อมูล user ตามวิทยาเขตแบบแบ่งหน้า
func (s *userService) GetAllUsersByCampus(ctx context.Context, campusID int, page int, limit int) ([]models.User, error) {
	if campusID <= 0 {
		return nil, errors.New("invalid campus id")
	}
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 6
	}
	if limit > 6 {
		limit = 6
	}

	users, err := s.repo.GetUserListByCampus(ctx, campusID, page, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}
