package handler

import (
	"context"
	inerr "errors"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/monotykamary/golang-rw-backend/config"
	"github.com/monotykamary/golang-rw-backend/model/errors"
	"github.com/monotykamary/golang-rw-backend/repo"
	"github.com/monotykamary/golang-rw-backend/util"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type bookingQueueRequest struct {
	Event     string    `json:"event"`
	UserId    uuid.UUID `json:"userId"`
	BookingId uuid.UUID `json:"bookingId"`
}

type getBookingQueueInfo struct {
	Event     string    `json:"event"`
	UserId    uuid.UUID `json:"userId"`
	BookingId uuid.UUID `json:"bookingId"`
}

type getBookingInfo struct {
	Id     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type getBookingQueueInfoResponse struct {
	Data getBookingQueueInfo `json:"data"`
}

type getBookingInfoResponse struct {
	Data getBookingInfo `json:"data"`
}

type getBookingsResponse struct {
	Data []getBookingInfo `json:"data"`
}

// GetBookingsHandler
// @Summary get all booking info
// @Description get all booking info
// @Accept	json
// @Produce  json
// @Success 200 {object} handler.getUsersResponse	"ok"
// @Failure 400 {object} errors.Error
// @Router /api/v1/booking [get]
func (h *Handler) GetBookingsHandler(c echo.Context) error {
	bookings, err := h.repo.Booking.GetAll(h.store)
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.GetBookingsHandler] Booking.GetAll() no bookings found")
			return errors.ErrBookingsNotfound
		}
		zap.L().Sugar().Infof("[handler.GetBookingsHandler] Booking.GetAll()")
		return errors.ErrInternalServerError
	}

	bookingsMap := make([]getBookingInfo, 0)

	for _, booking := range bookings {
		bookingsMap = append(bookingsMap, getBookingInfo{Id: booking.Id, Status: booking.Status})
	}

	return c.JSON(http.StatusOK, &getBookingsResponse{Data: bookingsMap})
}

// GetBookingInfoHandler
// @Summary get booking info
// @Description get booking info
// @Accept	json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} handler.getBookingInfoResponse	"ok"
// @Failure 400 {object} errors.Error
// @Router /api/v1/booking/{id} [get]
func (h *Handler) GetBookingInfoHandler(c echo.Context) error {
	id := c.QueryParam("id")
	booking, err := h.repo.Booking.GetById(h.store, id)
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.GetBookingInfoHandler] Booking.GetById() booking not found")
			return errors.ErrBookingNotfound
		}
		zap.L().Sugar().Infof("[handler.GetBookingInfoHandler] Booking.GetById()")
		return errors.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &getBookingInfoResponse{Data: getBookingInfo{Id: booking.Id, Status: booking.Status}})
}

// QueueBookingHandler
// @Summary queue booking item
// @Description queue booking item
// @Accept	json
// @Produce  json
// @Param body body handler.bookingQueueRequest true "booking queue request"
// @Success 200 {object} handler.getBookingInfoResponse	"ok"
// @Failure 400 {object} errors.Error
// @Router /api/v1/booking/queue [post]
func (h *Handler) QueueBookingHandler(cfg config.Config, s repo.DBRepo) func(c echo.Context) error {
	ctx := context.Background()
	addr := cfg.RedisHost + ":" + cfg.RedisPort
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	return func(c echo.Context) error {
		var request bookingQueueRequest
		if err := c.Bind(&request); err != nil {
			zap.L().Sugar().Infof("[handler.UpdateUserHandler] c.Bind()")
			return util.HandleError(c, err)
		}

		event := request.Event
		userId := request.UserId
		bookingId := request.BookingId

		_, err := redisClient.XAdd(ctx, &redis.XAddArgs{
			Stream: "booking",
			Values: map[string]interface{}{
				"event":     event,
				"userId":    userId,
				"bookingId": bookingId,
			},
		}).Result()

		if err != nil {
			zap.L().Error("cannot queue booking", zap.Error(err))
		}

		return c.JSON(http.StatusOK, &getBookingQueueInfoResponse{
			Data: getBookingQueueInfo{event, userId, bookingId},
		})
	}
}
