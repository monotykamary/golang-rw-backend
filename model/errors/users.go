package errors

import "net/http"

var (
	ErrUserNotfound       = NewStringError("User not found", http.StatusNotFound)
	ErrUsersNotfound      = NewStringError("Users not found", http.StatusNotFound)
	ErrUserAlreadyExisted = NewStringError("User already existed", http.StatusConflict)
)
