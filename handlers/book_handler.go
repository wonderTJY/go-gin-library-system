package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"trae-go/models"
)

type BookHandler struct {
	DB *gorm.DB
}

func NewBookHandler(db *gorm.DB) *BookHandler {
	return &BookHandler{DB: db}
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	var books []models.Book
	if err := h.DB.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := h.DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get book"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	book := models.Book{
		Title:  input.Title,
		Author: input.Author,
	}
	if err := h.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}
	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := h.DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get book"})
		return
	}
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	book.Title = input.Title
	book.Author = input.Author
	book.Stock = input.Stock
	if err := h.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Book{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete book"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BookHandler) BookABook(c *gin.Context) {
	stuidstr := c.Param("student_id")
	bookidstr := c.Param("book_id")

	stuid, err := strconv.ParseUint(stuidstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student_id"})
		return
	}
	bookid, err := strconv.ParseUint(bookidstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book_id"})
		return
	}

	var student models.Student
	var book models.Book

	if err := h.DB.First(&student, uint(stuid)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if err := h.DB.First(&book, uint(bookid)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if book.Stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book out of stock"})
		return
	}
	var book_student models.Book_Student
	book_student.BookID = book.ID
	book_student.StudentID = student.ID
	book_student.BorrowedAt = time.Now()
	book_student.Status = models.BorrowStatusBorrowed
	if err := h.DB.Create(&book_student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	book.Stock -= 1
	if err := h.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
func (h *BookHandler) ReturnABook(c *gin.Context) {
	stuidstr := c.Param("student_id")
	bookidstr := c.Param("book_id")

	stuid, err := strconv.ParseUint(stuidstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student_id"})
		return
	}
	bookid, err := strconv.ParseUint(bookidstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book_id"})
		return
	}

	var book_student models.Book_Student
	if err := h.DB.
		Where("student_id = ? AND book_id = ? AND status = ?", stuid, bookid, models.BorrowStatusBorrowed).
		First(&book_student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "borrow record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	book_student.Status = models.BorrowStatusReturned
	book_student.ReturnedAt = time.Now()
	if err := h.DB.Save(&book_student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	var book models.Book
	if err := h.DB.First(&book, uint(bookid)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	book.Stock += 1
	if err := h.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, book_student)
}

func (h *BookHandler) ListStudentBooks(c *gin.Context) {
	id := c.Param("id")
	var student models.Student
	if err := h.DB.Preload("Book_Student").First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, student)
}
