package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/middlewares"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// InitOrderRouter menginisialisasi rute-rute yang berkaitan dengan order
func InitOrderRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	orderRepo := repositories.NewOrderRepository(db)
	orderHandler := handlers.NewOrderHandlers(orderRepo)

	ordersGroup := router.Group("/orders").Use(middlewares.VerifyTokenWithBlacklist(rdb))
	// ordersGroup.POST("", middlewares.VerifyToken, middlewares.Access("user"), orderHandler.CreateOrder)
	ordersGroup.POST("", middlewares.Access("user"), orderHandler.CreateOrder)
	ordersGroup.GET("/history/:userId", middlewares.VerifyToken, middlewares.Access("user"), orderHandler.GetOrderHistory)
	ordersGroup.GET("/:orderId", middlewares.VerifyToken, middlewares.Access("user"), orderHandler.GetOrderByID)

}
