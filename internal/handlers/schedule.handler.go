package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
)

// ScheduleHandler adalah handler untuk semua rute jadwal
type ScheduleHandler struct {
	scheduleRepo *repositories.ScheduleRepository
}

// NewScheduleHandlers membuat instance baru dari ScheduleHandler
func NewScheduleHandlers(scheduleRepo *repositories.ScheduleRepository) *ScheduleHandler {
	return &ScheduleHandler{scheduleRepo: scheduleRepo}
}

// GetAllSchedules mengambil semua jadwal
// @summary                 Get all schedules
// @router                  /schedules [get]
// @Description             Get all available movie schedules
// @Tags                    Schedules
// @accept                  json
// @produce                 json
// @success                 200 {object} models.ResponseSchedule
// @failure                 500 {object} models.Response
func (s *ScheduleHandler) GetAllSchedules(ctx *gin.Context) {
	schedules, err := s.scheduleRepo.GetAllSchedules(ctx.Request.Context())
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "failed to get schedules",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseScheduleByMovie{
		Status:  "berhasil",
		Message: "successfully fetched all schedules",
		Data:    schedules,
	})
}

func (s *ScheduleHandler) GetSchedulesByMovieID(ctx *gin.Context) {
	movieID, err := strconv.Atoi(ctx.Param("movieId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid movie ID",
		})
		return
	}

	schedules, err := s.scheduleRepo.GetSchedulesByMovieID(ctx.Request.Context(), movieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "failed to get schedules",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseScheduleByMovie{
		Status:  "berhasil",
		Message: "successfully fetched schedules",
		Data:    schedules,
	})
}

// GetScheduleByMovieId mengambil jadwal berdasarkan ID film
// @summary                 Get schedules by movie ID
// @router                  /schedules/movie/{id} [get]
// @Description             Get all schedules for a specific movie by ID
// @Tags                    Schedules
// @Param                   id path int true "Movie ID"
// @accept                  json
// @produce                 json
// @success                 200 {object} models.ResponseScheduleByMovie
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
// @failure                 500 {object} models.Response
// func (s *ScheduleHandler) GetScheduleByMovieId(ctx *gin.Context) {
// 	paramId := ctx.Param("id")
// 	movieID, err := strconv.Atoi(paramId)
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "invalid movie ID",
// 		})
// 		return
// 	}

// 	schedules, err := s.scheduleRepo.GetScheduleByMovieId(ctx.Request.Context(), movieID)
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		if err == pgx.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, models.Response{
// 				Status:  "gagal",
// 				Message: "no schedules found for this movie",
// 			})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "failed to get schedules",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, models.ResponseScheduleByMovie{
// 		Status:  "berhasil",
// 		Message: "successfully fetched schedules for the movie",
// 		Data:    schedules,
// 	})
// }
