package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
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

}
