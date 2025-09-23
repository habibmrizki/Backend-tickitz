package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
)

type SeatHandler struct {
	seatRepo *repositories.SeatRepository
}

func NewSeatHandler(seatRepo *repositories.SeatRepository) *SeatHandler {
	return &SeatHandler{seatRepo: seatRepo}
}

// GetAvailableSeats mengambil daftar kursi yang tersedia untuk jadwal film
// @summary                 Get available seats
// @router                  /seats/:scheduleId/available [get]
// @Description             Get a list of available seats for a specific schedule
// @Tags                    Seats
// @Param                   scheduleId path int true "Schedule ID"
// @accept                  json
// @produce                 json
// @failure                 400 {object} models.Response
// @failure                 500 {object} models.Response
// @success                 200 {object} models.AvailableSeatsResponse
func (h *SeatHandler) GetAvailableSeats(ctx *gin.Context) {
	scheduleIDStr := ctx.Param("scheduleId")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid schedule ID",
		})
		return
	}

	seats, err := h.seatRepo.GetAvailableSeats(ctx.Request.Context(), scheduleID)
	if err != nil {
		log.Println("[ERROR]: Failed to get available seats:", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.AvailableSeatsResponse{
		Status:  "berhasil",
		Message: "successfully get available seats",
		Data:    seats,
	})
}

// GetSeatByID mengambil satu kursi berdasarkan ID
// @summary                 Get seat by ID
// @router                  /seats/:seatId [get]
// @Description             Get a single seat by its ID
// @Tags                    Seats
// @Param                   seatId path int true "Seat ID"
// @accept                  json
// @produce                 json
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
// @failure                 500 {object} models.Response
// @success                 200 {object} models.SeatStruct
// func (h *SeatHandler) GetSeatByID(ctx *gin.Context) {
// 	seatIDStr := ctx.Param("seatId")
// 	seatID, err := strconv.Atoi(seatIDStr)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "ID kursi tidak valid",
// 		})
// 		return
// 	}

// 	seat, err := h.seatRepo.GetSeatByID(ctx.Request.Context(), seatID)
// 	if err != nil {
// 		log.Println("[ERROR]: Gagal mengambil kursi:", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "Kesalahan server internal",
// 		})
// 		return
// 	}

// 	if seat == nil {
// 		ctx.JSON(http.StatusNotFound, models.Response{
// 			Status:  "gagal",
// 			Message: "Kursi tidak ditemukan",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, models.ResponseWithData{
// 		Status:  "berhasil",
// 		Message: "Berhasil mendapatkan detail kursi",
// 		Data:    seat,
// 	})
// }
