package redis

import "github.com/google/uuid"

type BookingQueueRequest struct {
	Event     string    `json:"event"`
	UserId    uuid.UUID `json:"userId"`
	BookingId uuid.UUID `json:"bookingId"`
}
