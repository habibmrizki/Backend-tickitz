// package models

// import (
// 	"time"
// )

// type MovieStruct struct {
// 	Id           int        `json:"id"`
// 	DirectorId   int        `json:"director_id"`
// 	Title        string     `json:"title"`
// 	Synopsis     string     `json:"synopsis"`
// 	Popularity   int        `json:"popularity"`
// 	BackdropPath string     `json:"backdrop_path"`
// 	Genres       []string   `json:"genres" form:"genres" binding:"required"`
// 	Casts        []string   `json:"casts" form:"casts" binding:"required"`
// 	PosterPath   string     `json:"poster_path"`
// 	Duration     int        `json:"duration"`
// 	ReleaseDate  time.Time  `json:"release_date"`
// 	DirectorName string     `json:"director_name"`
// 	CreatedAt    *time.Time `json:"created_at"`
// 	UpdateAt     *time.Time `json:"update_at"`
// }

// // type MoviesStruct struct {
// // 	Id            int                   `json:"id,omitempty"`
// // 	Title         string                `json:"title" form:"title" binding:"required"`
// // 	Release_date  time.Time             `json:"release_date" form:"release_date" binding:"required"`
// // 	Overview      string                `json:"overview" form:"overview" binding:"required"`
// // 	Duration      int                   `json:"duration" form:"duration" binding:"required"`
// // 	Director_name string                `json:"director_name" form:"director_name" binding:"required"`
// // 	Genres        []string              `json:"genres" form:"genres" binding:"required"`
// // 	Casts         []string              `json:"casts" form:"casts" binding:"required"`
// // 	Image_movie   string                `json:"image_movie"`
// // 	TotalSales    int                   `json:"total_sales,omitempty"`
// // 	Image_path    *multipart.FileHeader `json:"image_path,omitempty" form:"image_path,omitempty" binding:"required"`
// // 	Cinema_ids    []int                 `json:"cinema_ids,omitempty" form:"cinema_ids,omitempty"`
// // 	Movie_id      int                   `json:"movie_id,omitempty" form:"movie_id,omitempty"`
// // 	Location      string                `json:"location,omitempty" form:"location,omitempty"`
// // 	Date          time.Time             `json:"-" form:"date,omitempty"`
// // 	Times         []time.Time           `json:"time,omitempty" form:"time,omitempty"`
// // 	Price         int                   `json:"price,omitempty" form:"price,omitempty"`
// // }

// type MovieListResponse struct {
// 	Message string        `json:"message"`
// 	Status  string        `json:"status"`
// 	Data    []MovieStruct `json:"data"`
// }

// type MovieDetailStruct struct {
// 	ID           int       `json:"id"`
// 	Title        string    `json:"title"`
// 	Synopsis     string    `json:"synopsis"`
// 	Duration     *int      `json:"duration"`
// 	Director     string    `json:"director"`
// 	Cast         []string  `json:"cast"`
// 	Genre        []string  `json:"genre"`
// 	ReleaseDate  time.Time `json:"release_date"`
// 	PosterPath   string    `json:"poster_path"`
// 	BackdropPath *string   `json:"backdrop_path"`
// }

// // MovieDetailResponse is the handler's response structure
// type MovieDetailResponse struct {
// 	Message string            `json:"message"`
// 	Status  string            `json:"status"`
// 	Data    MovieDetailStruct `json:"data"`
// }

// // MovieUpdateRequest is the request body for updating a movie
// type MovieUpdateRequest struct {
// 	Title        string    `json:"title"`
// 	Synopsis     string    `json:"synopsis"`
// 	Duration     int       `json:"duration"`
// 	DirectorId   int       `json:"director_id"`
// 	ReleaseDate  time.Time `json:"release_date"`
// 	PosterPath   string    `json:"poster_path"`
// 	BackdropPath string    `json:"backdrop_path"`
// }

// type PaginationInfo struct {
// 	TotalItems   int `json:"totalItems"`
// 	TotalPages   int `json:"totalPages"`
// 	CurrentPage  int `json:"currentPage"`
// 	ItemsPerPage int `json:"itemsPerPage"`
// }

// // PaginatedMovieListResponse mewakili respons yang dipaginasi untuk daftar film
//
//	type PaginatedMovieListResponse struct {
//		Status     string         `json:"status"`
//		Message    string         `json:"message"`
//		Data       []MovieStruct  `json:"data"`
//		Pagination PaginationInfo `json:"pagination"`
//	}
package models

import (
	"mime/multipart"
	"time"
)

// Define structs for Genre and Cast to hold both ID and name
type Genre struct {
	ID   int    `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

type Cast struct {
	ID   int    `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

// ini yang menggunakan id nantinya dalam bentuk json { id : 1, genre : "action"}
// type MovieStruct struct {
// 	Id           int        `json:"id"`
// 	DirectorId   int        `json:"director_id"`
// 	Title        string     `json:"title"`
// 	Synopsis     string     `json:"synopsis"`
// 	Popularity   int        `json:"popularity"`
// 	BackdropPath string     `json:"backdrop_path"`
// 	Genres       []Genre   `json:"genres" form:"genres" binding:"required"`
// 	Casts        []Cast   `json:"casts" form:"casts" binding:"required"`
// 	PosterPath   string     `json:"poster_path"`
// 	Duration     int        `json:"duration"`
// 	ReleaseDate  time.Time  `json:"release_date"`
// 	DirectorName string     `json:"director_name"`
// 	CreatedAt    *time.Time `json:"created_at,omitempty" form:"created,omitempty"`
// 	UpdateAt     *time.Time `json:"update_at,omitempty" form:"update,omiempty"`
// }

type MovieStruct struct {
	Id           int        `json:"id" db:"id"`
	DirectorId   int        `json:"director_id" db:"director_id"`
	Title        string     `json:"title" db:"title"`
	Synopsis     string     `json:"synopsis" db:"synopsis"`
	Popularity   int        `json:"popularity" db:"popularity"`
	BackdropPath string     `json:"backdrop_path" db:"backdrop_path"`
	Genres       []string   `json:"genres" db:"genres" form:"genres" binding:"required"`
	Casts        []string   `json:"casts" db:"casts" form:"casts" binding:"required"`
	PosterPath   string     `json:"poster_path" db:"poster_path"`
	Duration     int        `json:"duration" db:"duration"`
	ReleaseDate  time.Time  `json:"release_date" db:"release_date"`
	DirectorName string     `json:"director_name" db:"director_name"`
	CreatedAt    *time.Time `json:"created_at,omitempty" db:"created_at" form:"created,omitempty"`
	UpdateAt     *time.Time `json:"update_at,omitempty" db:"update_at" form:"update,omiempty"`
	ArchivedAt   *time.Time `json:"archived_at,omitempty" db:"-"`
}

type MovieListResponse struct {
	Message string        `json:"message"`
	Status  string        `json:"status"`
	Data    []MovieStruct `json:"data"`
}

// Update MovieDetailStruct to use the new Genre and Cast structs
type MovieDetailStruct struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Synopsis     string    `json:"synopsis"`
	Duration     *int      `json:"duration"`
	Director     string    `json:"director"`
	Cast         []Cast    `json:"cast"`  // Updated to use the new Cast struct
	Genre        []Genre   `json:"genre"` // Updated to use the new Genre struct
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
// type MovieUpdateRequest struct {
// 	Title        string    `json:"title"`
// 	Synopsis     string    `json:"synopsis"`
// 	Duration     int       `json:"duration"`
// 	DirectorId   int       `json:"director_id"`
// 	ReleaseDate  time.Time `json:"release_date"`
// 	PosterPath   string    `json:"poster_path"`
// 	BackdropPath string    `json:"backdrop_path"`
// }

// type MovieCreateRequest struct {
// 	Title        string    `json:"title" binding:"required"`
// 	Synopsis     string    `json:"synopsis" binding:"required"`
// 	Duration     int       `json:"duration" binding:"required"`
// 	Popularity   int       `json:"popularity"`
// 	DirectorId   int       `json:"director_id" binding:"required"`
// 	ReleaseDate  time.Time `json:"release_date" binding:"required"`
// 	PosterPath   string    `json:"-" binding:"required"`
// 	BackdropPath string    `json:"-"`
// 	Genre        []Genre   `json:"genre" binding:"required"`
// 	Cast         []Cast    `json:"cast" binding:"required"`
// }

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

type NewMovieRequest struct {
	Title        string                `json:"title" form:"title" binding:"required"`
	Synopsis     string                `json:"synopsis" form:"synopsis" binding:"required"`
	Duration     int                   `json:"duration" form:"duration" binding:"required"`
	DirectorName string                `json:"director" form:"director" binding:"required"`
	ReleaseDate  time.Time             `json:"release_date" form:"release_date" binding:"required" time_format:"2006-01-02"`
	PosterPath   *multipart.FileHeader `json:"-" form:"poster_path" binding:"required"`
	BackdropPath *multipart.FileHeader `json:"-" form:"backdrop_path" binding:"required"`
	Popularity   int                   `json:"popularity" form:"popularity"`
	Genres       []string              `json:"genres" form:"genres" binding:"required"`
	Casts        []string              `json:"casts" form:"casts" binding:"required"`
}

type GenreUp struct {
	ID int `form:"id"`
}

type CastUp struct {
	ID int `form:"id"`
}

// type MovieUpdateRequest struct {
// 	Title        *string               `form:"title"`
// 	Synopsis     *string               `form:"synopsis"`
// 	Duration     *int                  `form:"duration"`
// 	DirectorId   *int                  `form:"director_id"`
// 	ReleaseDate  *time.Time            `form:"release_date"`
// 	PosterPath   *multipart.FileHeader `form:"poster_path"`
// 	BackdropPath *multipart.FileHeader `form:"backdrop_path"`
// 	Popularity   *int                  `form:"popularity"`
// 	GenreIDs     []Genre               `form:"genre_ids"`
// 	CastIDs      []Cast                `form:"cast_ids"`
// }

type MovieUpdateRequest struct {
	Title        *string               `form:"title" json:"title,omitempty"`
	Synopsis     *string               `form:"synopsis" json:"synopsis,omitempty"`
	Duration     *int                  `form:"duration" json:"duration,omitempty"`
	Director     *string               `form:"director" json:"director,omitempty"`
	ReleaseDate  *time.Time            `form:"release_date" time_format:"2006-01-02" json:"release_date,omitempty"`
	PosterPath   *multipart.FileHeader `form:"poster_path" json:"-"`
	BackdropPath *multipart.FileHeader `form:"backdrop_path" json:"-"`
	Popularity   *int                  `form:"popularity" json:"popularity,omitempty"`
	Genres       []string              `form:"genres" json:"genres,omitempty"`
	Casts        []string              `form:"casts" json:"casts,omitempty"`
}

type MovieUpdateRequestResponse struct {
	Message string             `json:"message"`
	Status  string             `json:"status"`
	Data    MovieUpdateRequest `json:"data"`
}

type MovieArchived struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Archived_at time.Time `json:"archived_at"`
}

type MovieListAdminStruct struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Synopsis    string     `json:"synopsis"`
	PosterPath  string     `json:"poster_path"`
	Duration    int        `json:"duration"`
	ReleaseDate time.Time  `json:"release_date"`
	Genres      []string   `json:"genres"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
}

type MovieListAdminResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []MovieListAdminStruct `json:"data"`
}

type MovieGetDetailAdmin struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Synopsis     string    `json:"synopsis"`
	Duration     int       `json:"duration"`
	Director     string    `json:"director"`
	ReleaseDate  time.Time `json:"release_date"`
	PosterPath   string    `json:"poster_path"`
	BackdropPath string    `json:"backdrop_path"`
	Popularity   int       `json:"popularity"`
	Genres       []string  `json:"genres"`
	Casts        []string  `json:"casts"`
}

type MovieGetDetailResponseAdmin struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    MovieGetDetailAdmin `json:"data"`
}
