package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/monotykamary/golang-rw-backend/config"
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

	result, err := r.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "booking",
		Values: map[string]interface{}{
			"event":     event,
			"userId":    userId,
			"bookingId": bookingId,
		},
	}).Result()

	return result, err
}
