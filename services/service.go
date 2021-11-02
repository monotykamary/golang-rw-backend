package services

import (
	"github.com/monotykamary/golang-rw-backend/config"
	"github.com/monotykamary/golang-rw-backend/repo"
	redisService "github.com/monotykamary/golang-rw-backend/services/redis"
)

type Services struct {
	Redis redisService.RedisService
}

func NewServices(cfg config.Config, store repo.DBRepo, repo *repo.Repo) *Services {
	return &Services{
		Redis: redisService.NewRedisService(cfg),
	}
}
