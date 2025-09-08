package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitSeatRouter menginisialisasi rute-rute yang berkaitan dengan kursi
func InitSeatRouter(router *gin.Engine, db *pgxpool.Pool) {
	seatRepo := repositories.NewSeatRepository(db)
	seatHandler := handlers.NewSeatHandler(seatRepo)

	seatGroup := router.Group("/seats")
	seatGroup.GET("/:scheduleId/available", seatHandler.GetAvailableSeats)
	seatGroup.GET("/seat/:seatId", seatHandler.GetSeatByID)
}
