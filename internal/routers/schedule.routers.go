package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/middlewares"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// InitScheduleRouter menginisialisasi rute-rute yang berkaitan dengan jadwal
func InitScheduleRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	scheduleRepo := repositories.NewScheduleRepository(db)
	scheduleHandler := handlers.NewScheduleHandlers(scheduleRepo)

	scheduleGroup := router.Group("/schedules").Use(middlewares.VerifyTokenWithBlacklist(rdb))
	scheduleGroup.GET("", middlewares.Access("user"), scheduleHandler.GetAllSchedules)
	scheduleGroup.GET("/movie/:movieId", middlewares.Access("user", "admin"), scheduleHandler.GetSchedulesByMovieID)
	// scheduleGroup.GET("/:id", scheduleHandler.GetScheduleByMovieId)
}
