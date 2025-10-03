// package middlewares

// import (
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/habibmrizki/back-end-tickitz/pkg"
// )

// func VerifyToken(ctx *gin.Context) {
// 	// ambil token dari header
// 	bearerToken := ctx.GetHeader("Authorization")
// 	// Bearer token
// 	token := strings.Split(bearerToken, " ")[1]
// 	if token == "" {
// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"success": false,
// 			"error":   "Silahkan login terlebih dahulu",
// 		})
// 		return
// 	}

//		var claims pkg.Claims
//		if err := claims.VerifyToken(token); err != nil {
//			if err == jwt.ErrTokenInvalidIssuer {
//				log.Println("JWT Error.\nCause: ", err.Error())
//				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//					"success": false,
//					"error":   "Silahkan login kembali",
//				})
//				return
//			}
//			if err == jwt.ErrTokenExpired {
//				log.Println("JWT Error.\nCause: ", err.Error())
//				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//					"success": false,
//					"error":   "Silahkan login kembali",
//				})
//				return
//			}
//			log.Println("Internal Server Error.\nCause: ", err.Error())
//			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
//				"success": false,
//				"error":   "Internal Server Error",
//			})
//			return
//		}
//		ctx.Set("claims", claims)
//		ctx.Set("role", claims.Role)
//		ctx.Next()
//	}
package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/habibmrizki/back-end-tickitz/internal/configs"
	"github.com/habibmrizki/back-end-tickitz/pkg"
	"github.com/redis/go-redis/v9"
)

func VerifyToken(ctx *gin.Context) {
	// ambil token dari header
	bearerToken := ctx.GetHeader("Authorization")
	bearerToken = strings.TrimSpace(bearerToken)
	// Periksa jika header kosong atau tidak diawali dengan "Bearer "
	if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer ") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authorization header is missing or malformed",
		})
		return
	}

	// Ambil token dari header setelah "Bearer "
	token := strings.Split(bearerToken, " ")[1]

	var claims pkg.Claims
	if err := claims.VerifyToken(token); err != nil {
		if err == jwt.ErrTokenInvalidIssuer {
			log.Println("JWT Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Silahkan login kembali",
			})
			return
		}
		if err == jwt.ErrTokenExpired {
			log.Println("JWT Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Silahkan login kembali",
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal Server Error",
		})
		return
	}
	ctx.Set("claims", claims)
	ctx.Set("role", claims.Role)
	ctx.Next()
}

// Tambahkan fungsi ini di bawah fungsi VerifyToken yang sudah ada
func VerifyTokenWithBlacklist(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		bearerToken = strings.TrimSpace(bearerToken)
		if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header is missing or malformed",
			})
			return
		}
		token := strings.Split(bearerToken, " ")[1]

		// Cek blacklist token di Redis
		blacklisted, err := configs.IsTokenBlacklisted(rdb, token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal Server Error",
			})
			return
		}
		if blacklisted {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token sudah logout, silahkan login kembali",
			})
			return
		}

		// Lanjutkan dengan verifikasi token seperti kode Anda
		var claims pkg.Claims
		if err := claims.VerifyToken(token); err != nil {
			if err == jwt.ErrTokenInvalidIssuer {
				log.Println("JWT Error.\nCause: ", err.Error())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Silahkan login kembali",
				})
				return
			}
			if err == jwt.ErrTokenExpired {
				log.Println("JWT Error.\nCause: ", err.Error())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Silahkan login kembali",
				})
				return
			}
			log.Println("Internal Server Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal Server Error",
			})
			return
		}
		ctx.Set("claims", claims)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
