// package handlers

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/habibmrizki/back-end-tickitz/internal/models"
// 	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
// )

// // OrderHandler adalah handler untuk semua rute order
// type OrderHandler struct {
// 	orderRepo *repositories.OrderRepository
// }

// // NewOrderHandlers membuat instance baru dari OrderHandler
// func NewOrderHandlers(orderRepo *repositories.OrderRepository) *OrderHandler {
// 	return &OrderHandler{orderRepo: orderRepo}
// }

// // CreateOrder membuat order baru
// // @summary                 Create a new order
// // @router                  /orders [post]
// // @Description             Create a new movie ticket order
// // @Tags                    Orders
// // @Param                   order body models.CreateOrderRequest true "Order data"
// // @accept                  json
// // @produce                 json
// // @success                 201 {object} models.ResponseCreateOrder
// // @failure                 400 {object} models.Response
// // @failure                 500 {object} models.Response
// func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
// 	var requestBody models.CreateOrderRequest
// 	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "invalid data sent",
// 		})
// 		return
// 	}

// 	order, err := o.orderRepo.CreateOrder(ctx.Request.Context(), requestBody)
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: "failed to create order",
// 		})
// 		return
// 	}

//		ctx.JSON(http.StatusCreated, models.ResponseCreateOrder{
//			Status:  "berhasil",
//			Message: "order created successfully",
//			Data:    order,
//		})
//	}
package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
)

// OrderHandler adalah handler untuk semua rute order
type OrderHandler struct {
	orderRepo *repositories.OrderRepository
}

// NewOrderHandlers membuat instance baru dari OrderHandler
func NewOrderHandlers(orderRepo *repositories.OrderRepository) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo}
}

// CreateOrder membuat order baru
// @summary                 Create a new order
// @router                  /orders [post]
// @Description             Create a new movie ticket order
// @Tags                    Orders
// @Param                   order body models.Order true "Order data"
// @accept                  json
// @produce                 json
// @success                 201 {object} models.ResponseCreateOrder
// @failure                 400 {object} models.Response
// @failure                 500 {object} models.Response
func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
	var requestBody models.Order
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid data sent",
		})
		return
	}
	log.Printf("[Handler] Received IsPaid: %v", requestBody.IsPaid)
	order, err := o.orderRepo.CreateOrder(ctx.Request.Context(), requestBody)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		// Mengembalikan pesan error dari repository jika spesifik
		if err.Error() == "one or more seats not found" {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  "gagal",
				Message: "Kursi yang dipilih tidak valid",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "Gagal membuat order",
		})
		return
	}

	// Tambahkan SeatsCode ke dalam objek 'order' sebelum mengirim respons.
	order.SeatsCode = requestBody.SeatsCode

	ctx.JSON(http.StatusCreated, models.ResponseCreateOrder{
		Status:  "berhasil",
		Message: "order created successfully",
		Data:    order,
	})
}

// GetOrderHistory mengambil riwayat order berdasarkan user ID
// @summary                 Get order history
// @router                  /orders/history/{userId} [get]
// @Description             Get a list of all orders for a specific user
// @Tags                    Orders
// @Param                   userId path int true "User ID"
// @accept                  json
// @produce                 json
// @success                 200 {object} models.OrderHistoryResponse
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
// @failure                 500 {object} models.Response
func (o *OrderHandler) GetOrderHistory(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "ID pengguna tidak valid",
		})
		return
	}

	orderHistory, err := o.orderRepo.GetOrderHistory(ctx.Request.Context(), userID)
	if err != nil {
		log.Println("[ERROR]: Gagal mengambil riwayat order:", err.Error())
		// Jika tidak ada order, kembalikan 404 atau pesan kosong.
		// Di sini  kembalikan respons 200 dengan data kosong jika tidak ada error lain.
		if len(orderHistory) == 0 {
			ctx.JSON(http.StatusOK, models.OrderHistoryResponse{
				Status:  "berhasil",
				Message: "Tidak ada riwayat order ditemukan untuk pengguna ini",
				Data:    []models.OrderHistory{},
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "Kesalahan server internal",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.OrderHistoryResponse{
		Status:  "berhasil",
		Message: "Riwayat order berhasil diambil",
		Data:    orderHistory,
	})
}

// GetOrderByID mengambil detail order berdasarkan order ID
// @summary                 Get order by ID
// @router                  /orders/{orderId} [get]
// @Description             Get details of a specific order by its ID
// @Tags                    Orders
// @Param                   orderId path int true "Order ID"
// @accept                  json
// @produce                 json
// @success                 200 {object} models.OrderHistoryResponse
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
// @failure                 500 {object} models.Response
func (o *OrderHandler) GetOrderByID(ctx *gin.Context) {
	orderIDStr := ctx.Param("orderId")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "ID order tidak valid",
		})
		return
	}

	orderHistory, err := o.orderRepo.GetOrderByID(ctx.Request.Context(), orderID)
	if err != nil {
		log.Println("[ERROR]: Gagal mengambil order:", err.Error())
		if err.Error() == "order not found" {
			ctx.JSON(http.StatusNotFound, models.Response{
				Status:  "gagal",
				Message: "Order tidak ditemukan",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "Kesalahan server internal",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.OrderHistoryResponse{
		Status:  "berhasil",
		Message: "Order berhasil diambil",
		Data:    []models.OrderHistory{orderHistory},
	})
}
