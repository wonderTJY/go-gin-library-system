package handlers

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveFile(c *gin.Context, fileName, dir string) (string, error) {
	file, err := c.FormFile(fileName)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	ext := filepath.Ext(file.Filename)
	now := time.Now().UnixNano()
	r := rand.Int63()
	newName := fmt.Sprintf("%d_%d%s", now, r, ext)
	dst := filepath.Join(dir, newName)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	return dst, nil
}
