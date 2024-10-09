package log

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapSugaredLogger struct {
	base Logger
}

func (l *zapSugaredLogger) Named(name string) SugarLogger {
	return &zapSugaredLogger{base: l.base.Named(name)}
}

func (l *zapSugaredLogger) With(keysAndValues ...interface{}) SugarLogger {
	return &zapSugaredLogger{base: l.base.With(l.sweetenFields(context.Background(), keysAndValues)...)}
}

func (l *zapSugaredLogger) Set(keysAndValues ...interface{}) {
	l.base = l.base.With(l.sweetenFields(context.Background(), keysAndValues)...)
}

func (l *zapSugaredLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.base.Debug(ctx, msg, fields...)
}

func (l *zapSugaredLogger) Info(ctx context.Context, msg string, fields ...Field) {
	l.base.Info(ctx, msg, fields...)
}

func (l *zapSugaredLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.base.Warn(ctx, msg, fields...)
}

func (l *zapSugaredLogger) Error(ctx context.Context, msg string, fields ...Field) {
	l.base.Error(ctx, msg, fields...)
}

func (l *zapSugaredLogger) Panic(ctx context.Context, msg string, fields ...Field) {
	l.base.Panic(ctx, msg, fields...)
}

func (l *zapSugaredLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.base.Fatal(ctx, msg, fields...)
}

func (l *zapSugaredLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.base.Debugf(ctx, format, args...)
}

func (l *zapSugaredLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.base.Infof(ctx, format, args...)
}

func (l *zapSugaredLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.base.Warnf(ctx, format, args...)
}

func (l *zapSugaredLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.base.Errorf(ctx, format, args...)
}

func (l *zapSugaredLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.base.Panicf(ctx, format, args...)
}

func (l *zapSugaredLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.base.Fatalf(ctx, format, args...)
}

func (l *zapSugaredLogger) Debugw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.base.Debug(ctx, msg, l.sweetenFields(ctx, keysAndValues)...)
}

func (l *zapSugaredLogger) Infow(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.base.Info(ctx, msg, l.sweetenFields(ctx, keysAndValues)...)
}

func (l *zapSugaredLogger) Warnw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.base.Warn(ctx, msg, l.sweetenFields(ctx, keysAndValues)...)
}

func (l *zapSugaredLogger) Errorw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.base.Error(ctx, msg, l.sweetenFields(ctx, keysAndValues)...)
}

func (l *zapSugaredLogger) Panicw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.base.Panic(ctx, msg, l.sweetenFields(ctx, keysAndValues)...)
}

func (l *zapSugaredLogger) Fatalw(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.base.Fatal(ctx, msg, l.sweetenFields(ctx, keysAndValues)...)
}

func (l *zapSugaredLogger) SetLevel(lvl string) error {
	return l.base.SetLevel(lvl)
}

func (l *zapSugaredLogger) Sync() {
	l.base.Sync()
}

func (l *zapSugaredLogger) sweetenFields(ctx context.Context, args []interface{}) []Field {
	if len(args) == 0 {
		return nil
	}

	// Allocate enough space for the worst case; if users pass only structured
	// fields, we shouldn't penalize them with extra allocations.
	fields := make([]Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			l.base.Error(ctx, "Ignored key without a value.",
				Any("ignored", args[i]))
			break
		}

		// Consume this value and the next, treating them as a key-value pair.
		// If the key isn't a string, add this pair to the slice of invalid
		// pairs.
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			// Subsequent errors are likely, so allocate once up front.
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{
				key:      key,
				value:    val,
				position: i,
			})
		} else {
			fields = append(fields, Any(keyStr, val))
		}
		i += 2
	}

	// If we encountered any invalid key-value pairs, log an error.
	if len(invalid) > 0 {
		l.base.Error(ctx, "Ignored key-value pairs with non-string keys.",
			Array("invalid", invalid))
	}

	return fields
}

type invalidPair struct {
	key, value interface{}
	position   int
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("position", int64(p.position))
	zap.Any("key", p.key).AddTo(enc)
	zap.Any("value", p.value).AddTo(enc)
	return nil
}

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var errs []error
	for i := range ps {
		if err := enc.AppendObject(ps[i]); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
