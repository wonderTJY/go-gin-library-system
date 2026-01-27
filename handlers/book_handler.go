package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"trae-go/middleware"
	"trae-go/models"
)

type BookHandler struct {
	DB *gorm.DB
}

func NewBookHandler(db *gorm.DB) *BookHandler {
	return &BookHandler{DB: db}
}

// ListBooks 获取书籍列表
// @Summary      获取书籍列表
// @Description  获取所有书籍信息
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Book
// @Failure      500  {object}  middleware.AppError
// @Router       /books [get]
func (h *BookHandler) ListBooks(c *gin.Context) {
	var books []models.Book
	if err := h.DB.Find(&books).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_LIST_BOOKS", "failed to list books"))
		return
	}
	c.JSON(http.StatusOK, books)
}

// GetBook 获取单本书籍
// @Summary      获取单本书籍
// @Description  根据 ID 获取书籍详情
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "书籍 ID"
// @Success      200  {object}  models.Book
// @Failure      400  {object}  middleware.AppError
// @Failure      404  {object}  middleware.AppError
// @Router       /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_ID", "invalid id"))
		return
	}
	var book models.Book
	if err := h.DB.First(&book, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(middleware.NewAppError(http.StatusNotFound, "BOOK_NOT_FOUND", "book not found"))
			return
		}
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_GET_BOOK", "failed to get book"))
		return
	}
	c.JSON(http.StatusOK, book)
}

// CreateBook 创建书籍
// @Summary      创建书籍
// @Description  添加一本新书
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.Book true "书籍信息"
// @Success      201  {object}  models.Book
// @Failure      400  {object}  middleware.AppError
// @Router       /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid json"))
		return
	}
	book := models.Book{
		Title:  input.Title,
		Author: input.Author,
	}
	if err := h.DB.Create(&book).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_CREATE_BOOK", "failed to create book"))
		return
	}
	c.JSON(http.StatusCreated, book)
}

// UpdateBook 更新书籍
// @Summary      更新书籍
// @Description  根据 ID 更新书籍信息
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int          true  "书籍 ID"
// @Param        request body      models.Book  true  "更新信息"
// @Success      200     {object}  models.Book
// @Failure      400     {object}  middleware.AppError
// @Failure      404     {object}  middleware.AppError
// @Router       /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_ID", "invalid id"))
		return
	}
	var book models.Book
	if err := h.DB.First(&book, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(middleware.NewAppError(http.StatusNotFound, "BOOK_NOT_FOUND", "book not found"))
			return
		}
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_GET_BOOK", "failed to get book"))
		return
	}
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid json"))
		return
	}
	book.Title = input.Title
	book.Author = input.Author
	book.Stock = input.Stock
	if err := h.DB.Save(&book).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_UPDATE_BOOK", "failed to update book"))
		return
	}
	c.JSON(http.StatusOK, book)
}

// DeleteBook 删除书籍
// @Summary      删除书籍
// @Description  根据 ID 删除书籍
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "书籍 ID"
// @Success      204  "No Content"
// @Failure      400  {object}  middleware.AppError
// @Failure      500  {object}  middleware.AppError
// @Router       /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_ID", "invalid id"))
		return
	}
	if err := h.DB.Delete(&models.Book{}, uint(id)).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_DELETE_BOOK", "failed to delete book"))
		return
	}
	c.Status(http.StatusNoContent)
}

// @Tags         borrow
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        student_id  path      int  true  "学生 ID"
// @Param        book_id     path      int  true  "书籍 ID"
func (h *BookHandler) BookABook(c *gin.Context) {
	stuidstr := c.Param("student_id")
	bookidstr := c.Param("book_id")

	stuid, err := strconv.ParseUint(stuidstr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_STUDENT_ID", "invalid student_id"))
		return
	}
	bookid, err := strconv.ParseUint(bookidstr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_BOOK_ID", "invalid book_id"))
		return
	}

	var student models.Student
	var book models.Book

	if err := h.DB.First(&student, uint(stuid)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(middleware.NewAppError(http.StatusNotFound, "STUDENT_NOT_FOUND", "student not found"))
			return
		}
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	if err := h.DB.First(&book, uint(bookid)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(middleware.NewAppError(http.StatusNotFound, "BOOK_NOT_FOUND", "book not found"))
			return
		}
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	if book.Stock <= 0 {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "BOOK_OUT_OF_STOCK", "book out of stock"))
		return
	}
	var book_student models.Book_Student
	book_student.BookID = book.ID
	book_student.StudentID = student.ID
	book_student.BorrowedAt = time.Now()
	book_student.Status = models.BorrowStatusBorrowed
	if err := h.DB.Create(&book_student).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	book.Stock -= 1
	if err := h.DB.Save(&book).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, book_student)
	// if book.StudentID != 0 {
	// 	var currentOwner models.Student
	// 	if err := h.DB.First(&currentOwner, book.StudentID).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"msg": "this book is already booked", "student": currentOwner})
	// 	return
	// }
	// book.StudentID = student.ID
	// if err := h.DB.Save(&book).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	// 	return
	// }
	// c.JSON(http.StatusOK, book)
}

// ReturnABook 归还书籍
// @Summary      归还书籍
// @Description  学生归还书籍
// @Tags         borrow
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        student_id  path      int  true  "学生 ID"
// @Param        book_id     path      int  true  "书籍 ID"
// @Success      200  {object}  models.Book_Student
// @Failure      400  {object}  middleware.AppError
// @Failure      404  {object}  middleware.AppError
// @Failure      500  {object}  middleware.AppError
// @Router       /students/{student_id}/books/{book_id}/return [post]
func (h *BookHandler) ReturnABook(c *gin.Context) {
	stuidstr := c.Param("student_id")
	bookidstr := c.Param("book_id")

	stuid, err := strconv.ParseUint(stuidstr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_STUDENT_ID", "invalid student_id"))
		return
	}
	bookid, err := strconv.ParseUint(bookidstr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_BOOK_ID", "invalid book_id"))
		return
	}

	var book_student models.Book_Student
	if err := h.DB.
		Where("student_id = ? AND book_id = ? AND status = ?", stuid, bookid, models.BorrowStatusBorrowed).
		First(&book_student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(middleware.NewAppError(http.StatusNotFound, "BORROW_RECORD_NOT_FOUND", "borrow record not found"))
			return
		}
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	book_student.Status = models.BorrowStatusReturned
	book_student.ReturnedAt = time.Now()
	if err := h.DB.Save(&book_student).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	var book models.Book
	if err := h.DB.First(&book, uint(bookid)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(middleware.NewAppError(http.StatusNotFound, "BOOK_NOT_FOUND", "book not found"))
			return
		}
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	book.Stock += 1
	if err := h.DB.Save(&book).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, book_student)
}

// ListStudentBooks 获取学生借书记录
// @Summary      获取学生借书记录
// @Description  获取指定学生的所有借阅记录（包含书籍信息）
// @Tags         borrow
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "学生 ID"
// @Success      200  {object}  models.Student
// @Failure      400  {object}  middleware.AppError
// @Failure      404  {object}  middleware.AppError
// @Router       /students/{id}/books [get]
func (h *BookHandler) ListStudentBooks(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_ID", "invalid id"))
		return
	}
	var student models.Student
	if err := h.DB.Preload("Book_Student").First(&student, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(middleware.NewAppError(http.StatusNotFound, "STUDENT_NOT_FOUND", "student not found"))
			return
		}
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
		return
	}
	c.JSON(http.StatusOK, student)
}
