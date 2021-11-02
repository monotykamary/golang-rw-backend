package errors

import "net/http"

// Error in server
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

// NewStringError new a error with message
func NewStringError(msg string, code int) error {
	return Error{
		Code:    code,
		Message: msg,
	}
}

// NewStringError new a error with message

var (
	ErrInternalServerError      = NewStringError("Internal server error", http.StatusInternalServerError)
	ErrIncorrectEmailOrPassword = NewStringError("Email or password is incorrect", http.StatusBadRequest)
)
