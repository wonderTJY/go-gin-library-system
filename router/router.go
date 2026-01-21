package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"trae-go/handlers"
	"trae-go/middleware"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	bookHandler := handlers.NewBookHandler(db)
	studentHandler := handlers.NewStudentHandler(db)
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RequestCountMiddleware())

	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RateLimiterMiddleware())
	r.Use(middleware.CorsMiddleware([]string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}))
	r.Use(middleware.ErrorHandlingMiddleware())
	r.Use(middleware.AuthenticationMiddleware(db))

	api := r.Group("/api")
	v1 := api.Group("/v1")
	//auth := r.Group("/auth").Use(middleware.AuthenticationMiddleware(db))

	books := v1.Group("/books")
	books.GET("", bookHandler.ListBooks)
	books.GET("/:id", bookHandler.GetBook)
	books.POST("", bookHandler.CreateBook)
	books.PUT("/:id", bookHandler.UpdateBook)
	books.DELETE("/:id", bookHandler.DeleteBook)

	students := v1.Group("/students")
	students.GET("/panic", studentHandler.PanicTest)
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
