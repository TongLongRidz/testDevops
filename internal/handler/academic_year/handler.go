package academic_year

import (
    "backend/internal/models"
    "backend/internal/usecase"
    "errors"

    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

type AcademicYearHandler struct {
    svc *usecase.AcademicYearService
}

func NewAcademicYearHandler(svc *usecase.AcademicYearService) *AcademicYearHandler {
    return &AcademicYearHandler{svc: svc}
}

func (h *AcademicYearHandler) Update(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

    var ay models.AcademicYear
    if err := c.BodyParser(&ay); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
    }
    ay.AcademicYearID = uint(id)

    out, err := h.svc.Update(c.UserContext(), &ay)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(out)
}

func (h *AcademicYearHandler) Delete(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

    if err := h.svc.Delete(c.UserContext(), uint(id)); err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return c.SendStatus(fiber.StatusNoContent)
}

func (h *AcademicYearHandler) Create(c *fiber.Ctx) error {
	var ay models.AcademicYear
	if err := c.BodyParser(&ay); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	// คุณอาจจะเพิ่ม validation เล็กๆ น้อยๆ ตรงนี้เหมือนตอนเช็ค Email/Password
	if ay.Year == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "year is required"})
	}

	out, err := h.svc.Create(c.UserContext(), &ay)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "created successfully",
		"data":    out,
	})
}

func (h *AcademicYearHandler) GetList(c *fiber.Ctx) error {
	out, err := h.svc.GetList(c.UserContext())
	if err != nil {
		// ถ้าเกิด Error ในระดับ Database/Service
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// ถ้าไม่มีข้อมูลใน list เลย (เป็น empty slice)
	if len(out) == 0 {
		return c.JSON(fiber.Map{
			"data": []models.AcademicYear{}, // ส่ง array เปล่ากลับไปเพื่อให้ Front-end loop ได้ไม่พัง
		})
	}

	return c.JSON(fiber.Map{
		"data": out,
	})
}

func (h *AcademicYearHandler) GetLatest(c *fiber.Ctx) error {
	out, err := h.svc.GetLatest(c.UserContext())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no academic year records found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": out,
	})
}

func (h *AcademicYearHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	out, err := h.svc.GetByID(c.UserContext(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "academic year not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": out,
	})
}