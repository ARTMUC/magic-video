package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init() {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	//config.DisableCaller = true
	config.DisableStacktrace = true
	logger, _ := config.Build()
	Log = logger
}
