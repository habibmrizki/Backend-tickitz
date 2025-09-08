// package repositories

// import (
// 	"context"
// 	"errors"
// 	"log"
// 	"time"

// 	"github.com/habibmrizki/back-end-tickitz/internal/models"
// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// // OrderRepository adalah repositori untuk tabel 'orders'
// type OrderRepository struct {
// 	db *pgxpool.Pool
// }

// // NewOrderRepository membuat instance baru dari OrderRepository
// func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
// 	return &OrderRepository{db: db}
// }

// // CreateOrder membuat order baru, seat yang dipilih, dan menambahkan poin ke profil pengguna
// func (o *OrderRepository) CreateOrder(ctx context.Context, orderData models.CreateOrderRequest) (models.Order, error) {
// 	// Mulai transaksi untuk memastikan semua operasi atomic
// 	tx, err := o.db.Begin(ctx)
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		return models.Order{}, err
// 	}
// 	defer tx.Rollback(ctx)

// 	// Persiapkan data untuk tabel 'orders'
// 	newOrder := models.Order{
// 		IDUsers:     orderData.IDUsers,
// 		IDSchedule:  orderData.IDSchedule,
// 		IDPayment:   orderData.IDPayment,
// 		FullName:    orderData.FullName,
// 		Email:       orderData.Email,
// 		PhoneNumber: orderData.PhoneNumber,
// 		TotalPrice:  orderData.TotalPrice,
// 		IsPaid:      false,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}

// 	// 1. Insert data ke tabel 'orders'
// 	queryOrders := `
// 		INSERT INTO orders (id_users, id_schedule, id_payment, full_name, email, phone_number, total_price, ispaid, created_at, update_at)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
// 		RETURNING id`

// 	err = tx.QueryRow(ctx, queryOrders,
// 		newOrder.IDUsers,
// 		newOrder.IDSchedule,
// 		newOrder.IDPayment,
// 		newOrder.FullName,
// 		newOrder.Email,
// 		newOrder.PhoneNumber,
// 		newOrder.TotalPrice,
// 		newOrder.IsPaid,
// 		newOrder.CreatedAt,
// 		newOrder.UpdatedAt).Scan(&newOrder.ID)
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		return models.Order{}, err
// 	}

// 	// 2. Ambil ID kursi dari kode kursi (seats_code)
// 	var seatsID []int
// 	for _, seatCode := range orderData.SeatsCode {
// 		var seatID int
// 		err := tx.QueryRow(ctx, "SELECT id FROM seats WHERE seats_code = $1", seatCode).Scan(&seatID)
// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				return models.Order{}, errors.New("seat with code " + seatCode + " not found")
// 			}
// 			log.Println("[ERROR] : ", err.Error())
// 			return models.Order{}, err
// 		}
// 		seatsID = append(seatsID, seatID)
// 	}

// 	// 3. Insert data ke tabel 'order_seats'
// 	queryOrderSeats := `
// 		INSERT INTO order_seats (orders_id, seats_id)
// 		VALUES ($1, $2)`

// 	for _, seatID := range seatsID {
// 		_, err := tx.Exec(ctx, queryOrderSeats, newOrder.ID, seatID)
// 		if err != nil {
// 			log.Println("[ERROR] : ", err.Error())
// 			return models.Order{}, err
// 		}
// 	}

// 	// 4. Update poin di tabel 'profile'
// 	queryUpdatePoints := `
// 		UPDATE profile
// 		SET point = point + 50
// 		WHERE users_id = $1`

// 	_, err = tx.Exec(ctx, queryUpdatePoints, newOrder.IDUsers)
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		return models.Order{}, err
// 	}

//		// Commit transaksi jika semua operasi berhasil
//		return newOrder, tx.Commit(ctx)
//	}
package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderRepository adalah repositori untuk tabel 'orders'
type OrderRepository struct {
	db *pgxpool.Pool
}

// NewOrderRepository membuat instance baru dari OrderRepository
func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder membuat order baru, seat yang dipilih, dan menambahkan poin ke profil pengguna
func (o *OrderRepository) CreateOrder(ctx context.Context, orderData models.CreateOrderRequest) (models.Order, error) {
	// Mulai transaksi untuk memastikan semua operasi atomic
	tx, err := o.db.Begin(ctx)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return models.Order{}, err
	}
	defer tx.Rollback(ctx)

	// Persiapkan data untuk tabel 'orders'
	newOrder := models.Order{
		IDUsers:     orderData.IDUsers,
		IDSchedule:  orderData.IDSchedule,
		IDPayment:   orderData.IDPayment,
		FullName:    orderData.FullName,
		Email:       orderData.Email,
		PhoneNumber: orderData.PhoneNumber,
		TotalPrice:  orderData.TotalPrice,
		IsPaid:      false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 1. Insert data ke tabel 'orders'
	queryOrders := `
		INSERT INTO orders (users_id, schedule_id, payment_method_id, total_price, ispaid, created_at, update_at, full_name, email, phone_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	err = tx.QueryRow(ctx, queryOrders,
		newOrder.IDUsers,
		newOrder.IDSchedule,
		newOrder.IDPayment,
		newOrder.TotalPrice,
		newOrder.IsPaid,
		newOrder.CreatedAt,
		newOrder.UpdatedAt,
		newOrder.FullName,
		newOrder.Email,
		newOrder.PhoneNumber).Scan(&newOrder.ID)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return models.Order{}, err
	}

	// 2. Ambil ID kursi dari kode kursi (seats_code)
	var seatsID []int
	for _, seatCode := range orderData.SeatsCode {
		var seatID int
		err := tx.QueryRow(ctx, "SELECT id FROM seats WHERE seats_code = $1", seatCode).Scan(&seatID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return models.Order{}, errors.New("seat with code " + seatCode + " not found")
			}
			log.Println("[ERROR] : ", err.Error())
			return models.Order{}, err
		}
		seatsID = append(seatsID, seatID)
	}

	// 3. Insert data ke tabel 'order_seats'
	queryOrderSeats := `
		INSERT INTO order_seats (orders_id, seats_id)
		VALUES ($1, $2)`

	for _, seatID := range seatsID {
		_, err := tx.Exec(ctx, queryOrderSeats, newOrder.ID, seatID)
		if err != nil {
			log.Println("[ERROR] : ", err.Error())
			return models.Order{}, err
		}
	}

	// 4. Update poin di tabel 'profile'
	queryUpdatePoints := `
		UPDATE profile
		SET point = point + 50
		WHERE users_id = $1`

	_, err = tx.Exec(ctx, queryUpdatePoints, newOrder.IDUsers)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return models.Order{}, err
	}

	// Commit transaksi jika semua operasi berhasil
	return newOrder, tx.Commit(ctx)
}

// GetOrderHistory mengambil riwayat order berdasarkan user ID.
func (o *OrderRepository) GetOrderHistory(ctx context.Context, userID int) ([]models.OrderHistory, error) {
	query := `
        SELECT
            o.id,
            o.total_price,
            o.ispaid,
            o.created_at,
            s.date AS show_date,
            t.time AS show_time,
            m.title AS movie_title,
            m.poster_path AS movie_image,
            ci.name AS cinema_name,  -- INI PERBAIKANNYA
            l.location AS cinema_location
        FROM orders AS o
        JOIN schedule AS s ON o.schedule_id = s.id
        JOIN movie AS m ON s.movie_id = m.id
        JOIN cinema AS ci ON s.cinema_id = ci.id
        JOIN location AS l ON s.location_id = l.id
        JOIN time AS t ON s.time_id = t.id
        WHERE o.users_id = $1
        ORDER BY o.created_at DESC
    `
	rows, err := o.db.Query(ctx, query, userID)
	if err != nil {
		log.Println("[ERROR]: Gagal mengambil riwayat order:", err.Error())
		return nil, err
	}
	defer rows.Close()

	var orderHistories []models.OrderHistory
	for rows.Next() {
		var history models.OrderHistory
		if err := rows.Scan(
			&history.OrderID,
			&history.TotalPrice,
			&history.IsPaid,
			&history.CreatedAt,
			&history.Schedule.Date,
			&history.Schedule.Time,
			&history.Movie.Title,
			&history.Movie.Image,
			&history.Cinema.CinemaName,
			&history.Cinema.CinemaLocation,
		); err != nil {
			log.Println("[ERROR]: Gagal memindai baris riwayat order:", err.Error())
			return nil, err
		}

		// Ambil kursi yang dipesan untuk order ini
		querySeats := `
            SELECT s.seats_code
            FROM order_seats AS os
            JOIN seats AS s ON os.seats_id = s.id
            WHERE os.orders_id = $1
        `
		seatRows, err := o.db.Query(ctx, querySeats, history.OrderID)
		if err != nil {
			log.Println("[ERROR]: Gagal mengambil kursi yang dipesan:", err.Error())
			return nil, err
		}
		defer seatRows.Close()

		var seats []string
		for seatRows.Next() {
			var seatCode string
			if err := seatRows.Scan(&seatCode); err != nil {
				log.Println("[ERROR]: Gagal memindai kode kursi:", err.Error())
				return nil, err
			}
			seats = append(seats, seatCode)
		}
		history.Seats = seats
		orderHistories = append(orderHistories, history)
	}

	return orderHistories, nil
}

func (o *OrderRepository) GetOrderByID(ctx context.Context, orderID int) (models.OrderHistory, error) {
	query := `
        SELECT
            o.id,
            o.total_price,
            o.ispaid,
            o.created_at,
            s.date AS show_date,
            t.time AS show_time,
            m.title AS movie_title,
            m.poster_path AS movie_image,
            ci.name AS cinema_name,
            l.location AS cinema_location
        FROM orders AS o
        JOIN schedule AS s ON o.schedule_id = s.id
        JOIN movie AS m ON s.movie_id = m.id
        JOIN cinema AS ci ON s.cinema_id = ci.id
        JOIN location AS l ON s.location_id = l.id
        JOIN time AS t ON s.time_id = t.id
        WHERE o.id = $1
        LIMIT 1
    `
	var history models.OrderHistory
	row := o.db.QueryRow(ctx, query, orderID)
	if err := row.Scan(
		&history.OrderID,
		&history.TotalPrice,
		&history.IsPaid,
		&history.CreatedAt,
		&history.Schedule.Date,
		&history.Schedule.Time,
		&history.Movie.Title,
		&history.Movie.Image,
		&history.Cinema.CinemaName,
		&history.Cinema.CinemaLocation,
	); err != nil {
		if err == pgx.ErrNoRows {
			return models.OrderHistory{}, errors.New("order not found")
		}
		log.Println("[ERROR]: Gagal memindai baris order:", err.Error())
		return models.OrderHistory{}, err
	}

	// Ambil kursi yang dipesan untuk order ini
	querySeats := `
        SELECT s.seats_code
        FROM order_seats AS os
        JOIN seats AS s ON os.seats_id = s.id
        WHERE os.orders_id = $1
    `
	seatRows, err := o.db.Query(ctx, querySeats, history.OrderID)
	if err != nil {
		log.Println("[ERROR]: Gagal mengambil kursi yang dipesan:", err.Error())
		return models.OrderHistory{}, err
	}
	defer seatRows.Close()

	var seats []string
	for seatRows.Next() {
		var seatCode string
		if err := seatRows.Scan(&seatCode); err != nil {
			log.Println("[ERROR]: Gagal memindai kode kursi:", err.Error())
			return models.OrderHistory{}, err
		}
		seats = append(seats, seatCode)
	}
	history.Seats = seats

	return history, nil
}
