package user

import (
	authDto "backend/internal/dto/auth_dto"
	userdto "backend/internal/dto/user_dto"
	"backend/internal/models"
	"backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService usecase.UserService
	AuthService usecase.AuthService
}

func NewUserHandler(us usecase.UserService) *UserHandler {
	return &UserHandler{UserService: us}
}

func NewUserHandlerWithAuth(us usecase.UserService, as usecase.AuthService) *UserHandler {
	return &UserHandler{UserService: us, AuthService: as}
}

// GET /users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	user, err := h.UserService.GetUserByID(c.Context(), uint(id64))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	response := authDto.MeResponse{
		UserID:       user.UserID,
		Prefix:       user.Prefix,
		Firstname:    user.Firstname,
		Lastname:     user.Lastname,
		Email:        user.Email,
		ImagePath:    user.ImagePath,
		Provider:     user.Provider,
		RoleID:       user.RoleID,
		CampusID:     user.CampusID,
		IsFirstLogin: user.IsFirstLogin,
		CreatedAt:    user.CreatedAt,
		LatestUpdate: user.LatestUpdate,
	}

	if h.AuthService != nil {
		if user.RoleID == 2 {
			hod, err := h.AuthService.GetHeadOfDepartmentByUserID(c.Context(), user.UserID)
			if err == nil && hod != nil {
				response.HeadOfDepartmentData = &authDto.HeadOfDepartmentMeData{
					HodID:        hod.HodID,
					UserID:       hod.UserID,
					FacultyID:    hod.FacultyID,
					DepartmentID: hod.DepartmentID,
				}
			}
		}
		if user.RoleID == 3 {
			ad, err := h.AuthService.GetAssociateDeanByUserID(c.Context(), user.UserID)
			if err == nil && ad != nil {
				response.AssociateDeanData = &authDto.AssociateDeanMeData{
					AdID:      ad.AdID,
					UserID:    ad.UserID,
					AdCode:    ad.AdCode,
					FacultyID: ad.FacultyID,
				}
			}
		}
		if user.RoleID == 4 {
			dean, err := h.AuthService.GetDeanByUserID(c.Context(), user.UserID)
			if err == nil && dean != nil {
				response.DeanData = &authDto.DeanMeData{
					DID:       dean.DID,
					UserID:    dean.UserID,
					FacultyID: dean.FacultyID,
				}
			}
		}
		if user.RoleID == 5 {
			sd, err := h.AuthService.GetStudentDevelopmentByUserID(c.Context(), user.UserID)
			if err == nil && sd != nil {
				response.StudentDevelopmentData = &authDto.StudentDevelopmentMeData{
					SdID:   sd.SdID,
					UserID: sd.UserID,
				}
			}
		}
		if user.RoleID == 6 {
			com, err := h.AuthService.GetCommitteeByUserID(c.Context(), user.UserID)
			if err == nil && com != nil {
				response.CommitteeData = &authDto.CommitteeMeData{
					ComID:      com.ComID,
					UserID:     com.UserID,
					IsChairman: com.IsChairman,
				}
			}
		}
		if user.RoleID == 7 {
			ch, err := h.AuthService.GetChancellorByUserID(c.Context(), user.UserID)
			if err == nil && ch != nil {
				response.ChancellorData = &authDto.ChancellorMeData{
					ChancellorID: ch.ChancellorID,
					UserID:       ch.UserID,
				}
			}
		}
	}

	return c.JSON(response)
}

// GET /users (ดึง user ตามวิทยาเขตของคนที่ login อยู่)
func (h *UserHandler) GetAllUsersByCampus(c *fiber.Ctx) error {
	currentUser, ok := c.Locals("current_user").(*models.User)

	if !ok || currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	const limit = 6

	users, err := h.UserService.GetAllUsersByCampus(c.Context(), currentUser.CampusID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var responses []userdto.UserResponse
	for _, user := range users {
		responses = append(responses, userdto.UserResponse{
			UserID:       user.UserID,
			Prefix:       user.Prefix,
			Firstname:    user.Firstname,
			Lastname:     user.Lastname,
			Email:        user.Email,
			ImagePath:    user.ImagePath,
			Provider:     user.Provider,
			RoleID:       user.RoleID,
			CampusID:     user.CampusID,
			IsFirstLogin: user.IsFirstLogin,
			CreatedAt:    user.CreatedAt,
			LatestUpdate: user.LatestUpdate,
		})
	}

	return c.JSON(fiber.Map{
		"page":  page,
		"limit": limit,
		"size":  len(responses),
		"data":  responses,
	})
}

// PUT /users/:id
func (h *UserHandler) UpdateUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	var req userdto.EditUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	updated, err := h.UserService.UpdateUserByID(c.Context(), uint(id64), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(userdto.UserResponse{
		UserID:       updated.UserID,
		Prefix:       updated.Prefix,
		Firstname:    updated.Firstname,
		Lastname:     updated.Lastname,
		Email:        updated.Email,
		ImagePath:    updated.ImagePath,
		Provider:     updated.Provider,
		RoleID:       updated.RoleID,
		CampusID:     updated.CampusID,
		IsFirstLogin: updated.IsFirstLogin,
		CreatedAt:    updated.CreatedAt,
		LatestUpdate: updated.LatestUpdate,
	})
}

// PUT /users/promote-chairman/:id
func (h *UserHandler) ChangeCommitteeRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	updated, err := h.UserService.ChangeCommitteeRoleByID(c.Context(), uint(id64), true)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(userdto.UserResponse{
		UserID:       updated.UserID,
		Prefix:       updated.Prefix,
		Firstname:    updated.Firstname,
		Lastname:     updated.Lastname,
		Email:        updated.Email,
		ImagePath:    updated.ImagePath,
		Provider:     updated.Provider,
		RoleID:       updated.RoleID,
		CampusID:     updated.CampusID,
		IsFirstLogin: updated.IsFirstLogin,
		CreatedAt:    updated.CreatedAt,
		LatestUpdate: updated.LatestUpdate,
	})
}
