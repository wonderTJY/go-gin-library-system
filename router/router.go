package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"trae-go/handlers"
	"trae-go/middleware"
)

func SetupRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	bookHandler := handlers.NewBookHandler(db)
	studentHandler := handlers.NewStudentHandler(db)
	userHanlder := handlers.NewUserHanlder(db, rdb)

	r.Static("/static/avatars", "./static/avatars")

	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RequestCountMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RedisRateLimiterMiddleware(rdb))
	r.Use(middleware.CorsMiddleware([]string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}))
	r.Use(middleware.ErrorHandlingMiddleware())

	api := r.Group("/api")
	v1 := api.Group("/v1")

	// 不需要登录的接口（公开接口）
	publicUser := v1.Group("/user")
	publicUser.POST("/register", userHanlder.UserRegister)
	publicUser.POST("/login", userHanlder.UserLogin)
	publicUser.POST("/uploadAvatar", userHanlder.UploadAvatar)

	// 需要登录的接口
	authRequired := v1.Group("")
	authRequired.Use(middleware.AuthenticationMiddleware(db, rdb))

	authUser := authRequired.Group("/user")
	authUser.PUT("/profile", userHanlder.UpdateUser)
	authUser.DELETE("/:user_name", userHanlder.UserDelte)

	books := authRequired.Group("/books")
	books.GET("", bookHandler.ListBooks)
	books.GET("/:id", bookHandler.GetBook)
	books.POST("", bookHandler.CreateBook)
	books.PUT("/:id", bookHandler.UpdateBook)
	books.DELETE("/:id", bookHandler.DeleteBook)

	students := authRequired.Group("/students")
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
