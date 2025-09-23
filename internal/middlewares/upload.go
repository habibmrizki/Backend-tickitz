package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateFileSize(maxSize int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
			return
		}
		if file.Size > maxSize {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "File terlalu besar"})
			return
		}
		ctx.Next()
	}
}
