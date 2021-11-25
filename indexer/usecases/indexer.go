package usecase

type IUsecase interface {
	Process(log *RedisEvent) error
	ShouldProcessLog(log *RedisEvent) bool
	GetStreamInfo() (string, string)
	Name() string
}
