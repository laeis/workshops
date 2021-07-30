package response

import (
	"go.uber.org/zap"
)

type ErrorLogger struct {
	message string
	code    int
	payload map[string]interface{}
}

func (l ErrorLogger) WithLogger(logger *zap.Logger) {
	logger.Error(l.message, zap.Int("status_code", l.code))
}

type SuccessLogger struct {
	message string
	code    int
}

func (l SuccessLogger) WithLogger(logger *zap.Logger) {
	logger.Info(l.message, zap.Int("status_code", l.code))
}
