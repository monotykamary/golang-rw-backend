package repo

import (
	"gorm.io/gorm"
)

type DBRepo interface {
	DB() *gorm.DB
	NewTransaction() (DBRepo, FinallyFunc)
}

type FinallyFunc = func(error) error
