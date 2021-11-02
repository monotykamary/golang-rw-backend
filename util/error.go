package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/monotykamary/golang-rw-backend/model/errors"
)

// ParseErrorCode parse error code from errors.Error
func ParseErrorCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch arg := err.(type) {
	case *errors.Error:
		return int(arg.Code)

	case errors.Error:
		return int(arg.Code)

	case error:
		return http.StatusInternalServerError

	default:
		return http.StatusOK
	}
}

func HandleError(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	switch arg := err.(type) {

	case *errors.Error:
		return c.JSON(int(arg.Code), arg)

	case errors.Error:
		return c.JSON(int(arg.Code), arg)

	case *echo.HTTPError:
		return c.JSON(int(arg.Code), errors.Error{
			Code:    arg.Code,
			Message: arg.Message.(string),
		})
	case error:
		return c.JSON(http.StatusInternalServerError, errors.Error{
			Code:    http.StatusInternalServerError,
			Message: arg.Error(),
		})
	}
	return err
}
