// package models

// import "time"

// // CreateOrderRequest adalah model untuk body permintaan API
// type CreateOrderRequest struct {
// 	IDUsers     int      `json:"users_id" binding:"required"`
// 	IDSchedule  int      `json:"schedule_id" binding:"required"`
// 	IDPayment   int      `json:"payment_id" binding:"required"`
// 	FullName    string   `json:"full_name" binding:"required"`
// 	Email       string   `json:"email" binding:"required,email"`
// 	PhoneNumber string   `json:"phone_number" binding:"required"`
// 	TotalPrice  int      `json:"total_price" binding:"required"`
// 	SeatsCode   []string `json:"seats_code" binding:"required"`
// }

// // Order adalah model untuk tabel 'orders'
// type Order struct {
// 	ID          int       `json:"id"`
// 	IDUsers     int       `json:"users_id"`
// 	IDSchedule  int       `json:"schedule_id"`
// 	IDPayment   int       `json:"payment_method_id"`
// 	FullName    string    `json:"full_name"`
// 	Email       string    `json:"email"`
// 	PhoneNumber string    `json:"phone_number"`
// 	TotalPrice  int       `json:"total_price"`
// 	IsPaid      bool      `json:"ispaid"`
// 	CreatedAt   time.Time `json:"created_at"`
// 	UpdatedAt   time.Time `json:"updated_at"`
// 	SeatsCode   []string  `json:"seats_code"`
// }

// // ResponseCreateOrder adalah model untuk respons API setelah berhasil membuat order
// type ResponseCreateOrder struct {
// 	Status  string `json:"status"`
// 	Message string `json:"message"`
// 	Data    Order  `json:"data"`
// }

// // OrderHistory mewakili data riwayat order yang akan ditampilkan kepada pengguna
// type OrderHistory struct {
// 	OrderID    int       `json:"users_id"`
// 	TotalPrice int       `json:"total_price"`
// 	IsPaid     bool      `json:"is_paid"`
// 	CreatedAt  time.Time `json:"created_at"`
// 	Seats      []string  `json:"seats"`
// 	Schedule   Schedule  `json:"schedule"`
// 	Movie      Movie     `json:"movie"`
// 	Cinema     Cinema    `json:"cinema"`
// }

// // Schedule adalah model untuk informasi jadwal
// // type Schedule struct {
// // 	Date time.Time `json:"date"`
// // 	Time time.Time `json:"time"`
// // }

// // Movie adalah model untuk informasi film
// type Movie struct {
// 	Title string `json:"title"`
// 	Image string `json:"image"`
// }

// // Cinema adalah model untuk informasi bioskop
// type Cinema struct {
// 	CinemaName     string `json:"cinema_name"`
// 	CinemaLocation string `json:"cinema_location"`
// }

// // OrderHistoryResponse adalah model untuk respons API riwayat order
//
//	type OrderHistoryResponse struct {
//		Status  string         `json:"status"`
//		Message string         `json:"message"`
//		Data    []OrderHistory `json:"data"`
//	}
package models

import "time"

// Order adalah model untuk tabel 'orders'
type Order struct {
	ID          int       `json:"id"`
	IDUsers     int       `json:"users_id"`
	IDSchedule  int       `json:"schedule_id"`
	IDPayment   int       `json:"payment_method_id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	TotalPrice  int       `json:"total_price"`
	IsPaid      bool      `json:"ispaid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SeatsCode   []string  `json:"seats_code"`
}

// ResponseCreateOrder adalah model untuk respons API setelah berhasil membuat order
type ResponseCreateOrder struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    Order  `json:"data"`
}

// OrderHistory mewakili data riwayat order yang akan ditampilkan kepada pengguna
type OrderHistory struct {
	OrderID    int       `json:"users_id"`
	TotalPrice int       `json:"total_price"`
	IsPaid     bool      `json:"is_paid"`
	CreatedAt  time.Time `json:"created_at"`
	Seats      []string  `json:"seats"`
	Schedule   Schedule  `json:"schedule"`
	Movie      Movie     `json:"movie"`
	Cinema     Cinema    `json:"cinema"`
}

// Schedule adalah model untuk informasi jadwal
// Movie adalah model untuk informasi film
type Movie struct {
	Title string `json:"title"`
	Image string `json:"image"`
}

// Cinema adalah model untuk informasi bioskop
type Cinema struct {
	CinemaName     string `json:"cinema_name"`
	CinemaLocation string `json:"cinema_location"`
}

// OrderHistoryResponse adalah model untuk respons API riwayat order
type OrderHistoryResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []OrderHistory `json:"data"`
}
