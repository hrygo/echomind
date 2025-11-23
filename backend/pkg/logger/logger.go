package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(production bool) (*zap.Logger, error) {
	var config zap.Config
	if production {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	return config.Build()
}
