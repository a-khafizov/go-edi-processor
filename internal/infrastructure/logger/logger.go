package logger

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap  *zap.Logger
	slog *slog.Logger
}

func New(level string, development bool) (*Logger, error) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	var config zap.Config
	if development {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return &Logger{
		zap:  zapLogger,
		slog: slogLogger,
	}, nil
}

func (l *Logger) Zap() *zap.Logger {
	return l.zap
}

func (l *Logger) Slog() *slog.Logger {
	return l.slog
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}

var globalLogger *Logger

func InitGlobal(level string, development bool) error {
	logger, err := New(level, development)
	if err != nil {
		return err
	}
	globalLogger = logger

	slog.SetDefault(logger.slog)
	return nil
}

func GetGlobal() *Logger {
	if globalLogger == nil {
		zapLogger, _ := zap.NewProduction()
		slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		globalLogger = &Logger{
			zap:  zapLogger,
			slog: slogLogger,
		}
	}
	return globalLogger
}
