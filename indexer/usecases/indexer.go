package usecase

type IUsecase interface {
	Process(log map[string]interface{}) error
	ShouldProcessLog(log map[string]interface{}) bool
	Name() string
}
