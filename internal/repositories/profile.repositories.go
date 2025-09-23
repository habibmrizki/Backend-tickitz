package repositories

// import (
// 	"context"

// 	"github.com/habibmrizki/back-end-tickitz/internal/models"
// 	"github.com/jackc/pgx/v5/pgconn/ctxwatch"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type ProfileRepository struct {
// 	db *pgxpool.Pool
// }

// func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
// 	return &NewMovieRepository{db: db}
// }

// func (p *ProfileRepository) EditImage(ctx context.Context) ([]models.Profile, error) {

// }

type ProfileUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Image string `json:"image"`
}
