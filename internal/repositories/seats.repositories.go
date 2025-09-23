package repositories

import (
	"context"
	"log"

	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatRepository struct {
	db *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{db: db}
}

// GetAvailableSeats mengambil daftar kursi yang tersedia untuk schedule_id tertentu.
func (s *SeatRepository) GetAvailableSeats(ctx context.Context, scheduleID int) ([]models.SeatStruct, error) {
	// Query untuk mengambil semua kursi yang belum dipesan
	query := `
        SELECT id, seats_code
        FROM seats
        WHERE id NOT IN (
            SELECT os.seats_id
            FROM order_seats os
            JOIN orders o ON o.id = os.orders_id
            WHERE o.schedule_id = $1
        );
    `

	rows, err := s.db.Query(ctx, query, scheduleID)
	if err != nil {
		log.Println("[ERROR]: Failed to query available seats:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var seats []models.SeatStruct
	for rows.Next() {
		var seat models.SeatStruct
		if err := rows.Scan(&seat.ID, &seat.SeatsCode); err != nil {
			log.Println("[ERROR]: Failed to scan seat row:", err.Error())
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, nil
}
