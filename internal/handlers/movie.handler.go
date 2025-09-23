package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
// @Tags                    Admin
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

	ctx.JSON(http.StatusOK, models.MovieListAdminResponse{
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
// func (m *MovieHandler) GetMoviesWithPagination(ctx *gin.Context) {
// 	pageStr := ctx.DefaultQuery("page", "1")
// 	limitStr := ctx.DefaultQuery("limit", "10")
// 	name := ctx.DefaultQuery("name", "") // Ambil parameter 'name' dari query

// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page <= 0 {
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "invalid page number",
// 		})
// 		return
// 	}

// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil || limit <= 0 {
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "invalid limit number",
// 		})
// 		return
// 	}

// 	movies, totalMovies, err := m.movieRepo.GetMoviesWithPagination(ctx.Request.Context(), limit, page, name)
// 	if err != nil {
// 		log.Println("[ERROR]: Failed to get all movies: ", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "failed to get all movies",
// 		})
// 		return
// 	}

// 	totalPages := (totalMovies + limit - 1) / limit

// 	ctx.JSON(http.StatusOK, models.PaginatedMovieListResponse{
// 		Status:  "berhasil",
// 		Message: "successfully retrieved all movies",
// 		Data:    movies,
// 		Pagination: models.PaginationInfo{
// 			TotalItems:   totalMovies,
// 			TotalPages:   totalPages,
// 			CurrentPage:  page,
// 			ItemsPerPage: limit,
// 		},
// 	})
// }

// GetMoviesWithPagination retrieves all movies with pagination and optional filtering
// @Summary      Get all movies with pagination and optional filtering
// @Description  Retrieves a paginated list of all movies, with optional filters by title and genre
// @Tags         Movies
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Number of items per page (default 10)"
// @Param        title query string false "Movie title to filter by"
// @Param        genre query string false "Genre name to filter by"
// @Accept       json
// @Produce      json
// @Success      200 {object} models.PaginatedMovieListResponse
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /movies [get]
// func (m *MovieHandler) GetMoviesWithPagination(ctx *gin.Context) {
// 	// pageStr := ctx.PostForm("page")
// 	// limitStr := ctx.PostForm("limit")
// 	// title := ctx.PostForm("title")
// 	// genreName := ctx.PostForm("genre")

// 	pageStr := ctx.Query("page")
// 	limitStr := ctx.Query("limit")
// 	title := ctx.Query("title")
// 	genreName := ctx.Query("genre")

// 	page, _ := strconv.Atoi(pageStr)
// 	if page <= 0 {
// 		page = 1
// 	}

// 	limit, _ := strconv.Atoi(limitStr)
// 	if limit <= 0 {
// 		limit = 10
// 	}

// 	movies, totalItems, err := m.movieRepo.GetMoviesWithPagination(
// 		ctx.Request.Context(),
// 		limit,
// 		page,
// 		title,
// 		genreName,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "error",
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	totalPages := (totalItems + limit - 1) / limit

// 	// c.JSON(http.StatusOK, gin.H{
// 	// 	"status":  "success",
// 	// 	"message": "movies retrieved successfully",
// 	// 	"data":    movies,
// 	// 	"pagination": gin.H{
// 	// 		"totalItems":   totalItems,
// 	// 		"totalPages":   totalPages,
// 	// 		"currentPage":  page,
// 	// 		"itemsPerPage": limit,
// 	// 	},
// 	// })

// 	pagination := models.PaginationInfo{
// 		TotalItems:   totalItems,
// 		TotalPages:   totalPages,
// 		CurrentPage:  page,
// 		ItemsPerPage: limit,
// 	}

// 	response := models.PaginatedMovieListResponse{
// 		Status:     "success",
// 		Message:    "movies retrieved successfully",
// 		Data:       movies,
// 		Pagination: pagination,
// 	}

// 	ctx.JSON(http.StatusOK, response)

// }

func (m *MovieHandler) GetMoviesWithPagination(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	title := ctx.Query("title")
	genresStr := ctx.Query("genre") // bisa banyak genre: "Action,Drama"

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 12 // ✅ samain dengan FE
	}

	var genres []string
	if genresStr != "" {
		parts := strings.Split(genresStr, ",")
		for _, g := range parts {
			trimmed := strings.TrimSpace(g)
			if trimmed != "" {
				genres = append(genres, trimmed)
			}
		}
	}

	movies, totalItems, err := m.movieRepo.GetMoviesWithPagination(
		ctx.Request.Context(),
		limit,
		page,
		title,
		genres,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	totalPages := (totalItems + limit - 1) / limit

	response := models.PaginatedMovieListResponse{
		Status:  "success",
		Message: "movies retrieved successfully",
		Data:    movies,
		Pagination: models.PaginationInfo{
			TotalItems:   totalItems,
			TotalPages:   totalPages,
			CurrentPage:  page,
			ItemsPerPage: limit,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteMovie deletes a movie by its ID (Admin only)
// @summary                 Delete a movie (Admin only)
// @router                  /admin/movies/:movieId [delete]
// @Description             Deletes a movie from the database by its ID (admin only)
// @Tags                    Admin
// @param                   movieId path int true "Movie ID"
// @accept                  json
// @produce                 json
// @security                ApiKeyAuth
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
// @failure                 500 {object} models.Response
// @success                 200 {object} models.Response
// func (m *MovieHandler) DeleteMovie(ctx *gin.Context) {
// 	movieIDStr := ctx.Param("movieId")
// 	log.Println("Received movie ID string:", movieIDStr)
// 	movieID, err := strconv.Atoi(movieIDStr)
// 	if err != nil {
// 		log.Println("Failed to convert ID:", err.Error())
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "invalid movie id",
// 		})
// 		return
// 	}

// 	err = m.movieRepo.DeleteMovie(ctx.Request.Context(), movieID)
// 	if err != nil {
// 		if err.Error() == "movie not found" {
// 			ctx.JSON(http.StatusNotFound, models.Response{
// 				Status:  "gagal",
// 				Message: "movie not found",
// 			})
// 			return
// 		}
// 		log.Println("[ERROR]: Failed to delete movie: ", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "internal server error",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, models.Response{
// 		Status:  "berhasil",
// 		Message: "successfully deleted movie",
// 	})
// }

func (h *MovieHandler) ArchiveMovie(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("movieId"))
	if err != nil {
		fmt.Println("error", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "Invalid match ID",
		})
		return
	}

	archivedMovie, err := h.movieRepo.ArchiveMovie(ctx.Request.Context(), id)
	if err != nil {
		fmt.Println("error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "Failed for deleted movie",
		})
		return
	}

	ctx.JSON(http.StatusOK, archivedMovie)
}

// UpdateMovie godoc
// @Summary      Update a movie (Admin only)
// @Description  Updates movie details, poster, and backdrop
// @Tags         Admin
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id             path     int    true "Movie ID"
// @Param        title          formData string false "Movie title"
// @Param        description    formData string false "Movie description"
// @Param        release_date   formData string false "Release date (YYYY-MM-DD)"
// @Param        poster         formData file   false "Poster image"
// @Param        backdrop       formData file   false "Backdrop image"
// @Param        genre_ids      formData []int  false "Genre IDs (comma separated)"
// @Param        cast_ids       formData []int  false "Cast IDs (comma separated)"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
// @Router 	/admin/movies/{id} [patch]
// func (m *MovieHandler) UpdateMovie(ctx *gin.Context) {
// 	movieIDStr := ctx.Param("id")
// 	movieID, err := strconv.Atoi(movieIDStr)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "ID film tidak valid",
// 		})
// 		return
// 	}

// 	var req models.MovieUpdateRequest
// 	if err := ctx.ShouldBind(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "Format data tidak valid: " + err.Error(),
// 		})
// 		return
// 	}

// 	// Tangani upload poster jika ada
// 	var posterPathString string
// 	if req.PosterPath != nil {
// 		ext := filepath.Ext(req.PosterPath.Filename)
// 		filename := fmt.Sprintf("poster_%d_%d%s", movieID, time.Now().UnixNano(), ext)
// 		location := filepath.Join("public", "images", filename)
// 		if err := ctx.SaveUploadedFile(req.PosterPath, location); err != nil {
// 			ctx.JSON(http.StatusInternalServerError, models.Response{
// 				Status:  "gagal",
// 				Message: "Gagal mengunggah poster",
// 			})
// 			return
// 		}
// 		posterPathString = "/img/" + filename
// 	}

// 	// Tangani upload backdrop jika ada
// 	var backdropPathString string
// 	if req.BackdropPath != nil {
// 		ext := filepath.Ext(req.BackdropPath.Filename)
// 		filename := fmt.Sprintf("backdrop_%d_%d%s", movieID, time.Now().UnixNano(), ext)
// 		location := filepath.Join("public", "images", filename)
// 		if err := ctx.SaveUploadedFile(req.BackdropPath, location); err != nil {
// 			ctx.JSON(http.StatusInternalServerError, models.Response{
// 				Status:  "gagal",
// 				Message: "Gagal mengunggah backdrop",
// 			})
// 			return
// 		}
// 		backdropPathString = "/img/" + filename
// 	}

// 	// Panggil repository untuk update
// 	err = m.movieRepo.UpdateMovie(ctx.Request.Context(), movieID, req, posterPathString, backdropPathString)
// 	if err != nil {
// 		if err.Error() == "movie not found" {
// 			ctx.JSON(http.StatusNotFound, models.Response{
// 				Status:  "gagal",
// 				Message: "Film tidak ditemukan",
// 			})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "Terjadi kesalahan internal",
// 		})
// 		return
// 	}

//		ctx.JSON(http.StatusOK, models.Response{
//			Status:  "berhasil",
//			Message: "Film berhasil diperbarui",
//		})
//	}
func (m *MovieHandler) UpdateMovie(ctx *gin.Context) {
	// --- Ambil param ID ---
	movieIDStr := ctx.Param("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "ID film tidak valid",
		})
		return
	}

	// --- Bind form ---
	var req models.MovieUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "Format data tidak valid: " + err.Error(),
		})
		return
	}

	// --- Handle upload poster ---
	var posterPathString string
	filePoster, err := ctx.FormFile("poster_path")
	if err == nil && filePoster != nil {
		ext := filepath.Ext(filePoster.Filename)
		filename := fmt.Sprintf("poster_%d_%d%s", movieID, time.Now().UnixNano(), ext)
		location := filepath.Join("public", "images", filename)

		if err := ctx.SaveUploadedFile(filePoster, location); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  "gagal",
				Message: "Gagal mengunggah poster",
			})
			return
		}
		posterPathString = "/img/" + filename
	}

	// --- Handle upload backdrop ---
	var backdropPathString string
	fileBackdrop, err := ctx.FormFile("backdrop_path")
	if err == nil && fileBackdrop != nil {
		ext := filepath.Ext(fileBackdrop.Filename)
		filename := fmt.Sprintf("backdrop_%d_%d%s", movieID, time.Now().UnixNano(), ext)
		location := filepath.Join("public", "images", filename)

		if err := ctx.SaveUploadedFile(fileBackdrop, location); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  "gagal",
				Message: "Gagal mengunggah backdrop",
			})
			return
		}
		backdropPathString = "/img/" + filename
	}

	// --- Panggil repository ---
	err = m.movieRepo.UpdateMovie(
		ctx.Request.Context(),
		movieID,
		req,
		posterPathString,
		backdropPathString,
	)
	if err != nil {
		log.Printf("❌ Error UpdateMovie Repo: %+v\n", err)
		if err.Error() == "movie not found" {
			ctx.JSON(http.StatusNotFound, models.Response{
				Status:  "gagal",
				Message: "Film tidak ditemukan",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "Terjadi kesalahan internal: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "berhasil",
		Message: "Film berhasil diperbarui",
	})
}

// AddNewMovie godoc
// @Summary      Create a new movie (Admin only)
// @Description  Creates a new movie with details, poster, and backdrop
// @Tags         Admin
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Param        title          formData string true "Movie title"
// @Param        description    formData string true "Movie description"
// @Param        release_date   formData string true "Release date (YYYY-MM-DD)"
// @Param        poster         formData file   true "Poster image"
// @Param        backdrop       formData file   true "Backdrop image"
// @Param        genre_ids      formData []int  true "Genre IDs (comma separated)"
// @Param        cast_ids       formData []int  true "Cast IDs (comma separated)"
// @Success      201 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /admin/movies [post]
// func (m *MovieHandler) AddNewMovie(ctx *gin.Context) {
// 	var req models.NewMovieRequest
// 	if err := ctx.ShouldBind(&req); err != nil {
// 		log.Println("[ERROR]: Gagal binding data -", err.Error())
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "Format data tidak valid: " + err.Error(),
// 		})
// 		return
// 	}

// 	// === Simpan Poster ===
// 	posterExt := filepath.Ext(req.PosterPath.Filename)
// 	posterFilename := fmt.Sprintf("poster_%s%s", time.Now().Format("20060102150405"), posterExt)
// 	posterLocation := filepath.Join("public", "images", posterFilename)

// 	if err := os.MkdirAll(filepath.Dir(posterLocation), 0755); err != nil {
// 		log.Println("[ERROR]: Gagal membuat direktori:", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "Gagal membuat direktori penyimpanan file",
// 		})
// 		return
// 	}

// 	if err := ctx.SaveUploadedFile(req.PosterPath, posterLocation); err != nil {
// 		log.Println("[ERROR]: Gagal upload poster:", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "Gagal mengunggah poster",
// 		})
// 		return
// 	}
// 	posterPathString := "/img/" + posterFilename

// 	// === Simpan Backdrop ===
// 	backdropExt := filepath.Ext(req.BackdropPath.Filename)
// 	backdropFilename := fmt.Sprintf("backdrop_%s%s", time.Now().Format("20060102150405"), backdropExt)
// 	backdropLocation := filepath.Join("public", "images", backdropFilename)

// 	if err := ctx.SaveUploadedFile(req.BackdropPath, backdropLocation); err != nil {
// 		log.Println("[ERROR]: Gagal upload backdrop:", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "Gagal mengunggah backdrop",
// 		})
// 		return
// 	}
// 	backdropPathString := "/img/" + backdropFilename

// 	// === Genre & Cast langsung pakai slice of int ===
// 	genreIDs := req.GenreIDs
// 	castIDs := req.CastIDs

// 	// === Simpan ke database ===
// 	err := m.movieRepo.AddNewMovie(ctx.Request.Context(), req, posterPathString, backdropPathString, genreIDs, castIDs)
// 	if err != nil {
// 		log.Println("[ERROR]: Gagal insert movie:", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "Gagal menambahkan film",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, models.Response{
// 		Status:  "berhasil",
// 		Message: "Film berhasil ditambahkan",
// 	})
// }

// func (m *MovieHandler) AddNewMovie(ctx *gin.Context) {
// 	var req models.NewMovieRequest
// 	if err := ctx.ShouldBind(&req); err != nil {
// 		log.Printf("[ERROR AddNewMovie][Handler][Binding]: %v\n", err)
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "Format data tidak valid: " + err.Error(),
// 		})
// 		return
// 	}

// 	// === Save Poster ===
// 	posterExt := filepath.Ext(req.PosterPath.Filename)
// 	posterFilename := fmt.Sprintf("poster_%d%s", time.Now().Unix(), posterExt)
// 	posterLocation := filepath.Join("public", "images", posterFilename)
// 	if err := os.MkdirAll(filepath.Dir(posterLocation), 0755); err != nil {
// 		log.Printf("[ERROR AddNewMovie][Handler][PosterDir]: %v\n", err)
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status: "gagal", Message: "Gagal membuat direktori poster",
// 		})
// 		return
// 	}
// 	if err := ctx.SaveUploadedFile(req.PosterPath, posterLocation); err != nil {
// 		log.Printf("[ERROR AddNewMovie][Handler][PosterUpload]: %v\n", err)
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status: "gagal", Message: "Gagal upload poster",
// 		})
// 		return
// 	}
// 	posterPath := "/img/" + posterFilename

// 	// === Save Backdrop ===
// 	backdropExt := filepath.Ext(req.BackdropPath.Filename)
// 	backdropFilename := fmt.Sprintf("backdrop_%d%s", time.Now().Unix(), backdropExt)
// 	backdropLocation := filepath.Join("public", "images", backdropFilename)
// 	if err := ctx.SaveUploadedFile(req.BackdropPath, backdropLocation); err != nil {
// 		log.Printf("[ERROR AddNewMovie][Handler][BackdropUpload]: %v\n", err)
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status: "gagal", Message: "Gagal upload backdrop",
// 		})
// 		return
// 	}
// 	backdropPath := "/img/" + backdropFilename

// 	// === Insert ke DB ===
// 	err := m.movieRepo.AddNewMovie(ctx.Request.Context(), req, posterPath, backdropPath)
// 	if err != nil {
// 		log.Printf("[ERROR AddNewMovie][Handler][Repo]: %v\n", err)
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "Gagal menambahkan film",
// 		})
// 		return
// 	}

// 	// kalau sukses
// 	log.Println("[INFO AddNewMovie][Handler]: Film berhasil ditambahkan")
// 	ctx.JSON(http.StatusCreated, models.Response{
// 		Status:  "berhasil",
// 		Message: "Film berhasil ditambahkan",
// 	})
// }

func (m *MovieHandler) AddNewMovie(ctx *gin.Context) {
	var req models.NewMovieRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("[ERROR AddNewMovie][Handler][Binding]: %v\n", err)
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "Format data tidak valid: " + err.Error(),
		})
		return
	}

	// === Save Poster ===
	posterExt := filepath.Ext(req.PosterPath.Filename)
	posterFilename := fmt.Sprintf("poster_%d%s", time.Now().Unix(), posterExt)
	posterLocation := filepath.Join("public", "images", posterFilename)
	if err := os.MkdirAll(filepath.Dir(posterLocation), 0755); err != nil {
		log.Printf("[ERROR AddNewMovie][Handler][PosterDir]: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status: "gagal", Message: "Gagal membuat direktori poster",
		})
		return
	}
	if err := ctx.SaveUploadedFile(req.PosterPath, posterLocation); err != nil {
		log.Printf("[ERROR AddNewMovie][Handler][PosterUpload]: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status: "gagal", Message: "Gagal upload poster",
		})
		return
	}
	posterPath := "/img/" + posterFilename

	// === Save Backdrop ===
	backdropExt := filepath.Ext(req.BackdropPath.Filename)
	backdropFilename := fmt.Sprintf("backdrop_%d%s", time.Now().Unix(), backdropExt)
	backdropLocation := filepath.Join("public", "images", backdropFilename)
	if err := ctx.SaveUploadedFile(req.BackdropPath, backdropLocation); err != nil {
		log.Printf("[ERROR AddNewMovie][Handler][BackdropUpload]: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status: "gagal", Message: "Gagal upload backdrop",
		})
		return
	}
	backdropPath := "/img/" + backdropFilename

	// === Insert ke DB ===
	err := m.movieRepo.AddNewMovie(ctx.Request.Context(), req, posterPath, backdropPath)
	if err != nil {
		log.Printf("[ERROR AddNewMovie][Handler][Repo]: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "Gagal menambahkan film",
		})
		return
	}

	log.Println("[INFO AddNewMovie][Handler]: Film berhasil ditambahkan")
	ctx.JSON(http.StatusCreated, models.Response{
		Status:  "berhasil",
		Message: "Film berhasil ditambahkan",
	})
}

func (m *MovieHandler) GetAdminMovieDetail(ctx *gin.Context) {
	movieIDStr := ctx.Param("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid movie id",
		})
		return
	}

	movie, err := m.movieRepo.GetAdminMovieDetail(ctx.Request.Context(), movieID)
	if err != nil {
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
		Status:  "berhasil",
		Message: "successfully retrieved admin movie detail",
		Data:    *movie,
	})
}
