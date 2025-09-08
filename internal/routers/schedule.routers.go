package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitScheduleRouter menginisialisasi rute-rute yang berkaitan dengan jadwal
func InitScheduleRouter(router *gin.Engine, db *pgxpool.Pool) {
	scheduleRepo := repositories.NewScheduleRepository(db)
	scheduleHandler := handlers.NewScheduleHandlers(scheduleRepo)

	scheduleGroup := router.Group("/schedules")
	scheduleGroup.GET("", scheduleHandler.GetAllSchedules)
	scheduleGroup.GET("/:id", scheduleHandler.GetScheduleByMovieId)
}
