package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/middlewares"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitUsersRouter menginisialisasi semua rute yang berkaitan dengan pengguna
func InitUsersRouter(router *gin.Engine, db *pgxpool.Pool) {
	usersGroup := router.Group("/auth")

	// Inisialisasi Repositori dan Handler
	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.NewUsersHandlers(userRepo)

	// Tambahkan rute register ke dalam grup
	usersGroup.POST("/register", userHandler.UserRegister)
	usersGroup.POST("/login", userHandler.UserLogin)

	usersGroup.GET("/:id", userHandler.GetProfileById)
	usersGroup.PUT("/:id", userHandler.UpdateProfile)

	// Grup rute yang memerlukan otentikasi JWT.
	// Middleware `VerifyToken` akan memeriksa token JWT dan menyimpan klaim (termasuk peran) di konteks.
	apiGroup := router.Group("/users").Use(middlewares.VerifyToken)

	// Rute ini dapat diakses oleh "user" dan "admin".
	// Middleware `middlewares.Access("user", "admin")` memastikan hanya user dengan salah satu dari peran ini yang bisa mengakses.
	apiGroup.GET("/:id", middlewares.Access("user", "admin"), userHandler.GetProfileById)

	// Rute ini hanya dapat diakses oleh "admin".
	// Middleware `middlewares.Access("admin")` membatasi akses hanya untuk admin.
	apiGroup.PUT("/:id", middlewares.Access("admin"), userHandler.UpdateProfile)

}
