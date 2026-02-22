package usecase

import (
	userdto "backend/internal/dto/user_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	UpdateUserByID(ctx context.Context, userID uint, req *userdto.EditUserRequest) (*models.User, error)
	GetAllUsersByCampus(ctx context.Context, campusID int) ([]models.User, error)
	CreateUser(ctx context.Context, req *userdto.CreateUserRequest) (*models.User, error)
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

// GetAllUsersByCampus ดึงข้อมูล user ทั้งหมดตามวิทยาเขต
func (s *userService) GetAllUsersByCampus(ctx context.Context, campusID int) ([]models.User, error) {
	if campusID <= 0 {
		return nil, errors.New("invalid campus id")
	}
	users, err := s.repo.GetUserListByCampus(ctx, campusID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser สร้าง user ใหม่ (สำหรับ Role 5)
func (s *userService) CreateUser(ctx context.Context, req *userdto.CreateUserRequest) (*models.User, error) {
	// 1. ตรวจสอบ field ที่จำเป็น (ทำใน handler หรือ validate struct แล้ว แต่ check อีกทีได้)
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		return nil, errors.New("email is required")
	}

	// 2. ตรวจสอบว่ามี Email นี้ซ้ำหรือไม่
	existing, err := s.repo.GetUserByEmail(email)
	if err == nil && existing != nil {
		return nil, errors.New("email already registered")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 3. ตรวจสอบ Password Match
	if req.Password != req.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	// 4. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 5. เตรียม Model
	user := &models.User{
		Email:          email,
		HashedPassword: string(hashed),
		Prefix:         strings.TrimSpace(req.Prefix),
		Firstname:      strings.TrimSpace(req.Firstname),
		Lastname:       strings.TrimSpace(req.Lastname),
		RoleID:         req.RoleID,
		CampusID:       req.CampusID,	
		Provider:       "manual",
		IsFirstLogin:   true, // ให้ user ใหม่ต้อง setup profile อีกที หรือแล้วแต่ requirement
		CreatedAt:      time.Now(),
		LatestUpdate:   time.Now(),
	}
	// ถ้าต้องการ ImagePath default
	// user.ImagePath = "..."

	// 6. Save to DB
	if err := s.repo.UpsertUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
