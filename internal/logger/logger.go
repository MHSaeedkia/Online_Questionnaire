package logger

import (
	"fmt"
	config "online-questionnaire/configs"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger *Logger
	once         sync.Once
)

type Logger struct {
	zapLogger *zap.Logger
	service   string
	// endpoint  string
}

type Logctx struct {
	Data map[string]interface{}
}

func NewLogger(cfg config.Config, service string) error {
	var err error
	once.Do(func() {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   cfg.Logging.Filename,
			MaxSize:    cfg.Logging.MaxSize,
			MaxBackups: cfg.Logging.MaxBackups,
			MaxAge:     cfg.Logging.MaxAge,
			Compress:   cfg.Logging.Compress,
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:       "timestamp",
			LevelKey:      "level",
			MessageKey:    "message",
			NameKey:       "service",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		}

		level := zap.DebugLevel
		if cfg.Logging.Level != "" {
			if err = level.UnmarshalText([]byte(cfg.Logging.Level)); err != nil {
				err = fmt.Errorf("failed to parse log level: %w", err)
				return
			}
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(lumberjackLogger),
			level,
		)

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

		globalLogger = &Logger{
			zapLogger: zapLogger,
			service:   service,
			// endpoint:  endpoint,
		}
	})
	return err
}

func GetLogger() *Logger {
	return globalLogger
}
func (l *Logger) Debug(message, userID string, context Logctx) {
	l.log(zap.DebugLevel, message, userID, context, "")
}

func (l *Logger) Info(message, userID string, context Logctx) {
	l.log(zap.InfoLevel, message, userID, context, "")
}

func (l *Logger) Warning(message, userID string, context Logctx) {
	l.log(zap.WarnLevel, message, userID, context, "")
}

func (l *Logger) Error(message, userID string, context Logctx, traceID string) {
	l.log(zap.ErrorLevel, message, userID, context, traceID)
}

func (l *Logger) Fatal(message, userID string, context Logctx, traceID string) {
	l.log(zap.FatalLevel, message, userID, context, traceID)
	panic(1)
}

func (l *Logger) log(level zapcore.Level, message, userID string, context Logctx, traceID string) {
	fields := []zap.Field{
		zap.String("service", l.service),
		// zap.String("endpoint", l.endpoint),
		zap.String("user_id", userID),
		zap.String("trace_id", traceID),
		zap.Any("context", context.Data),
	}
	if ce := l.zapLogger.Check(level, message); ce != nil {
		ce.Write(fields...)
	}
}
