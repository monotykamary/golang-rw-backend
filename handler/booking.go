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
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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

func (h *Handler) QueueBookingHandler(cfg config.Config, s repo.DBRepo) func(c echo.Context) error {
	ctx := context.Background()
	addr := cfg.RedisHost + ":" + cfg.RedisPort
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	return func(c echo.Context) error {
		event := c.QueryParam("event")
		userId, _ := uuid.FromString(c.QueryParam("userId"))
		bookingId, _ := uuid.FromString(c.QueryParam("bookingId"))

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
