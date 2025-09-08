package models

// SeatStruct mewakili struktur data untuk kursi tunggal
type SeatStruct struct {
	ID        int    `json:"id"`
	SeatsCode string `json:"seats_code"`
}

// AvailableSeatsResponse adalah struktur respons untuk daftar kursi yang tersedia
type AvailableSeatsResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    []SeatStruct `json:"data"`
}
