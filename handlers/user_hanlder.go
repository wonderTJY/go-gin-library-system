package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	"trae-go/config"
	"trae-go/middleware"
	"trae-go/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewUserHanlder(db *gorm.DB, rdb *redis.Client) UserHandler {
	return UserHandler{db, rdb}
}

type UserRegisterRequest struct {
	Name      string `json:"user_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Sex       string `json:"sex"`
	BornDate  string `json:"born_date"`
	Identify  string `json:"ide"`
	AvatarURL string `json:"avatar_url"`
}
type UserUpdateRequest struct {
	Name      *string `json:"user_name"`
	Sex       *string `json:"sex"`
	BornDate  *string `json:"born_date"`
	AvatarURL *string `json:"avatar_url"`
}
type UserLoginRequest struct {
	Name     string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (h *UserHandler) UserRegister(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid json"))
		return
	}

	var existing models.User
	if err := h.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "USER_ALREADY_EXISTS", "user already exists"))
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal error"))
		return
	}

	var bornTime time.Time
	if req.BornDate != "" {
		t, err := time.Parse("2006-01-02", req.BornDate)
		if err != nil {
			c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_BORN_DATE", "invalid born_date"))
			return
		}
		bornTime = t
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "PASSWORD_HASH_FAILED", "password hash failed"))
		return
	}

	user := models.User{
		Name:      req.Name,
		Password:  string(hashedPassword),
		Sex:       req.Sex,
		BornDate:  bornTime,
		Identify:  req.Identify,
		AvatarURL: req.AvatarURL,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_CREATE_USER", "failed to create user"))
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid json"))
		return
	}

	var user models.User
	if err := h.DB.Where("name = ?", req.Name).First(&user).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid credentials"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.Error(middleware.NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid credentials"))
		return
	}

	token, err := generateToken(32)
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "TOKEN_GENERATE_FAILED", "token generate failed"))
		return
	}

	ctx := c.Request.Context()
	key := "auth:token:" + token
	d, err := time.ParseDuration(config.AppConfig.Auth.TokenExpireHours)
	if err != nil {
		log.Printf("TokenExpireHours parse failed,use default setting")
		d = 24 * time.Hour
	}
	if err := h.RDB.Set(ctx, key, strconv.Itoa(user.ID), d).Err(); err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "TOKEN_STORE_FAILED", "token store failed"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *UserHandler) UserDelte(c *gin.Context) {
	userName := c.Param("user_name")

	result := h.DB.Where("name = ?", userName).Delete(&models.User{})
	if result.Error != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_DELETE_USER", "failed to delete user"))
		return
	}
	if result.RowsAffected == 0 {
		c.Error(middleware.NewAppError(http.StatusNotFound, "USER_NOT_FOUND", "user not found"))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid json"))
		return
	}
	userIDval, ok := c.Get("user_id")
	if !ok {
		c.Error(middleware.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized"))
		return
	}
	userID, ok := userIDval.(uint)
	if !ok {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal error"))
		return
	}
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusNotFound, "USER_NOT_FOUND", "user not found"))
		return
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Sex != nil {
		user.Sex = *req.Sex
	}
	if req.BornDate != nil {
		if *req.BornDate == "" {
			user.BornDate = time.Time{}
		} else {
			t, err := time.Parse("2006-01-02", *req.BornDate)
			if err != nil {
				c.Error(middleware.NewAppError(http.StatusBadRequest, "INVALID_BORN_DATE", "invalid born_date"))
				return
			}
			user.BornDate = t
		}
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}

	if err := h.DB.Save(&user).Error; err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "FAILED_UPDATE_USER", "failed to update user"))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	path, err := SaveFile(c, "avatar", "static/avatars")
	if err != nil {
		c.Error(middleware.NewAppError(http.StatusInternalServerError, "SAVE_FILE_FAILED", "save file failed"))
		return
	}
	filename := filepath.Base(path)
	avatarURL := "http://127.0.0.1:8080/static/avatars/" + filename
	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}
