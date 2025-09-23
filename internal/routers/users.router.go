package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/middlewares"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// InitUsersRouter menginisialisasi semua rute yang berkaitan dengan pengguna
func InitUsersRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	usersGroup := router.Group("/auth")

	// Inisialisasi Repositori dan Handler
	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.NewUsersHandlers(userRepo)

	// Tambahkan rute register ke dalam grup
	usersGroup.POST("/register", userHandler.UserRegister)
	usersGroup.POST("/login", userHandler.UserLogin)

	usersGroup.GET("/:id", userHandler.GetProfileById)
	// usersGroup.PATCH("/:id", userHandler.UpdateProfile)

	// Grup rute yang memerlukan otentikasi JWT.
	// Middleware `VerifyToken` akan memeriksa token JWT dan menyimpan klaim (termasuk peran) di konteks.
	apiGroup := router.Group("/users").Use(middlewares.VerifyTokenWithBlacklist(rdb))

	apiGroup.GET("", middlewares.Access("user"), userHandler.GetProfileById)

	// Middleware `middlewares.Access("user")` membatasi akses hanya untuk admin.
	// apiGroup.PATCH("/:id", middlewares.Access("user"), userHandler.UpdateProfile)
	apiGroup.PATCH("", middlewares.Access("user"), userHandler.UpdateProfile)

	// apiGroup.GET("/dashboard", userHandler.GetDashboard)
	apiGroup.POST("/logout", handlers.LogoutHandler(rdb))

}
