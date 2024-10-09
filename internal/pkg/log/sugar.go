package log

import (
	"context"
)

// SugarLogger wraps the base Logger functionality in a slower, but less
// verbose, API.
//
//nolint:interfacebloat // logger requires many methods
type SugarLogger interface {
	Named(string) SugarLogger
	With(...interface{}) SugarLogger
	Set(...interface{})

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

	Debugw(context.Context, string, ...interface{})
	Infow(context.Context, string, ...interface{})
	Warnw(context.Context, string, ...interface{})
	Errorw(context.Context, string, ...interface{})
	Panicw(context.Context, string, ...interface{})
	Fatalw(context.Context, string, ...interface{})

	SetLevel(string) error
	Sync()
}
