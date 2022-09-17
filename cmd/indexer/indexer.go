package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/monotykamary/golang-rw-backend/indexer"
	"github.com/monotykamary/golang-rw-backend/log"

	"github.com/monotykamary/golang-rw-backend/config"
	"github.com/monotykamary/golang-rw-backend/repo/pg"
	"go.uber.org/zap"
)

func main() {
	//init config
	cls := config.DefaultConfigLoaders()
	cfg := config.LoadConfig(cls)

	// init logger
	undo := log.New()
	defer zap.L().Sync()
	defer undo()

	s, close := pg.NewPostgresStore(&cfg)
	defer close()

	addr := cfg.RedisHost + ":" + cfg.RedisPort
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	idxSvc, err := indexer.NewIndexService(cfg, s, *redisClient)
	if err != nil {
		zap.L().Panic("cannot init indexer", zap.Error(err))
	}

	idxSvc.Index()
}
