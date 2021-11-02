package repo

import "github.com/monotykamary/golang-rw-backend/model"

type Repo struct {
	User    UserRepo
	Booking BookingRepo
}

type UserRepo interface {
	GetById(s DBRepo, id string) (*model.User, error)
	Create(s DBRepo, param model.User) (*model.User, error)
	Update(s DBRepo, param model.User, value string, id string) (*model.User, error)
	GetAll(s DBRepo) ([]model.User, error)
}

type BookingRepo interface {
	GetById(s DBRepo, id string) (*model.Booking, error)
	GetAll(s DBRepo) ([]model.Booking, error)
}
