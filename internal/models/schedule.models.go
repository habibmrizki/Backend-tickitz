package models

import "time"

type Schedule struct {
	Id       int `json:"id"`
	MovieID  int `json:"movie_id"`
	CinemaID int `json:"cinema_id"`
	Location int `json:"location"`
	TimeID   int `json:"time_id"`
	Time     time.Time
	Date     time.Time `json:"date"`
}

// ScheduleDetails adalah model untuk respons yang lebih detail
type ScheduleDetails struct {
	ID         int       `json:"id"`
	Date       time.Time `json:"date"`
	CinemaName string    `json:"cinema_name"`
	Location   string    `json:"location"`
	Time       time.Time `json:"time"`
}

// ResponseSchedule adalah model respons API untuk daftar jadwal
type ResponseSchedule struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []ScheduleDetails `json:"data"`
}

// ResponseScheduleByMovie adalah model respons API untuk jadwal berdasarkan film
type ResponseScheduleByMovie struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []ScheduleDetails `json:"data"`
}
