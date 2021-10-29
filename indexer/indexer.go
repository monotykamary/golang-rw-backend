package indexer

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/monotykamary/golang-rw-backend/config"
	usecase "github.com/monotykamary/golang-rw-backend/indexer/usecases"
	"github.com/monotykamary/golang-rw-backend/repo"
	"github.com/monotykamary/golang-rw-backend/repo/pg"
	"go.uber.org/zap"
)

type IndexService struct {
	cfg      config.Config
	store    repo.DBRepo
	repo     *repo.Repo
	usecases []usecase.IUsecase
}

func NewIndexService(cfg config.Config, store repo.DBRepo, client redis.Client) (*IndexService, error) {
	repo := pg.NewRepo()
	usecases := make([]usecase.IUsecase, 0)
	usecases = append(usecases, usecase.NewBookingUsecase(cfg, store, &client))

	return &IndexService{cfg, store, repo, usecases}, nil
}

func (svc *IndexService) Index() {
	stream, err := ConsumeStream(svc.cfg, "bookingGroup", "node")

	if err != nil {
		zap.L().Error("cannot read stream", zap.Error(err))
	}

OUTER:
	for _, item := range stream {
		for _, message := range item.Messages {
			log := message.Values
			for _, usecase := range svc.usecases {
				if usecase.ShouldProcessLog(log) {
					err := usecase.Process(log)
					if err != nil {
						zap.L().Fatal("panic when process block", zap.Any("debug data", log))
					}
					continue OUTER
				}
			}
		}
	}
}

func ConsumeStream(c config.Config, consumerGroup string, uniqueID string) ([]redis.XStream, error) {
	ctx := context.Background()
	addr := c.RedisHost + ":" + c.RedisPort
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: c.RedisPass,
		DB:       0,
	})

	err := redisClient.XGroupCreate(ctx, "booking", consumerGroup, "0").Err()
	if err != nil {
		fmt.Println(err)
	}

	return redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: uniqueID,
		Streams:  []string{"booking", ">"},
		Count:    10,
		Block:    0,
	}).Result()
}
