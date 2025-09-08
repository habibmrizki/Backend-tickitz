package models

import (
	"time"
)

type MovieStruct struct {
	Id           int        `json:"id"`
	DirectorId   int        `json:"director_id"`
	Title        string     `json:"title"`
	Synopsis     string     `json:"synopsis"`
	Popularity   int        `json:"popularity"`
	BackdropPath string     `json:"backdrop_path"`
	Genres       []string   `json:"genres" form:"genres" binding:"required"`
	Casts        []string   `json:"casts" form:"casts" binding:"required"`
	PosterPath   string     `json:"poster_path"`
	Duration     int        `json:"duration"`
	ReleaseDate  time.Time  `json:"release_date"`
	DirectorName string     `json:"director_name"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdateAt     *time.Time `json:"update_at"`
}

// type MoviesStruct struct {
// 	Id            int                   `json:"id,omitempty"`
// 	Title         string                `json:"title" form:"title" binding:"required"`
// 	Release_date  time.Time             `json:"release_date" form:"release_date" binding:"required"`
// 	Overview      string                `json:"overview" form:"overview" binding:"required"`
// 	Duration      int                   `json:"duration" form:"duration" binding:"required"`
// 	Director_name string                `json:"director_name" form:"director_name" binding:"required"`
// 	Genres        []string              `json:"genres" form:"genres" binding:"required"`
// 	Casts         []string              `json:"casts" form:"casts" binding:"required"`
// 	Image_movie   string                `json:"image_movie"`
// 	TotalSales    int                   `json:"total_sales,omitempty"`
// 	Image_path    *multipart.FileHeader `json:"image_path,omitempty" form:"image_path,omitempty" binding:"required"`
// 	Cinema_ids    []int                 `json:"cinema_ids,omitempty" form:"cinema_ids,omitempty"`
// 	Movie_id      int                   `json:"movie_id,omitempty" form:"movie_id,omitempty"`
// 	Location      string                `json:"location,omitempty" form:"location,omitempty"`
// 	Date          time.Time             `json:"-" form:"date,omitempty"`
// 	Times         []time.Time           `json:"time,omitempty" form:"time,omitempty"`
// 	Price         int                   `json:"price,omitempty" form:"price,omitempty"`
// }

type MovieListResponse struct {
	Message string        `json:"message"`
	Status  string        `json:"status"`
	Data    []MovieStruct `json:"data"`
}

type MovieDetailStruct struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Synopsis     string    `json:"synopsis"`
	Duration     *int      `json:"duration"`
	Director     string    `json:"director"`
	Cast         []string  `json:"cast"`
	Genre        []string  `json:"genre"`
	ReleaseDate  time.Time `json:"release_date"`
	PosterPath   string    `json:"poster_path"`
	BackdropPath *string   `json:"backdrop_path"`
}

// MovieDetailResponse is the handler's response structure
type MovieDetailResponse struct {
	Message string            `json:"message"`
	Status  string            `json:"status"`
	Data    MovieDetailStruct `json:"data"`
}

// MovieUpdateRequest is the request body for updating a movie
type MovieUpdateRequest struct {
	Title        string    `json:"title"`
	Synopsis     string    `json:"synopsis"`
	Duration     int       `json:"duration"`
	DirectorId   int       `json:"director_id"`
	ReleaseDate  time.Time `json:"release_date"`
	PosterPath   string    `json:"poster_path"`
	BackdropPath string    `json:"backdrop_path"`
}

type PaginationInfo struct {
	TotalItems   int `json:"totalItems"`
	TotalPages   int `json:"totalPages"`
	CurrentPage  int `json:"currentPage"`
	ItemsPerPage int `json:"itemsPerPage"`
}

// PaginatedMovieListResponse mewakili respons yang dipaginasi untuk daftar film
type PaginatedMovieListResponse struct {
	Status     string         `json:"status"`
	Message    string         `json:"message"`
	Data       []MovieStruct  `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
