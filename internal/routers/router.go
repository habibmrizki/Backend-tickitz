package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/docs"
	"github.com/habibmrizki/back-end-tickitz/internal/middlewares"
	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.MyLogger)
	router.Use(middlewares.CORSMiddleware)

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Catch all route
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, models.Response{
			Message: "Rute Salah",
			Status:  "Rute Tidak Ditemukan",
		})
	})

	// Panggil fungsi inisialisasi rute-rute khusus
	InitUsersRouter(router, db)
	InitMovieRouter(router, db)
	InitScheduleRouter(router, db)
	InitOrderRouter(router, db)
	InitSeatRouter(router, db)
	return router
}
