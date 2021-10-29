package pg

import (
	"fmt"
	"net/http"

	"github.com/monotykamary/golang-rw-backend/config"
	"github.com/monotykamary/golang-rw-backend/model/errors"
	"github.com/monotykamary/golang-rw-backend/repo"
	"github.com/monotykamary/golang-rw-backend/util"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

// store is implimentation of repository
type store struct {
	database *gorm.DB
}

// DB database connection
func (s *store) DB() *gorm.DB {
	return s.database
}

// NewTransaction for database connection
// return an db instance and a done(error) function, if the error != nil -> rollback
func (s *store) NewTransaction() (newRepo repo.DBRepo, finallyFn repo.FinallyFunc) {
	newDB := s.database.Begin()
	finallyFn = func(err error) error {
		if err != nil {
			nErr := newDB.Rollback().Error
			if nErr != nil {
				return errors.NewStringError(nErr.Error(), http.StatusInternalServerError)
			}
			return errors.NewStringError(err.Error(), util.ParseErrorCode(err))
		}

		cErr := newDB.Commit().Error
		if cErr != nil {
			return errors.NewStringError(cErr.Error(), http.StatusInternalServerError)
		}
		return nil
	}

	return &store{database: newDB}, finallyFn
}

func NewPostgresStore(cfg *config.Config) (repo.DBRepo, func() error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=disable",
			cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBHost, cfg.DBPort),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: glog.Default.LogMode(glog.Info),
	})
	if err != nil {
		zap.L().Panic("cannot connect to db", zap.Error(err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Panic("cannot connect to db", zap.Error(err))
	}
	return &store{database: db}, sqlDB.Close
}
