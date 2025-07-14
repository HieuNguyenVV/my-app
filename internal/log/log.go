package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"nvh/my-app/internal/pkg/config"
)

var SemanticVersion string

type Logger interface {
	// Basic loggig methods
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	// Format methods
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	// Line methods
	Println(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	// Field methods
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger

	Sync() error
}

type Fields map[string]interface{}

type ZapAdapter struct {
	logger *zap.SugaredLogger
}

func NewLogger(conf config.Config) (Logger, error) {
	zapConfig := zap.Config{}

	if conf.AppConfig.Debug {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	if conf.Log.Format != "" {
		zapConfig.Encoding = conf.Log.Format
	} else {
		zapConfig.Encoding = "console"
	}

	zapConfig.EncoderConfig = defaultEncoderConfig()
	level := parseLogLevel(conf.Log.Level)
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("logger.NewLogger: failed to build zap logger: %w", err)
	}
	if SemanticVersion != "" {
		logger = logger.With(zap.String("version", SemanticVersion))
	}
	return &ZapAdapter{
		logger: logger.Sugar(),
	}, nil
}

func defaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func (z *ZapAdapter) Debug(args ...interface{}) {
	z.logger.Debug(args...)
}

func (z *ZapAdapter) Info(args ...interface{}) {
	z.logger.Info(args...)
}

func (z *ZapAdapter) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}

func (z *ZapAdapter) Error(args ...interface{}) {
	z.logger.Error(args...)
}

func (z *ZapAdapter) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}

func (z *ZapAdapter) Panic(args ...interface{}) {
	z.logger.Panic(args...)
}

func (z *ZapAdapter) Debugf(format string, args ...interface{}) {
	z.logger.Debugf(format, args...)
}

func (z *ZapAdapter) Infof(format string, args ...interface{}) {
	z.logger.Infof(format, args...)
}

func (z *ZapAdapter) Warnf(format string, args ...interface{}) {
	z.logger.Warnf(format, args...)
}

func (z *ZapAdapter) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args...)
}

func (z *ZapAdapter) Fatalf(format string, args ...interface{}) {
	z.logger.Fatalf(format, args...)
}

func (z *ZapAdapter) Panicf(format string, args ...interface{}) {
	z.logger.Panicf(format, args...)
}

func (z *ZapAdapter) Println(args ...interface{}) {
	z.logger.Info(fmt.Sprintln(args...))
}

func (z *ZapAdapter) Debugln(args ...interface{}) {
	z.logger.Debug(fmt.Sprintln(args...))
}

func (z *ZapAdapter) Infoln(args ...interface{}) {
	z.logger.Info(fmt.Sprintln(args...))
}

func (z *ZapAdapter) Warnln(args ...interface{}) {
	z.logger.Warn(fmt.Sprintln(args...))
}

func (z *ZapAdapter) Errorln(args ...interface{}) {
	z.logger.Error(fmt.Sprintln(args...))
}

func (z *ZapAdapter) Fatalln(args ...interface{}) {
	z.logger.Fatal(fmt.Sprintln(args...))
}

func (z *ZapAdapter) Panicln(args ...interface{}) {
	z.logger.Panic(fmt.Sprintln(args...))
}

func (z *ZapAdapter) WithField(key string, value interface{}) Logger {
	return &ZapAdapter{
		logger: z.logger.With(key, value),
	}
}
func (z *ZapAdapter) WithFields(fields Fields) Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for key, value := range fields {
		args = append(args, key, value)
	}

	return &ZapAdapter{
		logger: z.logger.With(args...),
	}
}

func (z *ZapAdapter) WithError(err error) Logger {
	return &ZapAdapter{
		logger: z.logger.With("error", err),
	}
}

func (z *ZapAdapter) Sync() error {
	return z.logger.Sync()
}
