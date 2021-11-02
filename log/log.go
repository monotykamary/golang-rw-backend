package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() func() {
	var logger *zap.Logger
	switch os.Getenv("RUN_MODE") {
	case "local":
		atom := zap.NewAtomicLevel()
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger = zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		), zap.AddCaller())
	default:
		logger, _ = zap.NewProduction()
	}
	return zap.ReplaceGlobals(logger)
}
