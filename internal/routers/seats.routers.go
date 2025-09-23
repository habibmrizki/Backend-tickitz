package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/middlewares"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// InitSeatRouter menginisialisasi rute-rute yang berkaitan dengan kursi
func InitSeatRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	seatRepo := repositories.NewSeatRepository(db)
	seatHandler := handlers.NewSeatHandler(seatRepo)

	seatGroup := router.Group("/seats").Use(middlewares.VerifyTokenWithBlacklist(rdb))
	seatGroup.GET("/:scheduleId/available", middlewares.Access("user"), seatHandler.GetAvailableSeats)
	// seatGroup.GET("/seat/:seatId", seatHandler.GetSeatByID)
}
