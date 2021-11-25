package handler

import (
	inerr "errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/monotykamary/golang-rw-backend/model"
	"github.com/monotykamary/golang-rw-backend/model/errors"
	"github.com/monotykamary/golang-rw-backend/util"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRequest struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type getUserInfo struct {
	Id uuid.UUID `json:"id"`
}

type getUserInfoResponse struct {
	Data getUserInfo `json:"data"`
}

type getUsersResponse struct {
	Data []getUserInfo `json:"data"`
}

// GetUsersHandler
// @Summary get all user info
// @Description get all user info
// @Accept	json
// @Produce  json
// @Success 200 {object} handler.getUsersResponse	"ok"
// @Failure 400 {object} errors.Error
// @Router /api/v1/users [get]
func (h *Handler) GetUsersHandler(c echo.Context) error {
	users, err := h.repo.User.GetAll(h.store)
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.GetUsersHandler] User.GetAll() no users found")
			return util.HandleError(c, errors.ErrUserNotfound)
		}
		zap.L().Sugar().Infof("[handler.GetUsersHandler] User.GetAll()")
		return util.HandleError(c, errors.ErrInternalServerError)
	}

	usersMap := make([]getUserInfo, 0)

	for _, user := range users {
		usersMap = append(usersMap, getUserInfo{Id: user.Id})
	}

	return c.JSON(http.StatusOK, &getUsersResponse{Data: usersMap})
}

// GetUserInfoHandler
// @Summary get user info
// @Description get user info
// @Accept	json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} handler.getUserInfoResponse	"ok"
// @Failure 400 {object} errors.Error
// @Router /api/v1/users/{id} [get]
func (h *Handler) GetUserInfoHandler(c echo.Context) error {
	id := c.QueryParam("id")
	user, err := h.repo.User.GetById(h.store, id)
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.GetUserInfoHandler] User.GetById() user not found")
			return util.HandleError(c, errors.ErrUserNotfound)
		}
		zap.L().Sugar().Infof("[handler.GetUserInfoHandler] User.GetById()")
		return util.HandleError(c, errors.ErrInternalServerError)
	}

	return c.JSON(http.StatusOK, &getUserInfoResponse{Data: getUserInfo{Id: user.Id}})
}

// RegisterUserHandler
// @Summary register as new user
// @Description register as new user
// @Accept	json
// @Produce  json
// @Success 200 {object} handler.getUserInfoResponse	"ok"
// @Failure 400 {object} errors.Error
// @Router /api/v1/users/register [post]
func (h *Handler) RegisterUserHandler(c echo.Context) error {
	tx, done := h.store.NewTransaction()
	user, err := h.repo.User.Create(tx, model.User{})
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.RegisterUserHandler] User.Create() user not found")
			return util.HandleError(c, errors.ErrUserNotfound)
		}
		zap.L().Sugar().Infof("[handler.RegisterUserHandler] User.Create()")
		return util.HandleError(c, done(errors.ErrInternalServerError))
	}

	return c.JSON(http.StatusOK, &getUserInfoResponse{Data: getUserInfo{Id: user.Id}})
}

// UpdateUserHandler
// @Summary update user info
// @Description update user info
// @Accept	json
// @Produce  json
// @Param body body handler.userRequest true "user request"
// @Success 200 {object} handler.getUserInfoResponse	"ok"
// @Failure 400 {object} errors.Error
// @Router /api/v1/users/update [post]
func (h *Handler) UpdateUserHandler(c echo.Context) error {
	var request userRequest
	if err := c.Bind(&request); err != nil {
		zap.L().Sugar().Infof("[handler.UpdateUserHandler] c.Bind()")
		return util.HandleError(c, err)
	}

	user, err := h.repo.User.Update(h.store, model.User{}, request.Email, request.Id.String())
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Sugar().Infof("[handler.UpdateUserHandler] User.Update() user not found")
			return util.HandleError(c, errors.ErrUserNotfound)
		}
		zap.L().Sugar().Infof("[handler.UpdateUserHandler] User.Update()")
		return util.HandleError(c, errors.ErrInternalServerError)
	}

	return c.JSON(http.StatusOK, &getUserInfoResponse{Data: getUserInfo{Id: user.Id}})
}
