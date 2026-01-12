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
	studentHandler := handlers.NewStudentHandler(db)
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RequestCountMiddleware())
	r.Use(middleware.AuthenticationMiddleware())
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
	students.GET("", studentHandler.ListStudents)
	students.GET("/:id", studentHandler.GetStudent)
	students.POST("", studentHandler.CreatStudent)
	students.PUT("/:id", studentHandler.UpdateStudent)
	students.DELETE("/:id", studentHandler.DeleteStudent)
	students.GET("/:id/books", bookHandler.ListStudentBooks)
	students.POST("/:student_id/books/:book_id/borrow", bookHandler.BookABook)
	students.POST("/:student_id/books/:book_id/return", bookHandler.ReturnABook)

	return r
}
