package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/middlewares"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitAdminMovieRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	movieRepo := repositories.NewMovieRepository(db, rdb)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	adminMovieRouter := router.Group("/admin/movies", middlewares.VerifyTokenWithBlacklist(rdb), middlewares.AdminOnly)
	adminMovieRouter.POST("", movieHandler.AddNewMovie)
	adminMovieRouter.GET("", movieHandler.GetAllMovies)
	adminMovieRouter.DELETE("/:movieId/archive", movieHandler.ArchiveMovie)
	adminMovieRouter.PATCH("/:id", movieHandler.UpdateMovie)
	adminMovieRouter.GET("/:id", movieHandler.GetAdminMovieDetail)
	// adminMovieRouter.DELETE("/:movieId", movieHandler.DeleteMovie)

}
