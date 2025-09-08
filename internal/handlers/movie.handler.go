package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
)

type MovieHandler struct {
	movieRepo *repositories.MovieRepository
}

func NewMovieHandler(movieRepo *repositories.MovieRepository) *MovieHandler {
	return &MovieHandler{movieRepo: movieRepo}
}

// GetUpcomingMovies
// @summary                 Get upcoming movies
// @router                  /movies/upcoming [get]
// @Description             Get a list of upcoming movies
// @Tags                    Movies
// @accept                  json
// @produce                 json
// @success                 200 {object} models.MovieListResponse
// @failure                 500 {object} models.Response
func (m *MovieHandler) GetUpcomingMovies(ctx *gin.Context) {
	movies, err := m.movieRepo.GetUpcomingMovies(ctx.Request.Context())
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.MovieListResponse{
		Message: "successfully retrieve upcoming movies",
		Status:  "berhasil",
		Data:    movies,
	})
}

// GetPopularMovies
// @summary                 Get popular movies
// @router                  /movies/popular [get]
// @Description             Get a list of popular movies
// @Tags                    Movies
// @accept                  json
// @produce                 json
// @success                 200 {object} models.MovieListResponse
// @failure                 500 {object} models.Response
func (m *MovieHandler) GetPopularMovies(ctx *gin.Context) {
	movies, err := m.movieRepo.GetPopularMovies(ctx.Request.Context())
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.MovieListResponse{
		Message: "successfully retrieve popular movies",
		Status:  "berhasil",
		Data:    movies,
	})
}

// GetMovieDetail
// @summary                 Get movie detail
// @router                  /movies/:movieId [get]
// @Description             Get a single movie's details by its ID
// @Tags                    Movies
// @accept                  json
// @produce                 json
// @param                   movieId path int true "Movie ID"
// @success                 200 {object} models.MovieDetailResponse
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
// @failure                 500 {object} models.Response
func (m *MovieHandler) GetMovieDetail(ctx *gin.Context) {
	movieIDStr := ctx.Param("movieId")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid movie id",
		})
		return
	}

	movie, err := m.movieRepo.GetMovieDetail(ctx.Request.Context(), movieID)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "internal server error",
		})
		return
	}

	if movie == nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Status:  "gagal",
			Message: "movie not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.MovieDetailResponse{
		Message: "successfully retrieve movie details",
		Status:  "berhasil",
		Data:    *movie,
	})
}

// GetAllMovies retrieves all movies (Admin only)
// @summary                 Get all movies (Admin only)
// @router                  /admin/movies [get]
// @Description             Retrieves a list of all movies for admin
// @Tags                    Admin/Movies
// @accept                  json
// @produce                 json
// @security                ApiKeyAuth
// @success                 200 {object} models.MovieListResponse
// @failure                 500 {object} models.Response
func (m *MovieHandler) GetAllMovies(ctx *gin.Context) {
	movies, err := m.movieRepo.GetAllMovies(ctx.Request.Context())
	if err != nil {
		log.Println("[ERROR]: Failed to get all movies: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "failed to get all movies",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.MovieListResponse{
		Status:  "berhasil",
		Message: "successfully retrieved all movies",
		Data:    movies,
	})
}

// GetMoviesWithPagination retrieves all movies with pagination and optional filtering
// @summary                 Get all movies with pagination and optional filtering
// @router                  /movies [get]
// @Description             Retrieves a paginated list of all movies, with an optional filter by title
// @Tags                    Movies
// @param                   page query int false "Page number (default 1)"
// @param                   limit query int false "Number of items per page (default 10)"
// @param                   name query string false "Movie title to filter by"
// @accept                  json
// @produce                 json
// @success                 200 {object} models.PaginatedMovieListResponse
// @failure                 400 {object} models.Response
// @failure                 500 {object} models.Response
func (m *MovieHandler) GetMoviesWithPagination(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	name := ctx.DefaultQuery("name", "") // Ambil parameter 'name' dari query

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid page number",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid limit number",
		})
		return
	}

	movies, totalMovies, err := m.movieRepo.GetMoviesWithPagination(ctx.Request.Context(), limit, page, name)
	if err != nil {
		log.Println("[ERROR]: Failed to get all movies: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "failed to get all movies",
		})
		return
	}

	totalPages := (totalMovies + limit - 1) / limit

	ctx.JSON(http.StatusOK, models.PaginatedMovieListResponse{
		Status:  "berhasil",
		Message: "successfully retrieved all movies",
		Data:    movies,
		Pagination: models.PaginationInfo{
			TotalItems:   totalMovies,
			TotalPages:   totalPages,
			CurrentPage:  page,
			ItemsPerPage: limit,
		},
	})
}

// UpdateMovie updates an existing movie by its ID (Admin only)
// @summary                 Update a movie (Admin only)
// @router                  /admin/movies/:movieId [put]
// @Description             Updates a movie's details by its ID (admin only)
// @Tags                    Admin/Movies
// @param                   movieId path int true "Movie ID"
// @param                   movie body models.MovieUpdateRequest true "Updated movie details"
// @accept                  json
// @produce                 json
// @security                ApiKeyAuth
// @success                 200 {object} models.Response
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
// @failure                 500 {object} models.Response
func (m *MovieHandler) UpdateMovie(ctx *gin.Context) {
	movieIDStr := ctx.Param("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid movie ID",
		})
		return
	}

	var movie models.MovieUpdateRequest
	if err := ctx.ShouldBindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid request body",
		})
		return
	}

	err = m.movieRepo.UpdateMovie(ctx.Request.Context(), movieID, movie)
	if err != nil {
		if err.Error() == "movie not found" {
			ctx.JSON(http.StatusNotFound, models.Response{
				Status:  "gagal",
				Message: "movie not found",
			})
			return
		}
		log.Println("[ERROR]: Failed to update movie: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "berhasil",
		Message: "successfully updated movie",
	})
}

// DeleteMovie deletes a movie by its ID (Admin only)
// @summary                 Delete a movie (Admin only)
// @router                  /admin/movies/:movieId [delete]
// @Description             Deletes a movie from the database by its ID (admin only)
// @Tags                    Admin/Movies
// @param                   movieId path int true "Movie ID"
// @accept                  json
// @produce                 json
// @security                ApiKeyAuth
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
// @failure                 500 {object} models.Response
// @success                 200 {object} models.Response
func (m *MovieHandler) DeleteMovie(ctx *gin.Context) {
	movieIDStr := ctx.Param("movieId")
	log.Println("Received movie ID string:", movieIDStr)
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		log.Println("Failed to convert ID:", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid movie id",
		})
		return
	}

	err = m.movieRepo.DeleteMovie(ctx.Request.Context(), movieID)
	if err != nil {
		if err.Error() == "movie not found" {
			ctx.JSON(http.StatusNotFound, models.Response{
				Status:  "gagal",
				Message: "movie not found",
			})
			return
		}
		log.Println("[ERROR]: Failed to delete movie: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "berhasil",
		Message: "successfully deleted movie",
	})
}
