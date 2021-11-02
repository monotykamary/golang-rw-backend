package errors

import "net/http"

var (
	ErrBookingNotfound       = NewStringError("Booking not found", http.StatusNotFound)
	ErrBookingsNotfound      = NewStringError("Bookings not found", http.StatusNotFound)
	ErrBookingAlreadyExisted = NewStringError("Booking already existed", http.StatusConflict)
)
