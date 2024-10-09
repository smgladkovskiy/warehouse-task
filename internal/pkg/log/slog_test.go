package log

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type loggedEntry struct {
	level   zapcore.Level
	message string
	fields  map[string]any
}

func TestHandler(t *testing.T) {
	type dummy struct {
		name string
	}

	tests := []struct {
		name string
		run  func(*slog.Logger)
		want []loggedEntry
	}{
		{
			name: "debug no fields",
			run: func(l *slog.Logger) {
				l.Debug("debug")
			},
			want: []loggedEntry{{
				level:   zapcore.DebugLevel,
				message: "debug",
				fields:  map[string]any{},
			}},
		},
		{
			name: "debug with fields",
			run: func(l *slog.Logger) {
				l.Debug("debug",
					"foo", "bar",
					"id", 42,
					"n", uint64(100),
					"pi", 3.14,
					"flag", true,
					"time", time.UnixMicro(1024),
					"duration", time.Second,
					"dummy", dummy{name: "abacaba"},
					"group", slog.GroupValue(slog.String("a", "b")),
				)
			},
			want: []loggedEntry{{
				level:   zapcore.DebugLevel,
				message: "debug",
				fields: map[string]any{
					"foo":      "bar",
					"id":       int64(42),
					"n":        uint64(100),
					"pi":       float64(3.14),
					"flag":     true,
					"time":     time.UnixMicro(1024),
					"duration": time.Second,
					"dummy":    dummy{name: "abacaba"},
					"group":    map[string]any{"a": "b"},
				},
			}},
		},
		{
			name: "info",
			run: func(l *slog.Logger) {
				l.Info("info",
					slog.String("foo", "bar"),
					slog.Group("group",
						slog.String("a", "b")))
			},
			want: []loggedEntry{{
				level:   zapcore.InfoLevel,
				message: "info",
				fields: map[string]any{
					"foo":   "bar",
					"group": map[string]any{"a": "b"},
				},
			}},
		},
		{
			name: "info with group",
			run: func(l *slog.Logger) {
				l.With(slog.String("foo", "bar")).
					WithGroup("group").
					Info("info", slog.String("foo", "qux"))
			},
			want: []loggedEntry{{
				level:   zapcore.InfoLevel,
				message: "info",
				fields: map[string]any{
					"foo":   "bar",
					"group": map[string]any{"foo": "qux"},
				},
			}},
		},
		{
			name: "info with group fields",
			run: func(l *slog.Logger) {
				l.Info("info", slog.String("a", "b"),
					slog.Group("user",
						slog.Int("id", 42),
						slog.String("name", "alice"),
						slog.Group("pet", slog.String("name", "kitty"))))
			},
			want: []loggedEntry{{
				level:   zapcore.InfoLevel,
				message: "info",
				fields: map[string]any{
					"a": "b",
					"user": map[string]any{
						"id":   int64(42),
						"name": "alice",
						"pet": map[string]any{
							"name": "kitty",
						},
					},
				},
			}},
		},
		{
			name: "warn",
			run: func(l *slog.Logger) {
				l.Warn("warn", "hello", "world")
			},
			want: []loggedEntry{{
				level:   zapcore.WarnLevel,
				message: "warn",
				fields:  map[string]any{"hello": "world"},
			}},
		},
		{
			name: "error",
			run: func(l *slog.Logger) {
				l.Error("error", "operation", "read", "error", io.EOF)
			},
			want: []loggedEntry{{
				level:   zapcore.ErrorLevel,
				message: "error",
				fields: map[string]any{
					"operation": "read",
					"error":     io.EOF.Error(),
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zc, logs := observer.New(zap.DebugLevel)
			zapLogger := NewZapLogger(zap.New(zc))
			_ = zapLogger.SetLevel("debug")

			logger := slog.New(NewHandler(zapLogger))
			tt.run(logger)

			for i, entry := range logs.All() {
				assert.Equal(t, tt.want[i].level, entry.Level)
				assert.Equal(t, tt.want[i].message, entry.Message)
				assert.Equal(t, tt.want[i].fields, entry.ContextMap())
			}
		})
	}
}
