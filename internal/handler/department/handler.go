package department

import (
	departmentDTO "backend/internal/dto/department_dto"
	"backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type DepartmentHandler struct {
	service usecase.DepartmentService
}

func NewDepartmentHandler(service usecase.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// CreateDepartment สร้าง department ใหม่
func (h *DepartmentHandler) CreateDepartment(c *fiber.Ctx) error {
	req := new(departmentDTO.CreateDepartment)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	department, err := h.service.CreateDepartment(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &departmentDTO.DepartmentResponse{
		DepartmentID:   department.DepartmentID,
		FacultyID:      department.FacultyID,
		DepartmentName: department.DepartmentName,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Department created successfully",
		"data":    response,
	})
}

// GetAllDepartments ดึงข้อมูล department ทั้งหมด
func (h *DepartmentHandler) GetAllDepartments(c *fiber.Ctx) error {
	departments, err := h.service.GetAllDepartments(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var responses []departmentDTO.DepartmentResponse
	for _, d := range departments {
		responses = append(responses, departmentDTO.DepartmentResponse{
			DepartmentID:   d.DepartmentID,
			FacultyID:      d.FacultyID,
			DepartmentName: d.DepartmentName,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Departments retrieved successfully",
		"data":    responses,
	})
}

// GetDepartmentByID ดึงข้อมูล department ตามรหัส
func (h *DepartmentHandler) GetDepartmentByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid department ID",
		})
	}

	department, err := h.service.GetDepartmentByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Department not found",
		})
	}

	response := &departmentDTO.DepartmentResponse{
		DepartmentID:   department.DepartmentID,
		FacultyID:      department.FacultyID,
		DepartmentName: department.DepartmentName,
	}

	return c.JSON(fiber.Map{
		"message": "Department retrieved successfully",
		"data":    response,
	})
}

// GetDepartmentsByFacultyID ดึงข้อมูล department ตามรหัส faculty
func (h *DepartmentHandler) GetDepartmentsByFacultyID(c *fiber.Ctx) error {
	facultyID, err := strconv.ParseUint(c.Params("facultyId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid faculty ID",
		})
	}

	departments, err := h.service.GetDepartmentsByFacultyID(c.Context(), uint(facultyID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var responses []departmentDTO.DepartmentResponse
	for _, d := range departments {
		responses = append(responses, departmentDTO.DepartmentResponse{
			DepartmentID:   d.DepartmentID,
			FacultyID:      d.FacultyID,
			DepartmentName: d.DepartmentName,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Departments retrieved successfully",
		"data":    responses,
	})
}

// UpdateDepartment แก้ไข department
func (h *DepartmentHandler) UpdateDepartment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid department ID",
		})
	}

	req := new(departmentDTO.UpdateDepartment)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	department, err := h.service.UpdateDepartment(c.Context(), uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &departmentDTO.DepartmentResponse{
		DepartmentID:   department.DepartmentID,
		FacultyID:      department.FacultyID,
		DepartmentName: department.DepartmentName,
	}

	return c.JSON(fiber.Map{
		"message": "Department updated successfully",
		"data":    response,
	})
}

// DeleteDepartment ลบ department
func (h *DepartmentHandler) DeleteDepartment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid department ID",
		})
	}

	if err := h.service.DeleteDepartment(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Department deleted successfully",
	})
}
