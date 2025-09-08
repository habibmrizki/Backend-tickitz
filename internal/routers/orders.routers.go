package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitOrderRouter menginisialisasi rute-rute yang berkaitan dengan order
func InitOrderRouter(router *gin.Engine, db *pgxpool.Pool) {
	orderRepo := repositories.NewOrderRepository(db)
	orderHandler := handlers.NewOrderHandlers(orderRepo)

	ordersGroup := router.Group("/orders")
	ordersGroup.POST("", orderHandler.CreateOrder)
	ordersGroup.GET("/history/:userId", orderHandler.GetOrderHistory)
	ordersGroup.GET("/:orderId", orderHandler.GetOrderByID)
}
