package pg

import (
	"github.com/monotykamary/golang-rw-backend/model"
	"github.com/monotykamary/golang-rw-backend/repo"
)

type userRepo struct{}

func NewUserRepo() repo.UserRepo {
	return &userRepo{}
}

func (t *userRepo) GetById(s repo.DBRepo, id string) (*model.User, error) {
	db := s.DB()
	user := model.User{}
	return &user, db.First(&user, "id = ?", id).Error
}

func (t *userRepo) Create(s repo.DBRepo, param model.User) (*model.User, error) {
	db := s.DB()
	return &param, db.Create(&param).Error
}

func (t *userRepo) Update(s repo.DBRepo, param model.User, value string, id string) (*model.User, error) {
	db := s.DB()
	user := model.User{}
	return &param, db.Model(&user).Where("id = ?", id).Updates(map[string]interface{}{"email": value}).Error
}

func (r *userRepo) GetAll(s repo.DBRepo) ([]model.User, error) {
	db := s.DB()
	var results []model.User
	return results, db.Table("user").Select("id").Find(&results).Error
}
