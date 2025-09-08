package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/handlers"
	"github.com/habibmrizki/back-end-tickitz/internal/middlewares"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitMovieRouter menginisialisasi semua rute yang berkaitan dengan film
func InitMovieRouter(router *gin.Engine, db *pgxpool.Pool) {
	movieRepo := repositories.NewMovieRepository(db)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	movieGroup := router.Group("/movies")
	movieGroup.GET("/upcoming", movieHandler.GetUpcomingMovies)
	movieGroup.GET("/popular", movieHandler.GetPopularMovies)
	movieGroup.GET("", movieHandler.GetMoviesWithPagination)

	movieGroup.GET("/:movieId", movieHandler.GetMovieDetail)

	adminMovieRouter := router.Group("/admin/movies", middlewares.VerifyToken, middlewares.AdminOnly)
	adminMovieRouter.GET("", movieHandler.GetAllMovies)
	adminMovieRouter.PUT("/:id", movieHandler.UpdateMovie)
	adminMovieRouter.DELETE("/:movieId", movieHandler.DeleteMovie)
}
