package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/monotykamary/golang-rw-backend/config"
	"github.com/monotykamary/golang-rw-backend/model"
	"github.com/monotykamary/golang-rw-backend/repo"
	"github.com/monotykamary/golang-rw-backend/repo/pg"
	"github.com/qmuntal/stateless"
	"go.uber.org/zap"
)

const (
	stateBooking    = "Booking"
	stateProcessing = "Processing"
	stateCancelled  = "Cancelled"
	stateBooked     = "Booked"
)

const (
	triggerProcess = "Process"
	triggerPayment = "Payment"
	triggerCancel  = "Cancel"
	triggerRetry   = "Retry"
)

type RedisEvent struct {
	Event     string    `json:"event"`
	UserId    uuid.UUID `json:"userId"`
	BookingId uuid.UUID `json:"bookingId"`
}

type BookingUsecase struct {
	cfg    config.Config
	store  repo.DBRepo
	repo   *repo.Repo
	client *redis.Client
	stream string
	group  string
}

func NewBookingUsecase(cfg config.Config, store repo.DBRepo, client *redis.Client) *BookingUsecase {
	repo := pg.NewRepo()
	stream := "booking"
	group := "bookingGroup"

	bookingUsecase := &BookingUsecase{cfg, store, repo, client, stream, group}
	return bookingUsecase
}

func NewBookingStateMachine(initialState string) *stateless.StateMachine {
	booking := stateless.NewStateMachine(initialState)

	booking.Configure(stateBooking).
		Permit(triggerProcess, stateProcessing)

	booking.Configure(stateProcessing).
		Permit(triggerPayment, stateBooked).
		Permit(triggerCancel, stateCancelled)

	booking.Configure(stateCancelled).
		Permit(triggerRetry, stateBooking)

	return booking
}

func (uc BookingUsecase) Process(log *RedisEvent) error {
	db := uc.store.DB()

	var redisEventInterface map[string]interface{}
	redisEventJSON, _ := json.Marshal(log)
	json.Unmarshal(redisEventJSON, &redisEventInterface)

	id := redisEventInterface["bookingId"]
	booking, err := uc.repo.Booking.GetById(uc.store, fmt.Sprintf("%v", id))

	if err != nil {
		return db.Create(&model.Booking{}).Error
	}

	stateMachine := NewBookingStateMachine(booking.Status)
	stateMachine.Fire(redisEventInterface["event"])
	currentStatus := stateMachine.MustState()

	err = db.Model(booking).Update("status", currentStatus).Error
	if err != nil {
		zap.L().Error("cannot update booking", zap.Error(err))
	}

	return nil
}

func (uc BookingUsecase) ShouldProcessLog(log *RedisEvent) bool {
	// doesn't matter for our mock case
	return true
}

func (uc BookingUsecase) GetStreamInfo() (stream string, group string) {
	stream = uc.stream
	group = uc.group
	return
}

func (uc BookingUsecase) Name() string {
	return "BookingUsecase"
}
