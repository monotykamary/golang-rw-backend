package pg

import "github.com/monotykamary/golang-rw-backend/repo"

func NewRepo() *repo.Repo {
	return &repo.Repo{
		User:    NewUserRepo(),
		Booking: NewBookingRepo(),
	}
}
