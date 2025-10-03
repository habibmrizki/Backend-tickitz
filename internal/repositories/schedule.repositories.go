package repositories

import (
	"context"
	"log"

	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ScheduleRepository adalah repositori untuk tabel 'schedule'
type ScheduleRepository struct {
	db *pgxpool.Pool
}

// NewScheduleRepository membuat instance baru dari ScheduleRepository
func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// GetAllSchedules mengambil semua jadwal dari database
// func (s *ScheduleRepository) GetAllSchedules(ctx context.Context) ([]models.ScheduleDetails, error) {
// 	query := `
// 		SELECT
// 			S.id,
// 			S.date,
// 			C.name AS cinema_name,
// 			L.location,
// 			T.time
// 		FROM schedule AS S
// 		JOIN cinema AS C ON S.cinema_id = C.id
// 		JOIN location AS L ON S.location_id = L.id
// 		JOIN time AS T ON S.time_id = T.id
// 	`

// 	rows, err := s.db.Query(ctx, query)
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var schedules []models.ScheduleDetails
// 	for rows.Next() {
// 		var schedule models.ScheduleDetails
// 		if err := rows.Scan(&schedule.ID, &schedule.Date, &schedule.CinemaName, &schedule.Location, &schedule.Time); err != nil {
// 			log.Println("[ERROR] : ", err.Error())
// 			return nil, err
// 		}
// 		schedules = append(schedules, schedule)
// 	}

// 	return schedules, nil
// }

// mencoba schedule
func (s *ScheduleRepository) GetAllSchedules(ctx context.Context) ([]models.ScheduleDetails, error) {
	query := `
		SELECT
			S.id,
			S.date,
			C.name AS cinema_name,
			C.image_path,
			L.location,
			T.time,
			S.movie_id
		FROM schedule AS S
		JOIN cinema AS C ON S.cinema_id = C.id
		JOIN location AS L ON S.location_id = L.id
		JOIN time AS T ON S.time_id = T.id
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return nil, err
	}
	defer rows.Close()

	var schedules []models.ScheduleDetails
	for rows.Next() {
		var schedule models.ScheduleDetails
		if err := rows.Scan(
			&schedule.ID,
			&schedule.Date,
			&schedule.CinemaName,
			&schedule.ImagePath,
			&schedule.Location,
			&schedule.Time,
			&schedule.MovieID,
		); err != nil {
			log.Println("[ERROR] : ", err.Error())
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (s *ScheduleRepository) GetSchedulesByMovieID(ctx context.Context, movieID int) ([]models.ScheduleDetails, error) {
	query := `
		SELECT
			S.id,
			S.date,
			C.name AS cinema_name,
			C.image_path,
			L.location,
			T.time,
			S.movie_id
		FROM schedule AS S
		JOIN cinema AS C ON S.cinema_id = C.id
		JOIN location AS L ON S.location_id = L.id
		JOIN time AS T ON S.time_id = T.id
		WHERE S.movie_id = $1
	`
	rows, err := s.db.Query(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.ScheduleDetails
	for rows.Next() {
		var schedule models.ScheduleDetails
		if err := rows.Scan(
			&schedule.ID,
			&schedule.Date,
			&schedule.CinemaName,
			&schedule.ImagePath,
			&schedule.Location,
			&schedule.Time,
			&schedule.MovieID,
		); err != nil {
			log.Println("[ERROR] Scan schedule failed: ", err.Error())
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// GetScheduleByMovieId mengambil jadwal berdasarkan ID film
// func (s *ScheduleRepository) GetScheduleByMovieId(ctx context.Context, movieID int) ([]models.ScheduleDetails, error) {
// 	query := `
// 		SELECT
// 			S.id,
// 			S.date,
// 			C.name AS cinema_name,
// 			L.location,
// 			T.time
// 		FROM schedule AS S
// 		JOIN cinema AS C ON S.cinema_id = C.id
// 		JOIN location AS L ON S.location_id = L.id
// 		JOIN time AS T ON S.time_id = T.id
// 		WHERE S.movie_id = $1
// 	`

// 	rows, err := s.db.Query(ctx, query, movieID)
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var schedules []models.ScheduleDetails
// 	for rows.Next() {
// 		var schedule models.ScheduleDetails
// 		if err := rows.Scan(&schedule.ID, &schedule.Date, &schedule.CinemaName, &schedule.Location, &schedule.Time); err != nil {
// 			log.Println("[ERROR] : ", err.Error())
// 			return nil, err
// 		}
// 		schedules = append(schedules, schedule)
// 	}

// 	if len(schedules) == 0 {
// 		return nil, pgx.ErrNoRows
// 	}

// 	return schedules, nil
// }
