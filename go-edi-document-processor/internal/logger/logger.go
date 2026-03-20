package logger

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger представляет собой обёртку над zap и slog
type Logger struct {
	zap  *zap.Logger
	slog *slog.Logger
}

// New создаёт новый логгер на основе конфигурации
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

	// Создаём slog логгер
	slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return &Logger{
		zap:  zapLogger,
		slog: slogLogger,
	}, nil
}

// Zap возвращает zap логгер
func (l *Logger) Zap() *zap.Logger {
	return l.zap
}

// Slog возвращает slog логгер
func (l *Logger) Slog() *slog.Logger {
	return l.slog
}

// Sync закрывает логгер
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

// Global логгер (синглтон)
var globalLogger *Logger

// InitGlobal инициализирует глобальный логгер
func InitGlobal(level string, development bool) error {
	logger, err := New(level, development)
	if err != nil {
		return err
	}
	globalLogger = logger
	// Устанавливаем глобальный slog логгер
	slog.SetDefault(logger.slog)
	return nil
}

// GetGlobal возвращает глобальный логгер
func GetGlobal() *Logger {
	if globalLogger == nil {
		// fallback на стандартный логгер
		zapLogger, _ := zap.NewProduction()
		slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		globalLogger = &Logger{
			zap:  zapLogger,
			slog: slogLogger,
		}
	}
	return globalLogger
}
