package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/monotykamary/golang-rw-backend/config"
	usecase "github.com/monotykamary/golang-rw-backend/indexer/usecases"
)

type RedisService interface {
	XAddBooking(request BookingQueueRequest) (string, error)
}

type redisService struct {
	redisClient redis.Client
}

func NewRedisService(cfg config.Config) RedisService {
	addr := cfg.RedisHost + ":" + cfg.RedisPort
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	return &redisService{
		redisClient: *redisClient,
	}
}

func (r *redisService) XAddBooking(request BookingQueueRequest) (string, error) {
	ctx := context.Background()

	event := request.Event
	userId := request.UserId
	bookingId := request.BookingId

	var redisEventInterface map[string]interface{}
	redisEventStruct := &usecase.RedisEvent{
		Event:     event,
		UserId:    userId,
		BookingId: bookingId,
	}
	redisEventJSON, _ := json.Marshal(redisEventStruct)
	json.Unmarshal(redisEventJSON, &redisEventInterface)

	result, err := r.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "booking",
		Values: redisEventInterface,
	}).Result()

	return result, err
}
