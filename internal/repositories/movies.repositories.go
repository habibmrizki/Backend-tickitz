package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/redis/go-redis/v9"

	// "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewMovieRepository(db *pgxpool.Pool, rdb *redis.Client) *MovieRepository {
	return &MovieRepository{db: db, rdb: rdb}
}

func (m *MovieRepository) GetMovieDetail(ctx context.Context, movieID int) (*models.MovieDetailStruct, error) {
	sql := `
        SELECT
            mv.id,
            mv.title,
            mv.synopsis,
            mv.duration,
            d.name AS director_name,
            mv.release_date,
            mv.poster_path,
            mv.backdrop_path,
            COALESCE(json_agg(DISTINCT jsonb_build_object('id', g.id, 'name', g.name)) FILTER (WHERE g.id IS NOT NULL), '[]') AS genres,
            COALESCE(json_agg(DISTINCT jsonb_build_object('id', c.id, 'name', c.name)) FILTER (WHERE c.id IS NOT NULL), '[]') AS cast_list
        FROM
            movie mv
        JOIN
            director d ON mv.director_id = d.id
        LEFT JOIN
            movies_genre mg ON mv.id = mg.movie_id
        LEFT JOIN
            genre g ON mg.genre_id = g.id
        LEFT JOIN
            movie_cast mc ON mv.id = mc.movie_id
        LEFT JOIN
            "cast" c ON mc.cast_id = c.id
        WHERE
            mv.id = $1
        GROUP BY
            mv.id, d.name;
    `
	// AND mv.archived_at IS NULL
	var movie models.MovieDetailStruct
	var genresJSON []byte   // Use byte slice to scan JSON from the query
	var castListJSON []byte // Use byte slice to scan JSON from the query

	err := m.db.QueryRow(ctx, sql, movieID).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Synopsis,
		&movie.Duration,
		&movie.Director,
		&movie.ReleaseDate,
		&movie.PosterPath,
		&movie.BackdropPath,
		&genresJSON,
		&castListJSON,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Println("[ERROR]: Failed to get movie details: ", err.Error())
		return nil, err
	}

	// Unmarshal the JSON byte slices into the Go structs
	if len(genresJSON) > 0 {
		if err := json.Unmarshal(genresJSON, &movie.Genre); err != nil {
			log.Println("[ERROR]: Failed to unmarshal genres JSON: ", err.Error())
			return nil, err
		}
	}

	if len(castListJSON) > 0 {
		if err := json.Unmarshal(castListJSON, &movie.Cast); err != nil {
			log.Println("[ERROR]: Failed to unmarshal casts JSON: ", err.Error())
			return nil, err
		}
	}

	return &movie, nil
}

// func (m *MovieRepository) GetAdminMovieDetail(ctx context.Context, movieID int) (*models.MovieDetailStruct, error) {
// 	sql := `
//         SELECT
//             mv.id,
//             mv.title,
//             mv.synopsis,
//             mv.duration,
//             d.name AS director_name,
//             mv.release_date,
//             mv.poster_path,
//             mv.backdrop_path,
//             COALESCE(json_agg(DISTINCT g.name) FILTER (WHERE g.id IS NOT NULL), '[]') AS genres,
//             COALESCE(json_agg(DISTINCT c.name) FILTER (WHERE c.id IS NOT NULL), '[]') AS casts
//         FROM
//             movie mv
//         JOIN
//             director d ON mv.director_id = d.id
//         LEFT JOIN
//             movies_genre mg ON mv.id = mg.movie_id
//         LEFT JOIN
//             genre g ON mg.genre_id = g.id
//         LEFT JOIN
//             movie_cast mc ON mv.id = mc.movie_id
//         LEFT JOIN
//             "cast" c ON mc.cast_id = c.id
//         WHERE
//             mv.id = $1
//         GROUP BY
//             mv.id, d.name;
//     `

// 	var movie models.MovieDetailStruct
// 	var genresJSON []byte
// 	var castsJSON []byte

// 	err := m.db.QueryRow(ctx, sql, movieID).Scan(
// 		&movie.ID,
// 		&movie.Title,
// 		&movie.Synopsis,
// 		&movie.Duration,
// 		&movie.Director,
// 		&movie.ReleaseDate,
// 		&movie.PosterPath,
// 		&movie.BackdropPath,
// 		&genresJSON,
// 		&castsJSON,
// 	)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	if err := json.Unmarshal(genresJSON, &movie.Genre); err != nil {
// 		return nil, err
// 	}
// 	if err := json.Unmarshal(castsJSON, &movie.Cast); err != nil {
// 		return nil, err
// 	}

// 	return &movie, nil
// }

// percobaan kedua
// func (m *MovieRepository) GetMoviesWithPagination(
// 	ctx context.Context,
// 	limit, page int,
// 	title, genreName string,
// ) ([]models.MovieStruct, int, error) {
// 	offset := (page - 1) * limit

// 	var redisKey string
// 	useCache := false

// 	// ✅ Cache page pertama all movies, tanpa filter
// 	if page == 1 && title == "" && genreName == "" {
// 		redisKey = "tickitz:movies-all-first-page"
// 		useCache = true
// 	} else {
// 		// key unik untuk kombinasi filter + pagination
// 		redisKey = fmt.Sprintf("movies:limit=%d:page=%d:title=%s:genre=%s", limit, page, title, genreName)
// 	}

// 	// 🔍 cek Redis
// 	cache, err := m.rdb.Get(ctx, redisKey).Result()
// 	if err == nil {
// 		var movies []models.MovieStruct
// 		if err := json.Unmarshal([]byte(cache), &movies); err == nil {
// 			// hitung total
// 			var total int
// 			countQuery := `
// 				SELECT COUNT(DISTINCT m.id)
// 				FROM movie m
// 				LEFT JOIN movies_genre mg2 ON m.id = mg2.movie_id
// 				LEFT JOIN genre g2 ON mg2.genre_id = g2.id
// 				WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
// 				AND ($2 = '' OR EXISTS (
// 					SELECT 1
// 					FROM movies_genre mgx
// 					JOIN genre gx ON mgx.genre_id = gx.id
// 					WHERE mgx.movie_id = m.id
// 					AND gx.name ILIKE $2
// 				))
// 			`
// 			if err := m.db.QueryRow(ctx, countQuery, title, genreName).Scan(&total); err != nil {
// 				return nil, 0, err
// 			}
// 			return movies, total, nil
// 		}
// 	}

// 	// 📥 Query DB
// 	query := `
// 	SELECT
//     m.id,
//     m.director_id,
//     d.name AS director_name,
//     m.title,
//     m.synopsis,
//     m.popularity,
//     COALESCE(m.poster_path, '') AS poster_path,
//     COALESCE(m.backdrop_path, '') AS backdrop_path,
//     COALESCE(m.duration, 0) AS duration,
//     m.release_date,
//     COALESCE(string_agg(DISTINCT g.name, ','), '') AS genres,
//     COALESCE(string_agg(DISTINCT c.name, ','), '') AS casts
// FROM movie m
// LEFT JOIN director d ON m.director_id = d.id
// LEFT JOIN movies_genre mg ON m.id = mg.movie_id
// LEFT JOIN genre g ON mg.genre_id = g.id
// LEFT JOIN movie_cast mc ON m.id = mc.movie_id
// LEFT JOIN "cast" c ON mc.cast_id = c.id
// WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
// AND ($2 = '' OR EXISTS (
//     SELECT 1
//     FROM movies_genre mgx
//     JOIN genre gx ON mgx.genre_id = gx.id
//     WHERE mgx.movie_id = m.id
//     AND gx.name ILIKE $2
// ))
// GROUP BY m.id, d.name
// ORDER BY m.release_date DESC
// LIMIT $3 OFFSET $4

// 	`

// 	rows, err := m.db.Query(ctx, query, title, genreName, limit, offset)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	defer rows.Close()

// 	var movies []models.MovieStruct
// 	for rows.Next() {
// 		var movie models.MovieStruct
// 		var genresStr, castsStr string //  tambahan variabel string

// 		if err := rows.Scan(
// 			&movie.Id,
// 			&movie.DirectorId,
// 			&movie.DirectorName,
// 			&movie.Title,
// 			&movie.Synopsis,
// 			&movie.Popularity,
// 			&movie.PosterPath,
// 			&movie.BackdropPath,
// 			&movie.Duration,
// 			&movie.ReleaseDate,
// 			&genresStr, //  scan ke string
// 			&castsStr,  //  scan ke string
// 		); err != nil {
// 			return nil, 0, err
// 		}

// 		//  split manual biar jadi []string
// 		if genresStr != "" {
// 			movie.Genres = strings.Split(genresStr, ",")
// 		} else {
// 			movie.Genres = []string{}
// 		}

// 		if castsStr != "" {
// 			movie.Casts = strings.Split(castsStr, ",")
// 		} else {
// 			movie.Casts = []string{}
// 		}

// 		movies = append(movies, movie)
// 	}

// 	// 🔢 hitung total
// 	countQuery := `
// 		SELECT COUNT(DISTINCT m.id)
// 		FROM movie m
// 		LEFT JOIN movies_genre mg2 ON m.id = mg2.movie_id
// 		LEFT JOIN genre g2 ON mg2.genre_id = g2.id
// 		WHERE ($1 = '' OR m.title ILIKE '%' || $1 || '%')
// 		AND ($2 = '' OR EXISTS (
// 			SELECT 1
// 			FROM movies_genre mgx
// 			JOIN genre gx ON mgx.genre_id = gx.id
// 			WHERE mgx.movie_id = m.id
// 			AND gx.name ILIKE $2
// 		))
// 	`
// 	var total int
// 	if err := m.db.QueryRow(ctx, countQuery, title, genreName).Scan(&total); err != nil {
// 		return nil, 0, err
// 	}

// 	//  Simpan ke Redis
// 	if useCache || true { // simpan semua query, bisa diganti false kalau mau hanya page 1
// 		res, _ := json.Marshal(movies)
// 		_ = m.rdb.Set(ctx, redisKey, string(res), 20*time.Minute).Err()
// 	}

// 	return movies, total, nil
// }

func (m *MovieRepository) GetMoviesWithPagination(
	ctx context.Context,
	limit, page int,
	title string,
	genres []string, //  terima array genre
) ([]models.MovieStruct, int, error) {
	offset := (page - 1) * limit

	// --- Query utama ---
	baseQuery := `
		SELECT 
			m.id,
			m.director_id,
			d.name AS director_name,
			m.title,
			m.synopsis,
			m.popularity,
			COALESCE(m.poster_path, '') AS poster_path,
			COALESCE(m.backdrop_path, '') AS backdrop_path,
			COALESCE(m.duration, 0) AS duration,
			m.release_date,
			COALESCE(string_agg(DISTINCT g.name, ','), '') AS genres,
			COALESCE(string_agg(DISTINCT c.name, ','), '') AS casts
		FROM movie m
		LEFT JOIN director d ON m.director_id = d.id
		LEFT JOIN movies_genre mg ON m.id = mg.movie_id
		LEFT JOIN genre g ON mg.genre_id = g.id
		LEFT JOIN movie_cast mc ON m.id = mc.movie_id
		LEFT JOIN "cast" c ON mc.cast_id = c.id
		WHERE 1=1
	
	`
	// 	AND m.archived_at IS NULL
	args := []interface{}{}
	argID := 1

	// --- Filter judul ---
	if title != "" {
		baseQuery += fmt.Sprintf(" AND m.title ILIKE $%d", argID)
		args = append(args, "%"+title+"%")
		argID++
	}

	// --- Filter multi-genre (AND) ---
	for _, gname := range genres {
		baseQuery += fmt.Sprintf(`
			AND EXISTS (
				SELECT 1 FROM movies_genre mgx
				JOIN genre gx ON mgx.genre_id = gx.id
				WHERE mgx.movie_id = m.id AND gx.name ILIKE $%d
			)
		`, argID)
		args = append(args, gname)
		argID++
	}

	baseQuery += `
		GROUP BY m.id, d.name
		ORDER BY m.release_date DESC
	`

	// --- Hitung total ---
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") AS sub"
	var total int
	if err := m.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// --- Tambah LIMIT + OFFSET ---
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := m.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.MovieStruct
	for rows.Next() {
		var movie models.MovieStruct
		var genresStr, castsStr string

		if err := rows.Scan(
			&movie.Id,
			&movie.DirectorId,
			&movie.DirectorName,
			&movie.Title,
			&movie.Synopsis,
			&movie.Popularity,
			&movie.PosterPath,
			&movie.BackdropPath,
			&movie.Duration,
			&movie.ReleaseDate,
			&genresStr,
			&castsStr,
		); err != nil {
			return nil, 0, err
		}
		if genresStr != "" {
			movie.Genres = strings.Split(genresStr, ",")
		} else {
			movie.Genres = []string{} //  fallback kosong
		}
		if castsStr != "" {
			movie.Casts = strings.Split(castsStr, ",")
		}

		movies = append(movies, movie)
	}

	return movies, total, nil
}

// GetAllMovies retrieves all movies from the database
func (m *MovieRepository) GetAllMovies(ctx context.Context) ([]models.MovieListAdminStruct, error) {
	sql := `
		SELECT
			mv.id,
			mv.title,
			mv.synopsis,
			mv.poster_path,
			COALESCE(mv.duration, 0) AS duration,
			mv.release_date,
			COALESCE(
				json_agg(DISTINCT g.name) FILTER (WHERE g.id IS NOT NULL),
				'[]'
			) AS genres,
			 mv.archived_at 
		FROM
			movie mv
		LEFT JOIN
			movies_genre mg ON mv.id = mg.movie_id
		LEFT JOIN
			genre g ON mg.genre_id = g.id
		GROUP BY
			mv.id
		ORDER BY
			mv.created_at DESC;
	`
	// 	WHERE mv.archived_at IS NULL
	rows, err := m.db.Query(ctx, sql)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieListAdminStruct
	for rows.Next() {
		var movie models.MovieListAdminStruct
		var genresJSON []byte

		if err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.Synopsis,
			&movie.PosterPath,
			&movie.Duration,
			&movie.ReleaseDate,
			&genresJSON,
			&movie.ArchivedAt,
		); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return nil, err
		}

		// Unmarshal genres langsung ke []string
		if len(genresJSON) > 0 {
			if err := json.Unmarshal(genresJSON, &movie.Genres); err != nil {
				log.Println("[ERROR]: Failed to unmarshal genres JSON: ", err.Error())
				return nil, err
			}
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MovieRepository) ArchiveMovie(ctx context.Context, movieId int) (models.MovieArchived, error) {
	sql := `
        UPDATE movie
        SET archived_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND archived_at IS NULL
        RETURNING id, title, archived_at
    `
	var archivedMovie models.MovieArchived
	err := m.db.QueryRow(ctx, sql, movieId).Scan(
		&archivedMovie.Id,
		&archivedMovie.Title,
		&archivedMovie.Archived_at,
	)
	if err != nil {
		return models.MovieArchived{}, err
	}

	return archivedMovie, nil
}

// INI YANG AKAN DIPAKAI (FIX)
// redist
func (m *MovieRepository) GetPopularMovies(ctx context.Context) ([]models.MovieStruct, error) {
	redisKey := "moviePopular"

	// Cek Redis dulu
	cache, err := m.rdb.Get(ctx, redisKey).Result()
	if err == nil {
		var movies []models.MovieStruct
		if err := json.Unmarshal([]byte(cache), &movies); err == nil && len(movies) > 0 {
			return movies, nil
		}
	}

	// Query ke database berdasarkan popularitas
	query := `
			SELECT
				m.id,
				m.director_id,
				m.title,
				m.synopsis,
				m.popularity,
				COALESCE(m.poster_path, '') AS poster_path,
				COALESCE(m.backdrop_path, '') AS backdrop_path,
				COALESCE(m.duration, 0) AS duration,
				m.release_date,
				d.name AS director_name,
				COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.id IS NOT NULL), '{}') AS genres,
				COALESCE(array_agg(DISTINCT c.name) FILTER (WHERE c.id IS NOT NULL), '{}') AS casts
			FROM movie m
			LEFT JOIN director d ON m.director_id = d.id
			LEFT JOIN movies_genre mg ON m.id = mg.movie_id
			LEFT JOIN genre g ON mg.genre_id = g.id
			LEFT JOIN movie_cast mc ON m.id = mc.movie_id
			LEFT JOIN "cast" c ON mc.cast_id = c.id
			
			GROUP BY m.id, d.name
			ORDER BY m.popularity DESC;
    `
	// WHERE m.archived_at IS NULL
	rows, err := m.db.Query(ctx, query)
	if err != nil {
		log.Println("[ERROR]:", err)
		return nil, err
	}
	defer rows.Close()

	var result []models.MovieStruct
	for rows.Next() {
		var movie models.MovieStruct
		var genres, casts []string

		if err := rows.Scan(
			&movie.Id,
			&movie.DirectorId,
			&movie.Title,
			&movie.Synopsis,
			&movie.Popularity,
			&movie.PosterPath,
			&movie.BackdropPath,
			&movie.Duration,
			&movie.ReleaseDate,
			&movie.DirectorName,
			&genres,
			&casts,
		); err != nil {
			log.Println("[ERROR]:", err)
			return nil, err
		}

		movie.Genres = genres
		movie.Casts = casts

		result = append(result, movie)
	}

	// Simpan ke Redis
	res, err := json.Marshal(result)
	if err != nil {
		log.Println("[ERROR]: ", err.Error())
	}
	if err := m.rdb.Set(ctx, redisKey, string(res), time.Minute*20).Err(); err != nil {
		log.Println("[DEBUG] redis set", err.Error())
	}

	return result, nil
}

func (m *MovieRepository) GetUpcomingMovies(ctx context.Context) ([]models.MovieStruct, error) {
	redisKey := "movieUpcoming"

	// Cek Redis dulu
	cache, err := m.rdb.Get(ctx, redisKey).Result()
	if err == nil {
		var movies []models.MovieStruct
		if err := json.Unmarshal([]byte(cache), &movies); err == nil && len(movies) > 0 {
			return movies, nil
		}
	}

	// Query database
	query := `
SELECT
    m.id,
    m.director_id,
    d.name AS director_name,
    m.title,
    m.synopsis,
    m.popularity,
    COALESCE(m.poster_path, '') AS poster_path,
    COALESCE(m.backdrop_path, '') AS backdrop_path,
    COALESCE(m.duration, 0) AS duration,
    m.release_date,
    COALESCE(array_agg(DISTINCT g.name) FILTER (WHERE g.id IS NOT NULL), '{}') AS genres,
    COALESCE(array_agg(DISTINCT c.name) FILTER (WHERE c.id IS NOT NULL), '{}') AS casts
FROM movie m
LEFT JOIN director d ON m.director_id = d.id
LEFT JOIN movies_genre mg ON m.id = mg.movie_id
LEFT JOIN genre g ON mg.genre_id = g.id
LEFT JOIN movie_cast mc ON m.id = mc.movie_id
LEFT JOIN "cast" c ON mc.cast_id = c.id
WHERE m.release_date > CURRENT_DATE

GROUP BY m.id, d.name
ORDER BY m.release_date ASC;
`
	// AND m.archived_at IS NULL
	rows, err := m.db.Query(ctx, query)
	if err != nil {
		log.Println("[ERROR]:", err)
		return nil, err
	}
	defer rows.Close()

	var result []models.MovieStruct
	for rows.Next() {
		var movie models.MovieStruct
		var genres, casts []string

		// var archivedAt *time.Time
		var directorName string

		if err := rows.Scan(
			&movie.Id,
			&movie.DirectorId,
			&directorName,
			&movie.Title,
			&movie.Synopsis,
			&movie.Popularity,
			&movie.PosterPath,
			&movie.BackdropPath,
			&movie.Duration,
			&movie.ReleaseDate,
			&genres,
			&casts,
		); err != nil {
			return nil, err
		}

		movie.DirectorName = directorName
		movie.Genres = genres
		movie.Casts = casts
		// movie.ArchivedAt = archivedAt

		result = append(result, movie)
	}

	// Simpan hasil ke Redis
	res, err := json.Marshal(result)
	if err != nil {
		log.Println("[ERROR]: ", err.Error())
	}

	if err := m.rdb.Set(ctx, redisKey, string(res), time.Minute*20).Err(); err != nil {
		log.Println("[DEBUG] redis set", err.Error())
	}

	return result, nil
}

// func (m *MovieRepository) UpdateMovie(ctx context.Context, movieID int, updatedMovie models.MovieUpdateRequest, posterPathString, backdropPathString string) error {
// 	tx, err := m.db.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("gagal memulai transaksi: %w", err)
// 	}
// 	defer tx.Rollback(ctx)

// 	setParts := []string{}
// 	args := []interface{}{}
// 	argID := 1

// 	if updatedMovie.Title != nil {
// 		setParts = append(setParts, fmt.Sprintf("title = $%d", argID))
// 		args = append(args, *updatedMovie.Title)
// 		argID++
// 	}

// 	if updatedMovie.Synopsis != nil {
// 		setParts = append(setParts, fmt.Sprintf("synopsis = $%d", argID))
// 		args = append(args, *updatedMovie.Synopsis)
// 		argID++
// 	}

// 	if updatedMovie.Duration != nil {
// 		setParts = append(setParts, fmt.Sprintf("duration = $%d", argID))
// 		args = append(args, *updatedMovie.Duration)
// 		argID++
// 	}

// 	if updatedMovie.DirectorId != nil {
// 		setParts = append(setParts, fmt.Sprintf("director_id = $%d", argID))
// 		args = append(args, *updatedMovie.DirectorId)
// 		argID++
// 	}

// 	if updatedMovie.ReleaseDate != nil {
// 		setParts = append(setParts, fmt.Sprintf("release_date = $%d", argID))
// 		args = append(args, *updatedMovie.ReleaseDate)
// 		argID++
// 	}

// 	if posterPathString != "" {
// 		setParts = append(setParts, fmt.Sprintf("poster_path = $%d", argID))
// 		args = append(args, posterPathString)
// 		argID++
// 	}

// 	if backdropPathString != "" {
// 		setParts = append(setParts, fmt.Sprintf("backdrop_path = $%d", argID))
// 		args = append(args, backdropPathString)
// 		argID++
// 	}

// 	if updatedMovie.Popularity != nil {
// 		setParts = append(setParts, fmt.Sprintf("popularity = $%d", argID))
// 		args = append(args, *updatedMovie.Popularity)
// 		argID++
// 	}

// 	// kalau tidak ada field yang diupdate
// 	if len(setParts) == 0 {
// 		return errors.New("tidak ada data yang diupdate")
// 	}

// 	// tambahkan updated_at
// 	setParts = append(setParts, "update_at = NOW()")

// 	query := fmt.Sprintf(`UPDATE movie SET %s WHERE id = $%d`, strings.Join(setParts, ", "), argID)
// 	args = append(args, movieID)

// 	_, err = tx.Exec(ctx, query, args...)
// 	if err != nil {
// 		return fmt.Errorf("gagal update movie: %w", err)
// 	}

// 	// update genre jika dikirim
// 	if updatedMovie.GenreIDs != nil {
// 		_, err := tx.Exec(ctx, "DELETE FROM movies_genre WHERE movie_id = $1", movieID)
// 		if err != nil {
// 			return err
// 		}
// 		for _, genre := range updatedMovie.GenreIDs {
// 			_, err := tx.Exec(ctx, "INSERT INTO movies_genre (movie_id, genre_id) VALUES ($1, $2)", movieID, genre.ID)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	// update cast jika dikirim
// 	if updatedMovie.CastIDs != nil {
// 		_, err := tx.Exec(ctx, "DELETE FROM movie_cast WHERE movie_id = $1", movieID)
// 		if err != nil {
// 			return err
// 		}
// 		for _, cast := range updatedMovie.CastIDs {
// 			_, err := tx.Exec(ctx, "INSERT INTO movie_cast (movie_id, cast_id) VALUES ($1, $2)", movieID, cast.ID)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return tx.Commit(ctx)
// }

func (m *MovieRepository) UpdateMovie(
	ctx context.Context,
	movieID int,
	updatedMovie models.MovieUpdateRequest,
	posterPathString, backdropPathString string,
) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer tx.Rollback(ctx)

	setParts := []string{}
	args := []interface{}{}
	argID := 1

	// === Title ===
	if updatedMovie.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argID))
		args = append(args, *updatedMovie.Title)
		argID++
	}

	// === Synopsis ===
	if updatedMovie.Synopsis != nil {
		setParts = append(setParts, fmt.Sprintf("synopsis = $%d", argID))
		args = append(args, *updatedMovie.Synopsis)
		argID++
	}

	// === Duration ===
	if updatedMovie.Duration != nil {
		setParts = append(setParts, fmt.Sprintf("duration = $%d", argID))
		args = append(args, *updatedMovie.Duration)
		argID++
	}

	// === Director (pakai name → id) ===
	if updatedMovie.Director != nil {
		var directorID int
		err := tx.QueryRow(ctx,
			`INSERT INTO director (name)
			 VALUES ($1)
			 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			 RETURNING id`, *updatedMovie.Director,
		).Scan(&directorID)
		if err != nil {
			return fmt.Errorf("gagal insert/get director: %w", err)
		}

		setParts = append(setParts, fmt.Sprintf("director_id = $%d", argID))
		args = append(args, directorID)
		argID++
	}

	// === Release Date ===
	if updatedMovie.ReleaseDate != nil {
		setParts = append(setParts, fmt.Sprintf("release_date = $%d", argID))
		args = append(args, *updatedMovie.ReleaseDate)
		argID++
	}

	// === Poster ===
	if posterPathString != "" {
		setParts = append(setParts, fmt.Sprintf("poster_path = $%d", argID))
		args = append(args, posterPathString)
		argID++
	}

	// === Backdrop ===
	if backdropPathString != "" {
		setParts = append(setParts, fmt.Sprintf("backdrop_path = $%d", argID))
		args = append(args, backdropPathString)
		argID++
	}

	// === Popularity ===
	if updatedMovie.Popularity != nil {
		setParts = append(setParts, fmt.Sprintf("popularity = $%d", argID))
		args = append(args, *updatedMovie.Popularity)
		argID++
	}

	// kalau tidak ada kolom lain yang diupdate, tetap update updated_at
	if len(setParts) == 0 {
		query := `UPDATE movie SET update_at = NOW() WHERE id = $1`
		args = append(args, movieID)
		_, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("gagal update movie: %w", err)
		}
	} else {
		query := fmt.Sprintf(
			`UPDATE movie SET %s, update_at = NOW() WHERE id = $%d`,
			strings.Join(setParts, ", "),
			argID,
		)
		args = append(args, movieID)
		_, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("gagal update movie: %w", err)
		}
	}

	// === Update Genres ===
	if updatedMovie.Genres != nil {
		_, err := tx.Exec(ctx, "DELETE FROM movies_genre WHERE movie_id = $1", movieID)
		if err != nil {
			return fmt.Errorf("gagal hapus genres lama: %w", err)
		}
		for _, g := range updatedMovie.Genres {
			var genreID int
			err := tx.QueryRow(ctx,
				`INSERT INTO genre (name)
				 VALUES ($1)
				 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
				 RETURNING id`, g).Scan(&genreID)
			if err != nil {
				return fmt.Errorf("gagal insert/get genre: %w", err)
			}
			_, err = tx.Exec(ctx,
				`INSERT INTO movies_genre (movie_id, genre_id) VALUES ($1, $2)
				 ON CONFLICT DO NOTHING`, movieID, genreID)
			if err != nil {
				return fmt.Errorf("gagal insert relasi genre: %w", err)
			}
		}
	}

	// === Update Casts ===
	if updatedMovie.Casts != nil {
		_, err := tx.Exec(ctx, "DELETE FROM movie_cast WHERE movie_id = $1", movieID)
		if err != nil {
			return fmt.Errorf("gagal hapus casts lama: %w", err)
		}
		for _, c := range updatedMovie.Casts {
			var castID int
			err := tx.QueryRow(ctx,
				`INSERT INTO "cast" (name)
				 VALUES ($1)
				 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
				 RETURNING id`, c).Scan(&castID)
			if err != nil {
				return fmt.Errorf("gagal insert/get cast: %w", err)
			}
			_, err = tx.Exec(ctx,
				`INSERT INTO movie_cast (movie_id, cast_id) VALUES ($1, $2)
				 ON CONFLICT DO NOTHING`, movieID, castID)
			if err != nil {
				return fmt.Errorf("gagal insert relasi cast: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}

// mencoba new movie
// func (m *MovieRepository) AddNewMovie(ctx context.Context, movie models.NewMovieRequest, posterPath, backdropPath string, genreIDs, castIDs []int) error {
// 	tx, err := m.db.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("gagal memulai transaksi: %w", err)
// 	}
// 	defer func() {
// 		_ = tx.Rollback(ctx) // rollback jika belum commit
// 	}()

// 	var movieID int
// 	queryMovie := `
// 		INSERT INTO movie (
// 			title, synopsis, duration, director_id, release_date, poster_path, backdrop_path, popularity, created_at, update_at
// 		)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
// 		RETURNING id`

// 	err = tx.QueryRow(ctx, queryMovie,
// 		movie.Title,
// 		movie.Synopsis,
// 		movie.Duration,
// 		movie.DirectorId,
// 		movie.ReleaseDate,
// 		posterPath,
// 		backdropPath,
// 		movie.Popularity,
// 	).Scan(&movieID)
// 	if err != nil {
// 		log.Println("[ERROR]: Gagal insert movie:", err.Error())
// 		return err
// 	}

// 	for _, genreID := range genreIDs {
// 		_, err := tx.Exec(ctx, "INSERT INTO movies_genre (movie_id, genre_id) VALUES ($1, $2)", movieID, genreID)
// 		if err != nil {
// 			log.Println("[ERROR]: Gagal insert genre:", err.Error())
// 			return err
// 		}
// 	}

// 	for _, castID := range castIDs {
// 		_, err := tx.Exec(ctx, "INSERT INTO movie_cast (movie_id, cast_id) VALUES ($1, $2)", movieID, castID)
// 		if err != nil {
// 			log.Println("[ERROR]: Gagal insert cast:", err.Error())
// 			return err
// 		}
// 	}

// 	return tx.Commit(ctx)
// }

// func (m *MovieRepository) getOrCreateGenre(ctx context.Context, tx pgx.Tx, name string) (int, error) {
// 	var id int
// 	err := tx.QueryRow(ctx, `SELECT id FROM genre WHERE name = $1`, name).Scan(&id)
// 	if err != nil {
// 		// kalau belum ada, insert baru
// 		err = tx.QueryRow(ctx, `INSERT INTO genre (name) VALUES ($1) RETURNING id`, name).Scan(&id)
// 		if err != nil {
// 			return 0, err
// 		}
// 	}
// 	return id, nil
// }

// func (m *MovieRepository) getOrCreateCast(ctx context.Context, tx pgx.Tx, name string) (int, error) {
// 	var id int
// 	err := tx.QueryRow(ctx, `SELECT id FROM cast WHERE name = $1`, name).Scan(&id)
// 	if err != nil {
// 		// kalau belum ada, insert baru
// 		err = tx.QueryRow(ctx, `INSERT INTO cast (name) VALUES ($1) RETURNING id`, name).Scan(&id)
// 		if err != nil {
// 			return 0, err
// 		}
// 	}
// 	return id, nil
// }

// func (m *MovieRepository) AddNewMovie(
// 	ctx context.Context,
// 	movie models.NewMovieRequest,
// 	posterPath, backdropPath string,
// ) error {
// 	tx, err := m.db.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("gagal memulai transaksi: %w", err)
// 	}
// 	defer func() { _ = tx.Rollback(ctx) }()

// 	// === Insert Movie ===
// 	var movieID int
// 	queryMovie := `
// 		INSERT INTO movie (
// 			title, synopsis, duration, director_id, release_date,
// 			poster_path, backdrop_path, popularity, created_at, update_at
// 		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())
// 		RETURNING id`
// 	err = tx.QueryRow(ctx, queryMovie,
// 		movie.Title, movie.Synopsis, movie.Duration, movie.DirectorId,
// 		movie.ReleaseDate, posterPath, backdropPath, movie.Popularity,
// 	).Scan(&movieID)
// 	if err != nil {
// 		return fmt.Errorf("gagal insert movie: %w", err)
// 	}

// 	// === Handle Genres ===
// 	for _, g := range movie.Genres {
// 		var genreID int
// 		err := tx.QueryRow(ctx,
// 			`INSERT INTO genre (name)
//          VALUES ($1)
//          ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
//          RETURNING id`, g).Scan(&genreID)
// 		if err != nil {
// 			return fmt.Errorf("gagal insert/get genre: %w", err)
// 		}

// 		_, err = tx.Exec(ctx,
// 			`INSERT INTO movies_genre (movie_id, genre_id)
//          VALUES ($1, $2)
//          ON CONFLICT DO NOTHING`, // <== supaya tidak error kalau sudah ada
// 			movieID, genreID,
// 		)
// 		if err != nil {
// 			return fmt.Errorf("gagal insert relasi genre: %w", err)
// 		}
// 	}

// 	// === Handle Casts ===
// 	for _, c := range movie.Casts {
// 		var castID int
// 		err := tx.QueryRow(ctx,
// 			`INSERT INTO "cast" (name)
//          VALUES ($1)
//          ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
//          RETURNING id`, c).Scan(&castID)
// 		if err != nil {
// 			return fmt.Errorf("gagal insert/get cast: %w", err)
// 		}

// 		_, err = tx.Exec(ctx,
// 			`INSERT INTO movie_cast (movie_id, cast_id)
//          VALUES ($1, $2)
//          ON CONFLICT DO NOTHING`, // <== sama, biar aman dari duplicate
// 			movieID, castID,
// 		)
// 		if err != nil {
// 			return fmt.Errorf("gagal insert relasi cast: %w", err)
// 		}
// 	}

// 	return tx.Commit(ctx)
// }

func (m *MovieRepository) AddNewMovie(
	ctx context.Context,
	movie models.NewMovieRequest,
	posterPath, backdropPath string,
) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// === Handle Director ===
	var directorID int
	err = tx.QueryRow(ctx,
		`INSERT INTO director (name)
         VALUES ($1)
         ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
         RETURNING id`, movie.DirectorName,
	).Scan(&directorID)
	if err != nil {
		return fmt.Errorf("gagal insert/get director: %w", err)
	}

	// === Insert Movie ===
	var movieID int
	queryMovie := `
		INSERT INTO movie (
			title, synopsis, duration, director_id, release_date,
			poster_path, backdrop_path, popularity, created_at, update_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())
		RETURNING id`
	err = tx.QueryRow(ctx, queryMovie,
		movie.Title, movie.Synopsis, movie.Duration, directorID,
		movie.ReleaseDate, posterPath, backdropPath, movie.Popularity,
	).Scan(&movieID)
	if err != nil {
		return fmt.Errorf("gagal insert movie: %w", err)
	}

	// === Handle Genres ===
	for _, g := range movie.Genres {
		var genreID int
		err := tx.QueryRow(ctx,
			`INSERT INTO genre (name)
			 VALUES ($1)
			 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			 RETURNING id`, g).Scan(&genreID)
		if err != nil {
			return fmt.Errorf("gagal insert/get genre: %w", err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO movies_genre (movie_id, genre_id) 
             VALUES ($1, $2)
             ON CONFLICT DO NOTHING`,
			movieID, genreID,
		)
		if err != nil {
			return fmt.Errorf("gagal insert relasi genre: %w", err)
		}
	}

	// === Handle Casts ===
	for _, c := range movie.Casts {
		var castID int
		err := tx.QueryRow(ctx,
			`INSERT INTO "cast" (name)
			 VALUES ($1)
			 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			 RETURNING id`, c).Scan(&castID)
		if err != nil {
			return fmt.Errorf("gagal insert/get cast: %w", err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO movie_cast (movie_id, cast_id) 
             VALUES ($1, $2)
             ON CONFLICT DO NOTHING`,
			movieID, castID,
		)
		if err != nil {
			return fmt.Errorf("gagal insert relasi cast: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// ININ HANG PAKAI JSON {ID: 1, GENRE:"ACTION"}
// func (m *MovieRepository) GetPopularMovies(ctx context.Context) ([]models.MovieStruct, error) {
// 	// Cek Redis terlebih dahulu, jika ada nilainya, maka gunakan nilai dari Redis
// 	redisKey := "popularMovies"
// 	cache, err := m.rdb.Get(ctx, redisKey).Result()

// 	if err != nil {
// 		log.Println("[ERROR]: ", err.Error())
// 		if err == redis.Nil {
// 			log.Printf("\nkey %s does not exist\n", redisKey)
// 		} else {
// 			log.Println("Redis not working")
// 		}
// 	} else {
// 		var movies []models.MovieStruct
// 		log.Println(cache)
// 		if err := json.Unmarshal([]byte(cache), &movies); err != nil {
// 			log.Println("[ERROR]: ", err.Error())
// 			return []models.MovieStruct{}, err
// 		}
// 		if len(movies) > 0 {
// 			return movies, nil
// 		}
// 	}

// 	// Mengambil data film populer dari database
// 	query := `
//         SELECT
//             m.id, m.director_id, m.title, m.synopsis, m.popularity,
//             COALESCE(m.backdrop_path, '') AS backdrop_path,
//             m.poster_path, COALESCE(m.duration, 0) AS duration, m.release_date, m.created_at, m.update_at,
//             d.name AS director_name,
//             COALESCE(json_agg(DISTINCT jsonb_build_object('id', g.id, 'name', g.name)) FILTER (WHERE g.id IS NOT NULL), '[]') AS genres,
//             COALESCE(json_agg(DISTINCT jsonb_build_object('id', c.id, 'name', c.name)) FILTER (WHERE c.id IS NOT NULL), '[]') AS casts
//         FROM
//             movie m
//         JOIN
//             director d ON m.director_id = d.id
//         LEFT JOIN
//             movies_genre mg ON m.id = mg.movie_id
//         LEFT JOIN
//             genre g ON mg.genre_id = g.id
//         LEFT JOIN
//             movie_cast mc ON m.id = mc.movie_id
//         LEFT JOIN
//             "cast" c ON mc.cast_id = c.id
//         GROUP BY
//             m.id, d.name
//         ORDER BY m.popularity DESC
//     `

// 	rows, err := m.db.Query(ctx, query)
// 	if err != nil {
// 		log.Println("[ERROR]: ", err.Error())
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var result []models.MovieStruct
// 	for rows.Next() {
// 		var movie models.MovieStruct
// 		var genresJSON, castsJSON []byte

// 		if err := rows.Scan(
// 			&movie.Id, &movie.DirectorId, &movie.Title, &movie.Synopsis, &movie.Popularity,
// 			&movie.BackdropPath, &movie.PosterPath, &movie.Duration, &movie.ReleaseDate, &movie.CreatedAt, &movie.UpdateAt,
// 			&movie.DirectorName, &genresJSON, &castsJSON,
// 		); err != nil {
// 			log.Println("[ERROR]: ", err.Error())
// 			return nil, err
// 		}

// 		// Unmarshal the JSON byte slices into the Go structs
// 		if err := json.Unmarshal(genresJSON, &movie.Genres); err != nil {
// 			log.Println("[ERROR]: Failed to unmarshal genres JSON", err.Error())
// 			return nil, err
// 		}
// 		if err := json.Unmarshal(castsJSON, &movie.Casts); err != nil {
// 			log.Println("[ERROR]: Failed to unmarshal casts JSON", err.Error())
// 			return nil, err
// 		}

// 		result = append(result, movie)
// 	}

// 	// Simpan hasil ke Redis
// 	res, err := json.Marshal(result)
// 	if err != nil {
// 		log.Println("[ERROR]: ", err.Error())
// 	}

// 	if err := m.rdb.Set(ctx, redisKey, string(res), time.Minute*20).Err(); err != nil {
// 		log.Println("[DEBUG] redis set", err.Error())
// 	}

// 	return result, nil
// }

// inin sama juga
// func (m *MovieRepository) GetUpcomingMovies(ctx context.Context) ([]models.MovieStruct, error) {
// 	// Cek redis terlebih dahulu. Jika nilainya ada , maka gunakan nilai dari redis
// 	// Buat kata kunci untuk di Redis
// 	redisKey := "movieUpcoming"

// 	// Ambil di penyimpanan di redis dengan kunci yang sudah dibuat
// 	cache, err := m.rdb.Get(ctx, redisKey).Result()

// 	// Error handling jika terjadi error : Kunci redis tidak ditemukan, atau bernilai nil, atau kunci redis tidak bekerja
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		if err == redis.Nil {
// 			log.Printf("\nkey %s does not exist\n", redisKey)
// 		} else {
// 			log.Println("Redis not working")
// 		}
// 	} else {
// 		// Lakukam parsing data dari redis menjadi bentuk JSON
// 		var movies []models.MovieStruct
// 		if err := json.Unmarshal([]byte(cache), &movies); err != nil {
// 			log.Println("[ERROR] :", err.Error())
// 			return []models.MovieStruct{}, nil
// 		}

// 		// Jika berhasil mendapatkan data di redis maka tampilkan hasil dari redis
// 		if len(movies) > 0 {
// 			return movies, nil
// 		}
// 	}

// 	sql := `
// 		SELECT
// 			m.id, m.title, m.poster_path, COALESCE(m.backdrop_path, '') AS backdrop_path, m.release_date,
// 			COALESCE(json_agg(DISTINCT jsonb_build_object('id', g.id, 'name', g.name)) FILTER (WHERE g.id IS NOT NULL), '[]') AS genres,
// 			COALESCE(json_agg(DISTINCT jsonb_build_object('id', c.id, 'name', c.name)) FILTER (WHERE c.id IS NOT NULL), '[]') AS cast_list
// 		FROM
// 			movie m
// 		LEFT JOIN
// 			movies_genre mg ON m.id = mg.movie_id
// 		LEFT JOIN
// 			genre g ON mg.genre_id = g.id
// 		LEFT JOIN
// 			movie_cast mc ON m.id = mc.movie_id
// 		LEFT JOIN
// 			"cast" c ON mc.cast_id = c.id
// 		WHERE
// 			m.release_date > CURRENT_DATE
// 		GROUP BY
// 			m.id
// 		ORDER BY
// 			m.release_date ASC;
// 	`

// 	rows, err := m.db.Query(ctx, sql)
// 	if err != nil {
// 		log.Println("[ERROR] :", err.Error())
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var result []models.MovieStruct
// 	for rows.Next() {
// 		var movie models.MovieStruct
// 		var genresJSON, castsJSON []byte

// 		// Sesuaikan jumlah variabel di sini agar sesuai dengan query SQL
// 		if err := rows.Scan(
// 			&movie.Id,
// 			&movie.Title,
// 			&movie.PosterPath,
// 			&movie.BackdropPath,
// 			&movie.ReleaseDate,
// 			&genresJSON,
// 			&castsJSON,
// 		); err != nil {
// 			log.Println("[ERROR]: ", err.Error())
// 			return nil, err
// 		}

// 		// Unmarshal the JSON byte slices into the Go structs
// 		if err := json.Unmarshal(genresJSON, &movie.Genres); err != nil {
// 			log.Println("[ERROR]: Failed to unmarshal genres JSON", err.Error())
// 			return nil, err
// 		}
// 		if err := json.Unmarshal(castsJSON, &movie.Casts); err != nil {
// 			log.Println("[ERROR]: Failed to unmarshal casts JSON", err.Error())
// 			return nil, err
// 		}

// 		result = append(result, movie)
// 	}

// 	// Simpan hasil ke Redis
// 	res, err := json.Marshal(result)
// 	if err != nil {
// 		log.Println("[ERROR]: ", err.Error())
// 	}

// 	if err := m.rdb.Set(ctx, redisKey, string(res), time.Minute*20).Err(); err != nil {
// 		log.Println("[DEBUG] redis set", err.Error())
// 	}

//		return result, nil
//	}
func (m *MovieRepository) GetAdminMovieDetail(ctx context.Context, movieID int) (*models.MovieDetailStruct, error) {
	sql := `
        SELECT
            mv.id,
            mv.title,
            mv.synopsis,
            mv.duration,
            d.name AS director_name,
            mv.release_date,
            mv.poster_path,
            mv.backdrop_path,
            COALESCE(
              json_agg(DISTINCT jsonb_build_object('id', g.id, 'name', g.name))
              FILTER (WHERE g.id IS NOT NULL), '[]'
            ) AS genres,
            COALESCE(
              json_agg(DISTINCT jsonb_build_object('id', c.id, 'name', c.name))
              FILTER (WHERE c.id IS NOT NULL), '[]'
            ) AS cast_list
        FROM
            movie mv
        JOIN
            director d ON mv.director_id = d.id
        LEFT JOIN
            movies_genre mg ON mv.id = mg.movie_id
        LEFT JOIN
            genre g ON mg.genre_id = g.id
        LEFT JOIN
            movie_cast mc ON mv.id = mc.movie_id
        LEFT JOIN
            "cast" c ON mc.cast_id = c.id
        WHERE
            mv.id = $1
        GROUP BY
            mv.id, d.name;
    `
	var movie models.MovieDetailStruct
	var genresJSON, castListJSON []byte

	err := m.db.QueryRow(ctx, sql, movieID).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Synopsis,
		&movie.Duration,
		&movie.Director,
		&movie.ReleaseDate,
		&movie.PosterPath,
		&movie.BackdropPath,
		&genresJSON,
		&castListJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Println("[ERROR]: Failed to get admin movie detail:", err.Error())
		return nil, err
	}

	if len(genresJSON) > 0 {
		if err := json.Unmarshal(genresJSON, &movie.Genre); err != nil {
			return nil, err
		}
	}
	if len(castListJSON) > 0 {
		if err := json.Unmarshal(castListJSON, &movie.Cast); err != nil {
			return nil, err
		}
	}

	return &movie, nil
}
