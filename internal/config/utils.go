package config

import (
	"strconv"
	"time"

	"github.com/ARTMUC/magic-video/internal/logger"
	"go.uber.org/zap"
)

// toTimeDuration transforms type from string (seconds) to time.Duration with default value (seconds)
func toTimeDuration(e string, def int) time.Duration {
	eInt, err := strconv.Atoi(e)
	if err != nil {
		logger.Log.Error("Error converting env to int", zap.Error(err))
		return time.Duration(def) * time.Second
	}
	return time.Duration(eInt) * time.Second
}

func toInt(e string, def int) int {
	eInt, err := strconv.Atoi(e)
	if err != nil {
		logger.Log.Error("Error converting env to int", zap.Error(err))
		return def
	}
	return eInt
}
