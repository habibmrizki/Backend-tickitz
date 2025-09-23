package utils

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context, formKey string, folder string, prefix string) (string, error) {
	file, err := ctx.FormFile(formKey)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), ext)
	location := filepath.Join(folder, filename)
	if err := ctx.SaveUploadedFile(file, location); err != nil {
		return "", err
	}
	return filename, nil
}
