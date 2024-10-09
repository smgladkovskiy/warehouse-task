/*
Package log implements a logging package. Thin wrapper over zap.
It defines a type, Logger, with methods for formatting output.
*/
package log

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a logger interface.
//
//go:generate mockgen -source=log.go -destination=log_mock.go -package=log -mock_names Logger=LogMock
//nolint:interfacebloat // logger requires many methods
type Logger interface {
	Named(string) Logger
	With(...Field) Logger
	Set(...Field)

	Debug(context.Context, string, ...Field)
	Info(context.Context, string, ...Field)
	Warn(context.Context, string, ...Field)
	Error(context.Context, string, ...Field)
	Panic(context.Context, string, ...Field)
	Fatal(context.Context, string, ...Field)

	Debugf(context.Context, string, ...interface{})
	Infof(context.Context, string, ...interface{})
	Warnf(context.Context, string, ...interface{})
	Errorf(context.Context, string, ...interface{})
	Panicf(context.Context, string, ...interface{})
	Fatalf(context.Context, string, ...interface{})

	SetLevel(string) error

	Sugar() SugarLogger

	Sync()
}

var (
	once          sync.Once
	currentLogger unsafe.Pointer

	loggerConfig atomic.Pointer[ConfigLogger]

	zapLoggerOnce    sync.Once
	currentZapLogger *zap.Logger
)

func init() {
	cfg := ConfigLogger{Level: "DEBUG"}

	l, err := NewLogger(cfg)
	if err != nil {
		panic(fmt.Errorf("init logger: %w", err))
	}
	_ = SetLogger(l)

	if l, ok := l.(*zapLogger); ok {
		currentZapLogger = l.logger
	}
}

// GetLogger creates once with ConfigLogger and returns main application logger.
func GetLogger(config ConfigLogger) Logger {
	once.Do(func() {
		loggerConfig.Store(&config)

		if l := logger(); l != nil {
			l.Sync()
		}

		l, err := NewLogger(config)
		if err != nil {
			panic(err)
		}

		_ = SetLogger(l)
	})

	return logger()
}

// NewLogger creates new Logger.
func NewLogger(config ConfigLogger) (Logger, error) {
	var err error

	config.MergeDefault()
	logger := &zapLogger{
		level:  zap.NewAtomicLevel(),
		levels: NewLevelRegistry(),
	}

	logger.levels.Set(logger.name, logger.level)
	if err = logger.SetLevel(config.Level); err != nil {
		return nil, err
	}

	config.Level = "DEBUG"

	if logger.config, logger.logger, err = initZap(config); err != nil {
		return nil, err
	}

	return logger, nil
}

// GetLoggerInstance returns current Logger.
func GetLoggerInstance() Logger {
	return logger()
}

// NewZapLogger creates new zapLogger.
//
//nolint:revive // backward compatibility
func NewZapLogger(l *zap.Logger) *zapLogger {
	logger := &zapLogger{
		logger: l,
		level:  zap.NewAtomicLevel(),
		levels: NewLevelRegistry(),
	}

	logger.levels.Set(logger.name, logger.level)

	return logger
}

// SetLogger sets new logger as current.
func SetLogger(newLog Logger) Logger {
	if l := (*Logger)(atomic.SwapPointer(&currentLogger, unsafe.Pointer(&newLog))); l != nil {
		return *l
	}

	return nil
}

func logger() Logger {
	if l := (*Logger)(atomic.LoadPointer(&currentLogger)); l != nil {
		return *l
	}

	return nil
}

// SetSentry sets Sentry client to send errors to the Sentry server.
func SetSentry(client *sentry.Client) error {
	if h, ok := logger().(interface {
		SetSentry(*sentry.Client) error
	}); ok {
		return h.SetSentry(client)
	}
	return fmt.Errorf("sentry not implement for this logger")
}

// OnFatal clones logger with different behaviour on fatal logs.
func OnFatal(action CheckWriteAction) Logger {
	if h, ok := logger().(interface {
		OnFatal(CheckWriteAction) Logger
	}); ok {
		return h.OnFatal(action)
	}

	return logger()
}

// IsDebugLevel returns true if current logger level is DEBUG.
func IsDebugLevel() bool {
	if v, ok := logger().(*zapLogger); ok && v.level.Level() == zapcore.DebugLevel {
		return true
	}
	return false
}

// Named creates a child logger and adds a new path segment to the logger's
// name. Segments are joined by periods. By default, Loggers are unnamed.
func Named(name string) Logger {
	return logger().Named(name)
}

// Set sets structured context to it. Fields added to the current logger.
func Set(fields ...Field) {
	logger().Set(fields...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(fields ...Field) Logger {
	return logger().With(fields...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(ctx context.Context, msg string, fields ...Field) {
	logger().Debug(ctx, msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(ctx context.Context, msg string, fields ...Field) {
	logger().Info(ctx, msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(ctx context.Context, msg string, fields ...Field) {
	logger().Warn(ctx, msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(ctx context.Context, msg string, fields ...Field) {
	logger().Error(ctx, msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(ctx context.Context, msg string, fields ...Field) {
	logger().Panic(ctx, msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(ctx context.Context, msg string, fields ...Field) {
	logger().Fatal(ctx, msg, fields...)
}

// Debugf logs a message at DebugLevel.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(ctx context.Context, format string, a ...interface{}) {
	logger().Debugf(ctx, format, a...)
}

// Infof logs a message at InfoLevel.
// Arguments are handled in the manner of fmt.Printf.
func Infof(ctx context.Context, format string, a ...interface{}) {
	logger().Infof(ctx, format, a...)
}

// Warnf logs a message at WarnLevel.
// Arguments are handled in the manner of fmt.Printf.
func Warnf(ctx context.Context, format string, a ...interface{}) {
	logger().Warnf(ctx, format, a...)
}

// Errorf logs a message at ErrorLevel.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(ctx context.Context, format string, a ...interface{}) {
	logger().Errorf(ctx, format, a...)
}

// Panicf logs a message at PanicLevel.
// Arguments are handled in the manner of fmt.Printf.
func Panicf(ctx context.Context, format string, a ...interface{}) {
	logger().Panicf(ctx, format, a...)
}

// Fatalf logs a message at FatalLevel.
// Arguments are handled in the manner of fmt.Printf.
func Fatalf(ctx context.Context, format string, a ...interface{}) {
	logger().Fatalf(ctx, format, a...)
}

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower, API.
func Sugar() SugarLogger {
	return &zapSugaredLogger{base: logger()}
}

// Sync flushes any buffered log entries. Applications should take care to call
// Sync before exiting.
func Sync() {
	logger().Sync()
}
