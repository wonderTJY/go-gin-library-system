package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"trae-go/handlers"
	"trae-go/middleware"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	bookHandler := handlers.NewBookHandler(db)

	r.Use(middleware.LoggingMiddleware())

	api := r.Group("/api")
	v1 := api.Group("/v1")

	books := v1.Group("/books")
	books.GET("", bookHandler.ListBooks)
	books.GET("/:id", bookHandler.GetBook)
	books.POST("", bookHandler.CreateBook)
	books.PUT("/:id", bookHandler.UpdateBook)
	books.DELETE("/:id", bookHandler.DeleteBook)

	students := v1.Group("/students")
	students.GET("/:id/books", bookHandler.ListStudentBooks)
	students.POST("/:student_id/books/:book_id/borrow", bookHandler.BookABook)
	students.POST("/:student_id/books/:book_id/return", bookHandler.ReturnABook)

	return r
}
