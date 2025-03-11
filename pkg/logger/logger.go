package logger

import (
	"context"
	"log"
	"sync"

	"github.com/0sokrat0/BookAPI/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger — обертка для zap.SugaredLogger
type Logger struct {
	*zap.SugaredLogger
}

var (
	instance *Logger
	once     sync.Once
)

type ctxKey string

const loggerCtxKey ctxKey = "logger"

// NewLogger создает логгер в зависимости от режима (prod/dev)
func NewLogger(cfg *config.Config) *Logger {
	once.Do(func() {
		var l *zap.Logger
		var err error

		if cfg.Logger.Level == "production" {
			l, err = zap.NewProduction()
		} else {
			// Настраиваем цветной вывод для VS Code
			config := zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Цвета уровней логов
			config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder      // Короткий путь к файлу
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Человеческое время
			config.DisableStacktrace = true                                     // Отключаем stacktrace для INFO и WARN

			// Добавляем цветной вывод в терминал
			config.EncoderConfig.ConsoleSeparator = " | "
			config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

			l, err = config.Build(zap.WithCaller(true)) // Включаем caller (пути файлов)
		}

		if err != nil {
			panic(err)
		}

		instance = &Logger{l.Sugar()}
	})
	return instance
}

// Sync вызывает flush буфера логов
func (l *Logger) Sync() {
	_ = l.SugaredLogger.Sync()
}

func FromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(loggerCtxKey).(*Logger)
	if !ok || l == nil {
		log.Fatal("Логгер не найден в контексте. Проверьте, что вы передаёте контекст, содержащий логгер")
	}
	return l
}

// WithLogger добавляет логгер в контекст
func WithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}
