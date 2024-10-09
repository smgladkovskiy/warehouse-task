// TODO: сделать сравнение вывода логов.
package log

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log/internal/testdata"
)

type testStruct struct {
	Name string
	Age  int
}

func (s testStruct) String() string {
	return fmt.Sprintf("%s: %d", s.Name, s.Age)
}

func bodyLog(l Logger) {
	str := testStruct{Name: "Bill", Age: 76}
	now := time.Now()
	protoMessage := &testdata.Request{Name: "Alice", Age: 10}

	fields := []Field{
		String("string", "test value"),
		Bool("bool", true),
		Binary("binary", []byte{1, 0, 0, 23}),
		Float64("float64", 3.14),
		Int("int", 2323),
		Int64("int64", 2222222222),
		Reflect("reflect", str),
		Stringer("stringer", str),
		Time("time", now),
		Timep("timep", &now),
		Duration("duration", time.Second*60),
		Any("any", str),
		Proto("proto", protoMessage),
	}

	type traceTestKey struct{}
	ctx := context.WithValue(context.Background(), traceTestKey{}, "from test")

	l.Debug(ctx, "Debug msg", fields...)
	l.Info(ctx, "Info msg", fields...)
	l.Warn(ctx, "Warn msg", fields...)
	l.Error(ctx, "Info msg", fields...)
	l.Info(ctx, "Info msg", fields...)

	l.Debugf(ctx, "Debug %s", "msg")
	l.Infof(ctx, "Debug %s", "msg")
	l.Warnf(ctx, "Debug %s", "msg")
	l.Errorf(ctx, "Debug %s", "msg")

	nl := l.With(String("WITH KEY", "WITH VALUE"))

	nl.Debug(ctx, "Debug msg", fields...)
	nl.Info(ctx, "Info msg", fields...)
	nl.Warn(ctx, "Warn msg", fields...)
	nl.Error(ctx, "Info msg", fields...)
	nl.Info(ctx, "Info msg", fields...)

	nl.Debugf(ctx, "Debug %s", "msg")
	nl.Infof(ctx, "Debug %s", "msg")
	nl.Warnf(ctx, "Debug %s", "msg")
	nl.Errorf(ctx, "Debug %s", "msg")

	Debug(ctx, "Debug msg", fields...)
	Info(ctx, "Info msg", fields...)
	Warn(ctx, "Warn msg", fields...)
	Error(ctx, "Info msg", fields...)
	Info(ctx, "Info msg", fields...)

	Debugf(ctx, "Debug %s", "msg")
	Infof(ctx, "Debug %s", "msg")
	Warnf(ctx, "Debug %s", "msg")
	Errorf(ctx, "Debug %s", "msg")
}

func bodySugarLog(l SugarLogger) {
	str := testStruct{Name: "Bill", Age: 76}
	now := time.Now()
	protoMessage := &testdata.Request{Name: "Alice", Age: 10}

	fields := []Field{
		String("string", "test value"),
		Bool("bool", true),
		Binary("binary", []byte{1, 0, 0, 23}),
		Float64("float64", 3.14),
		Int("int", 2323),
		Int64("int64", 2222222222),
		Reflect("reflect", str),
		Stringer("stringer", str),
		Time("time", now),
		Timep("timep", &now),
		Duration("duration", time.Second*60),
		Any("any", str),
		Proto("proto", protoMessage),
	}

	keysAndValues := []interface{}{
		"stringw", "test value",
		"boolw", true,
		"binaryw", []byte{1, 0, 0, 23},
		"float64w", 3.14,
		"intw", 2323,
		"int64w", 2222222222,
		"reflectw", str,
		"stringerw", str,
		"timew", now,
		"timepw", &now,
		"durationw", time.Second * 60,
		"anyw", str,
		"protow", protoMessage,
	}

	for _, field := range fields {
		keysAndValues = append(keysAndValues, field)
	}

	type traceTestKey struct{}
	ctx := context.WithValue(context.Background(), traceTestKey{}, "from test")

	l.Debug(ctx, "Debug msg", fields...)
	l.Info(ctx, "Info msg", fields...)
	l.Warn(ctx, "Warn msg", fields...)
	l.Error(ctx, "Info msg", fields...)
	l.Info(ctx, "Info msg", fields...)

	l.Debugf(ctx, "Debug %s", "msg")
	l.Infof(ctx, "Debug %s", "msg")
	l.Warnf(ctx, "Debug %s", "msg")
	l.Errorf(ctx, "Debug %s", "msg")

	nl := l.With(String("WITH KEY", "WITH VALUE"))

	nl.Debug(ctx, "Debug msg", fields...)
	nl.Info(ctx, "Info msg", fields...)
	nl.Warn(ctx, "Warn msg", fields...)
	nl.Error(ctx, "Info msg", fields...)
	nl.Info(ctx, "Info msg", fields...)

	nl.Debugw(ctx, "Debug msg", keysAndValues...)
	nl.Infow(ctx, "Info msg", keysAndValues...)
	nl.Warnw(ctx, "Warn msg", keysAndValues...)
	nl.Errorw(ctx, "Info msg", keysAndValues...)
	nl.Infow(ctx, "Info msg", keysAndValues...)

	nl.Debugf(ctx, "Debug %s", "msg")
	nl.Infof(ctx, "Debug %s", "msg")
	nl.Warnf(ctx, "Debug %s", "msg")
	nl.Errorf(ctx, "Debug %s", "msg")
}

func TestLoggerJson(_ *testing.T) {
	l := GetLogger(ConfigLogger{EncodingType: "json"})
	defer l.Sync()
	bodyLog(l)
	bodySugarLog(l.Sugar())
}

func TestLoggerConsole(_ *testing.T) {
	l := GetLogger(ConfigLogger{})
	defer l.Sync()
	bodyLog(l)
	bodySugarLog(l.Sugar())
}

func TestLoggerConsoleDebug(_ *testing.T) {
	l := GetLogger(ConfigLogger{Level: "DEBUG"})
	defer l.Sync()
	bodyLog(l)
	bodySugarLog(l.Sugar())
}

func TestLoggerFatal(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
		}
	}()
	l := GetLogger(ConfigLogger{})
	defer l.Sync()
	l.Panicf(context.TODO(), "fatal!!! %s", "error")
}

func TestInitLogs(_ *testing.T) {
	Info(context.TODO(), "test log console")
	GetLogger(ConfigLogger{EncodingType: "json"})
	Info(context.TODO(), "test log json")
}

func TestLoggerNamed(_ *testing.T) {
	l := GetLogger(ConfigLogger{EncodingType: "json"}).Named("test")
	l.Info(context.TODO(), "named logger", String("foo", "bar"))
	l.Sync()
}

func TestSugarLoggerNamed(_ *testing.T) {
	l := GetLogger(ConfigLogger{EncodingType: "json"}).Sugar().Named("test")
	l.Infow(context.TODO(), "named logger", "foo", "bar")
	l.Sync()
}

func TestDefaultSugarLogger(_ *testing.T) {
	l := Sugar()
	l.Infow(context.TODO(), "default sugar logger", "foo", "bar")
	l.Sync()
}

func TestFieldsWithInternalKeys(t *testing.T) {
	ctx := context.Background()

	const (
		msg     = "some msg"
		warnFmt = "Attempting to set a log field with a keys %v that is used internally, skipping."
	)

	tests := map[string]struct {
		fields                []Field
		excludedKeys, expKeys []string
		expWarn               bool
	}{
		"ok": {
			fields:       []Field{String("key", "value")},
			excludedKeys: []string{},
			expKeys:      []string{"key"},
			expWarn:      false,
		},
		timestampKey: {
			fields:       []Field{Time(timestampKey, time.Now())},
			excludedKeys: []string{timestampKey},
			expKeys:      []string{},
			expWarn:      true,
		},
		messageKey: {
			fields:       []Field{String(messageKey, msg), Bool("key", true)},
			excludedKeys: []string{messageKey},
			expKeys:      []string{"key"},
			expWarn:      true,
		},
		spanKey: {
			fields:       []Field{Int("int_key", 10), String(spanKey, spanKey), Bool("key", true)},
			excludedKeys: []string{spanKey},
			expKeys:      []string{"int_key", "key"},
			expWarn:      true,
		},
		traceKey: {
			fields:       []Field{Int("int_key", 10), String("key", ""), String(traceKey, traceKey)},
			excludedKeys: []string{traceKey},
			expKeys:      []string{"int_key", "key"},
			expWarn:      true,
		},
		levelKey: {
			fields:       []Field{Int("int_key", 10), String(levelKey, "info"), String(spanKey, spanKey)},
			excludedKeys: []string{levelKey, spanKey},
			expKeys:      []string{"int_key"},
			expWarn:      true,
		},
	}

	check := func(t *testing.T, entries []observer.LoggedEntry, expWarn bool, excludedKeys, expKeys []string) {
		require.NotEmpty(t, entries)

		if expWarn {
			assert.Equal(t, fmt.Sprintf(warnFmt, excludedKeys), entries[0].Message)
		} else {
			assert.NotEqual(t, fmt.Sprintf(warnFmt, excludedKeys), entries[0].Message)
		}

		fields := entries[len(entries)-1].ContextMap()
		for _, k := range expKeys {
			assert.Condition(t, func() bool {
				_, ok := fields[k]
				return ok
			}, "Absent key: '%s'", k)
		}
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			withLogger(func(l Logger, ol *observer.ObservedLogs) {
				l.Debug(ctx, msg, tc.fields...)
				check(t, ol.TakeAll(), tc.expWarn, tc.excludedKeys, tc.expKeys)

				l.Info(ctx, msg, tc.fields...)
				check(t, ol.TakeAll(), tc.expWarn, tc.excludedKeys, tc.expKeys)

				l.Warn(ctx, msg, tc.fields...)
				check(t, ol.TakeAll(), tc.expWarn, tc.excludedKeys, tc.expKeys)

				l.Error(ctx, msg, tc.fields...)
				check(t, ol.TakeAll(), tc.expWarn, tc.excludedKeys, tc.expKeys)

				l.With(tc.fields...).Info(ctx, msg)
				check(t, ol.TakeAll(), tc.expWarn, tc.excludedKeys, tc.expKeys)

				l.Set(tc.fields...)
				l.Info(ctx, msg)
				check(t, ol.TakeAll(), tc.expWarn, tc.excludedKeys, tc.expKeys)
			})
		})
	}
}

func withLogger(f func(Logger, *observer.ObservedLogs)) {
	zc, logs := observer.New(zap.DebugLevel)
	zl := zap.New(zc)

	logger := NewZapLogger(zl)
	_ = logger.SetLevel("debug")

	old := SetLogger(logger)
	defer SetLogger(old)

	f(GetLoggerInstance(), logs)
}

func TestOnFatal(t *testing.T) {
	l := GetLogger(ConfigLogger{})
	defer l.Sync()
	ctx := context.Background()
	l = OnFatal(CheckWriteAction(zapcore.WriteThenPanic))

	defer func() {
		require.NotNil(t, recover())
	}()

	l.Fatal(ctx, "Test")
}

func BenchmarkLog(b *testing.B) {
	l := GetLogger(ConfigLogger{
		Level:        "error",
		EncodingType: "json",
	})
	old := SetLogger(l)
	defer SetLogger(old)

	exp := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exp))
	otel.SetTracerProvider(tp)

	ctx, span := otel.Tracer("test").Start(context.Background(), "test")
	defer span.End()

	b.Run("info", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Info(ctx, "message", String("test", "test"))
		}
	})

	b.Run("infof", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Infof(ctx, "message: %s", "test")
		}
	})

	b.Run("with", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			With(String("test", "test")).Info(ctx, "Message")
		}
	})
}
