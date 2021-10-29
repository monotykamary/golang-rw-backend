package handler

import (
	inerr "errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/monotykamary/golang-rw-backend/model"
	"github.com/monotykamary/golang-rw-backend/model/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type getUserInfo struct {
	Id uuid.UUID `json:"id"`
}

type getUserInfoResponse struct {
	Data getUserInfo `json:"data"`
}

type getUsersResponse struct {
	Data []getUserInfo `json:"data"`
}

func (h *Handler) GetUsersHandler(c echo.Context) error {
	users, err := h.repo.User.GetAll(h.store)
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.GetUsersHandler] User.GetAll() no users found")
			return errors.ErrUsersNotfound
		}
		zap.L().Sugar().Infof("[handler.GetUsersHandler] User.GetAll()")
		return errors.ErrInternalServerError
	}

	usersMap := make([]getUserInfo, 0)

	for _, user := range users {
		usersMap = append(usersMap, getUserInfo{Id: user.Id})
	}

	return c.JSON(http.StatusOK, &getUsersResponse{Data: usersMap})
}

func (h *Handler) GetUserInfoHandler(c echo.Context) error {
	id := c.QueryParam("id")
	user, err := h.repo.User.GetById(h.store, id)
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.GetUserInfoHandler] User.GetById() user not found")
			return errors.ErrUserNotfound
		}
		zap.L().Sugar().Infof("[handler.GetUserInfoHandler] User.GetById()")
		return errors.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &getUserInfoResponse{Data: getUserInfo{Id: user.Id}})
}

func (h *Handler) RegisterUserHandler(c echo.Context) error {
	user, err := h.repo.User.Create(h.store, model.User{})
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.RegisterUserHandler] User.Create() user not found")
			return errors.ErrUserNotfound
		}
		zap.L().Sugar().Infof("[handler.RegisterUserHandler] User.Create()")
		return errors.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &getUserInfoResponse{Data: getUserInfo{Id: user.Id}})
}

func (h *Handler) UpdateUserHandler(c echo.Context) error {
	user, err := h.repo.User.Update(h.store, model.User{})
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.UpdateUserHandler] User.Update() user not found")
			return errors.ErrUserNotfound
		}
		zap.L().Sugar().Infof("[handler.UpdateUserHandler] User.Update()")
		return errors.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, &getUserInfoResponse{Data: getUserInfo{Id: user.Id}})
}
