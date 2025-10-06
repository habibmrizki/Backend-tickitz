package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// InitMovieRouter menginisialisasi semua rute yang berkaitan dengan film
func InitMovieRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	movieRepo := repositories.NewMovieRepository(db, rdb)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	movieGroup := router.Group("/movies")
	movieGroup.GET("/upcoming", movieHandler.GetUpcomingMovies)
	movieGroup.GET("/popular", movieHandler.GetPopularMovies)
	movieGroup.GET("", movieHandler.GetMoviesWithPagination)

	movieGroup.GET("/:movieId", movieHandler.GetMovieDetail)
	genreGroup := router.Group("/genres")
	genreGroup.GET("", movieHandler.GetGenres)
}
