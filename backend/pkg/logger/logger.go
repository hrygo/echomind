package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// New creates a new zap logger with log rotation support
// production: if true, outputs JSON format and writes to file; if false, outputs console format
func New(production bool) (*zap.Logger, error) {
	var config zap.Config
	var core zapcore.Core

	if production {
		// Production: JSON format + File output with rotation
		config = zap.NewProductionConfig()

		// Configure log rotation
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "logs/backend.log",
			MaxSize:    100, // MB
			MaxBackups: 3,
			MaxAge:     7, // days
			Compress:   true,
		})

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			w,
			config.Level,
		)
	} else {
		// Development: Console format + Stdout
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		encoderConfig := config.EncoderConfig
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			config.Level,
		)
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}
