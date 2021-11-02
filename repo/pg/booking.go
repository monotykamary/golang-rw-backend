package pg

import (
	"github.com/monotykamary/golang-rw-backend/model"
	"github.com/monotykamary/golang-rw-backend/repo"
)

type bookingRepo struct{}

func NewBookingRepo() repo.BookingRepo {
	return &bookingRepo{}
}

func (t *bookingRepo) GetById(s repo.DBRepo, id string) (*model.Booking, error) {
	db := s.DB()
	booking := model.Booking{}
	return &booking, db.First(&booking, "id = ?", id).Error
}

func (r *bookingRepo) GetAll(s repo.DBRepo) ([]model.Booking, error) {
	db := s.DB()
	var results []model.Booking
	return results, db.Table("user").Select("id").Find(&results).Error
}
