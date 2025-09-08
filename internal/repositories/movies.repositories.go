package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/habibmrizki/back-end-tickitz/internal/models"
	// "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{db: db}
}

// GetUpcomingMovies mengambil daftar film yang akan datang dari database
func (m *MovieRepository) GetUpcomingMovies(ctx context.Context) ([]models.MovieStruct, error) {
	sql := `
    SELECT id, title, poster_path, COALESCE(backdrop_path, '') AS backdrop_path, release_date
    FROM movie
    WHERE release_date > CURRENT_DATE
    ORDER BY release_date ASC;
`

	rows, err := m.db.Query(ctx, sql)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieStruct
	for rows.Next() {
		var movie models.MovieStruct
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.PosterPath, &movie.BackdropPath, &movie.ReleaseDate); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MovieRepository) GetPopularMovies(ctx context.Context) ([]models.MovieStruct, error) {
	sql := `
        SELECT
            m.id, m.director_id, m.title, m.synopsis, m.popularity,
			COALESCE(m.backdrop_path, '') AS backdrop_path, -- Tambahkan COALESCE di sini
			m.poster_path, COALESCE(m.duration, 0) AS duration, m.release_date, m.created_at, m.update_at, d.name AS director_name
        FROM movie m
        JOIN director d ON m.director_id = d.id
        ORDER BY m.popularity DESC;
    `

	rows, err := m.db.Query(ctx, sql)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieStruct
	for rows.Next() {
		var movie models.MovieStruct
		if err := rows.Scan(
			&movie.Id,
			&movie.DirectorId,
			&movie.Title,
			&movie.Synopsis,
			&movie.Popularity,
			&movie.BackdropPath,
			&movie.PosterPath,
			&movie.Duration,
			&movie.ReleaseDate,
			&movie.CreatedAt,
			&movie.UpdateAt,
			&movie.DirectorName,
		); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MovieRepository) GetMovieDetail(ctx context.Context, movieID int) (*models.MovieDetailStruct, error) {
	// First query: Get main movie details
	sql := `
		SELECT
			mv.id,
			mv.title,
			mv.synopsis,
			mv.duration,
			d.name director_name,
			mv.release_date,
			mv.poster_path,
			mv.backdrop_path
		FROM
			movie mv
		JOIN
			director d ON mv.director_id = d.id
		WHERE
			mv.id = $1
	`

	var movie models.MovieDetailStruct
	err := m.db.QueryRow(ctx, sql, movieID).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Synopsis,
		&movie.Duration,
		&movie.Director,
		&movie.ReleaseDate,
		&movie.PosterPath,
		&movie.BackdropPath,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}

	// Second query: Get the genres for the movie
	genreSql := `
		SELECT
			g.name
		FROM
			genre g
		JOIN
			movies_genre mg ON g.id = mg.genre_id
		WHERE
			mg.movie_id = $1
	`
	genreRows, err := m.db.Query(ctx, genreSql, movieID)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}
	defer genreRows.Close()

	var genres []string
	for genreRows.Next() {
		var genreName string
		if err := genreRows.Scan(&genreName); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return nil, err
		}
		genres = append(genres, genreName)
	}
	movie.Genre = genres

	// Third query: Get the cast for the movie
	castSql := `
		SELECT
			c.name
		FROM
			"cast" c
		JOIN
			movie_cast mc ON c.id = mc.cast_id
		WHERE
			mc.movie_id = $1
		ORDER BY c.name ASC
	`

	castRows, err := m.db.Query(ctx, castSql, movieID)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}
	defer castRows.Close()

	var castList []string
	for castRows.Next() {
		var castName string
		if err := castRows.Scan(&castName); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return nil, err
		}
		castList = append(castList, castName)
	}
	movie.Cast = castList

	return &movie, nil
}

// GetMoviesWithPagination mengambil semua film dengan paginasi
func (m *MovieRepository) GetMoviesWithPagination(ctx context.Context, limit, page int) ([]models.MovieStruct, int, error) {
	offset := (page - 1) * limit

	var totalMovies int
	totalSql := `SELECT COUNT(*) FROM movie;`
	err := m.db.QueryRow(ctx, totalSql).Scan(&totalMovies)
	if err != nil {
		log.Println("[ERROR]: Failed to get total movies count: ", err.Error())
		return nil, 0, err
	}

	sql := `
		SELECT
			id, title, synopsis, poster_path, COALESCE(duration, 0) AS duration, release_date
		FROM
			movie
		ORDER BY
			created_at DESC
		LIMIT $1 OFFSET $2;
	`
	rows, err := m.db.Query(ctx, sql, limit, offset)
	if err != nil {
		log.Println("[ERROR]: Failed to get all movies with pagination: ", err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.MovieStruct
	for rows.Next() {
		var movie models.MovieStruct
		if err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Synopsis,
			&movie.PosterPath,
			&movie.Duration,
			&movie.ReleaseDate,
		); err != nil {
			log.Println("[ERROR]: Failed to scan movie row: ", err.Error())
			return nil, 0, err
		}
		movies = append(movies, movie)
	}

	return movies, totalMovies, nil
}

// GetAllMovies retrieves all movies from the database
func (m *MovieRepository) GetAllMovies(ctx context.Context) ([]models.MovieStruct, error) {
	sql := `
		SELECT
			id, title, synopsis, poster_path, COALESCE(duration, 0) AS duration, release_date
		FROM
			movie
		ORDER BY
			created_at DESC;
	`
	rows, err := m.db.Query(ctx, sql)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieStruct
	for rows.Next() {
		var movie models.MovieStruct
		if err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Synopsis,
			&movie.PosterPath,
			&movie.Duration, // Baris ini sekarang akan aman
			&movie.ReleaseDate,
		); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return nil, err
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// UpdateMovie updates a movie in the database
func (m *MovieRepository) UpdateMovie(ctx context.Context, movieID int, movie models.MovieUpdateRequest) error {
	sql := `
		UPDATE movie
		SET
			title = $1,
			synopsis = $2,
			duration = $3,
			release_date = $4,
			director_id = $5,
			poster_path = $6,
			backdrop_path = $7,
		update_at = $8 
		WHERE
			id = $9

	`
	cmd, err := m.db.Exec(
		ctx,
		sql,
		movie.Title,
		movie.Synopsis,
		movie.Duration,
		movie.ReleaseDate,
		movie.DirectorId,
		movie.PosterPath,
		movie.BackdropPath,
		time.Now(),
		movieID,
	)
	if err != nil {
		log.Println("[ERROR]: Failed to update movie: ", err.Error())
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("movie not found")
	}

	return nil
}

// DeleteMovie deletes a movie from the database
func (m *MovieRepository) DeleteMovie(ctx context.Context, movieID int) error {
	sql := `DELETE FROM movie WHERE id = $1`
	cmd, err := m.db.Exec(ctx, sql, movieID)
	if err != nil {
		log.Println("[ERROR]: Failed to delete movie: ", err.Error())
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("movie not found")
	}
	return nil
}
