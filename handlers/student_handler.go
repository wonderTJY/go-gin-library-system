package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"trae-go/models"

	"trae-go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StudentHandler struct {
	DB *gorm.DB
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	return &StudentHandler{DB: db}
}

func (h *StudentHandler) ListStudents(c *gin.Context) {
	var students []models.Student
	if err := h.DB.Find(&students).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_LIST_STUDENTS", "failed to list students"))
		return
	}
	c.JSON(http.StatusOK, students)
}
func (h *StudentHandler) GetStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_ID", "invalid id"))
		return
	}
	var student models.Student
	if err := h.DB.First(&student, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(middleware.NewAppError(http.StatusNotFound, "STUDENT_NOT_FOUND", "student not found"))
			return
		}
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, student)
}
func (h *StudentHandler) CreatStudent(c *gin.Context) {
	var input models.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid json"))
		return
	}
	student := models.Student{
		Name:  input.Name,
		Email: input.Email,
	}
	if err := h.DB.Create(&student).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "FAILED_CREATE_STUDENT", "failed to create student"))
		return
	}
	c.JSON(http.StatusCreated, student)
}
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_ID", "invalid id"))
		return
	}
	var input models.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid json"))
		return
	}
	var student models.Student
	if err := h.DB.First(&student, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(middleware.NewAppError(http.StatusNotFound, "STUDENT_NOT_FOUND", "student not found"))
			return
		}
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	student.Name = input.Name
	student.Email = input.Email
	if err := h.DB.Save(&student).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, student)
}
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_ID", "invalid id"))
		return
	}
	result := h.DB.Delete(&models.Student{}, uint(id))
	if result.Error != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	if result.RowsAffected == 0 {
		c.Error(middleware.NewAppError(http.StatusNotFound, "STUDENT_NOT_FOUND", "student not found"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "student deleted"})
}
func (h *StudentHandler) PanicTest(c *gin.Context) {
	panic("kkk")
}
