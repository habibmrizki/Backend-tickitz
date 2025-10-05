package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/docs"
	"github.com/habibmrizki/back-end-tickitz/internal/middlewares"
	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.MyLogger)
	router.Use(middlewares.CORSMiddleware)

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// router.Static("/img", "./public/images")
	router.Static("/img", "public/images")
	// Panggil fungsi inisialisasi rute-rute khusus
	InitUsersRouter(router, db, rdb)
	InitAdminMovieRouter(router, db, rdb)
	InitMovieRouter(router, db, rdb)
	InitScheduleRouter(router, db, rdb)
	InitOrderRouter(router, db, rdb)
	InitSeatRouter(router, db, rdb)

	// Catch all route
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, models.Response{
			Message: "Rute Salah",
			Status:  "Rute Tidak Ditemukan",
		})
	})

	return router
}
