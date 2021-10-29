package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/monotykamary/golang-rw-backend/config"
	"github.com/monotykamary/golang-rw-backend/handler"
	"github.com/monotykamary/golang-rw-backend/repo"
)

func NewRoutes(e *echo.Echo, h *handler.Handler, cfg config.Config, s repo.DBRepo) {
	apiV1Group := e.Group("/api/v1")

	userGroup := apiV1Group.Group("/users")
	{
		userGroup.GET("", h.GetUsersHandler)
		userGroup.GET("/:id", h.GetUserInfoHandler)
		userGroup.POST("/register", h.RegisterUserHandler)
		userGroup.POST("/update", h.UpdateUserHandler)
	}

	bookingGroup := apiV1Group.Group("/booking")
	{
		bookingGroup.GET("", h.GetBookingsHandler)
		bookingGroup.GET("/:id", h.GetBookingInfoHandler)
		bookingGroup.GET("/queue", h.QueueBookingHandler(cfg, s))
	}
}
