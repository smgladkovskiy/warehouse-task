package log

import (
	"context"
)

//go:generate mockgen -source=wrapper.go -destination=./wrapper_mock.go -package=log -mock_names WithLoggerable=WrapperMock
type WithLoggerable interface {
	Logger
	SetLogger(logger Logger)
}

// WithLogger is a wrapper for Logger
type WithLogger struct {
	logger Logger
}

const defaultLoggerName = "unnamedLogger"

var _ WithLoggerable = (*WithLogger)(nil)

func (w *WithLogger) Logger() *WithLogger {
	if w.logger == nil {
		w.logger = Named(defaultLoggerName)
	}

	return w
}

func (w *WithLogger) SetLogger(logger Logger) {
	if logger == nil {
		logger = Named(defaultLoggerName)
	}

	w.logger = logger
}

func (w *WithLogger) Named(s string) Logger {
	if w.logger == nil {
		return Named(defaultLoggerName)
	}

	return w.logger.Named(s)
}

func (w *WithLogger) With(field ...Field) Logger {
	if w.logger == nil {
		w.logger = Named(defaultLoggerName)
	}

	return w.logger.With(field...)
}

func (w *WithLogger) Set(field ...Field) {
	if w.logger == nil {
		return
	}

	w.logger.Set(field...)
}

func (w *WithLogger) Debug(ctx context.Context, s string, field ...Field) {
	if w.logger == nil {
		return
	}

	w.logger.Debug(ctx, s, field...)
}

func (w *WithLogger) Info(ctx context.Context, s string, field ...Field) {
	if w.logger == nil {
		return
	}

	w.logger.Info(ctx, s, field...)
}

func (w *WithLogger) Warn(ctx context.Context, s string, field ...Field) {
	if w.logger == nil {
		return
	}

	w.logger.Warn(ctx, s, field...)
}

func (w *WithLogger) Error(ctx context.Context, s string, field ...Field) {
	if w.logger == nil {
		return
	}

	w.logger.Error(ctx, s, field...)
}

func (w *WithLogger) Panic(ctx context.Context, s string, field ...Field) {
	if w.logger == nil {
		return
	}

	w.logger.Panic(ctx, s, field...)
}

func (w *WithLogger) Fatal(ctx context.Context, s string, field ...Field) {
	if w.logger == nil {
		return
	}

	w.logger.Fatal(ctx, s, field...)
}

func (w *WithLogger) Debugf(ctx context.Context, s string, i ...interface{}) {
	if w.logger == nil {
		return
	}

	w.logger.Debugf(ctx, s, i...)
}

func (w *WithLogger) Infof(ctx context.Context, s string, i ...interface{}) {
	if w.logger == nil {
		return
	}

	w.logger.Infof(ctx, s, i...)
}

func (w *WithLogger) Warnf(ctx context.Context, s string, i ...interface{}) {
	if w.logger == nil {
		return
	}

	w.logger.Warnf(ctx, s, i...)
}

func (w *WithLogger) Errorf(ctx context.Context, s string, i ...interface{}) {
	if w.logger == nil {
		return
	}

	w.logger.Errorf(ctx, s, i...)
}

func (w *WithLogger) Panicf(ctx context.Context, s string, i ...interface{}) {
	if w.logger == nil {
		return
	}

	w.logger.Panicf(ctx, s, i...)
}

func (w *WithLogger) Fatalf(ctx context.Context, s string, i ...interface{}) {
	if w.logger == nil {
		return
	}

	w.logger.Fatalf(ctx, s, i...)
}

func (w *WithLogger) SetLevel(s string) error {
	if w.logger == nil {
		w.logger = Named(defaultLoggerName)
	}

	return w.logger.SetLevel(s)
}

func (w *WithLogger) Sugar() SugarLogger {
	if w.logger == nil {
		w.logger = Named(defaultLoggerName)
	}

	return w.logger.Sugar()
}

func (w *WithLogger) Sync() {
	if w.logger == nil {
		return
	}

	w.logger.Sync()
}
