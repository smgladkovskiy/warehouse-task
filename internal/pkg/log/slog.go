package log

import (
	"context"
	"log/slog"

	"go.uber.org/zap/zapcore"
)

// Handler is a slog.Handler that forwards logs to the Logger.
type Handler struct {
	log Logger
}

// NewHandler creates a new slog.Handler.
func NewHandler(logger Logger) Handler {
	return Handler{log: logger}
}

var _ slog.Handler = (*Handler)(nil)

// Enabled reports whether the handler handles records at the given level.
// Always returns true. Level is controlled by the logger.
func (h Handler) Enabled(context.Context, slog.Level) bool {
	return true
}

// Handle handles the Record.
func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	var handle func(context.Context, string, ...Field)

	switch record.Level {
	case slog.LevelDebug:
		handle = h.log.Debug
	case slog.LevelInfo:
		handle = h.log.Info
	case slog.LevelWarn:
		handle = h.log.Warn
	case slog.LevelError:
		handle = h.log.Error
	}

	fields := make([]Field, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		fields = append(fields, convertAttrToField(attr))
		return true
	})

	handle(ctx, record.Message, fields...)

	return nil
}

// WithAttrs returns a new Handler whose attributes consist of both the
// receiver's attributes and the arguments.
func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]Field, 0, len(attrs))
	for _, attr := range attrs {
		fields = append(fields, convertAttrToField(attr))
	}
	return Handler{h.log.With(fields...)}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h Handler) WithGroup(name string) slog.Handler {
	return Handler{h.log.With(Namespace(name))}
}

func convertAttrToField(attr slog.Attr) Field {
	if attr.Equal(slog.Attr{}) {
		// Ignore empty attrs.
		return Skip()
	}

	switch attr.Value.Kind() {
	case slog.KindBool:
		return Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		return Duration(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		return Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		return Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		return String(attr.Key, attr.Value.String())
	case slog.KindTime:
		return Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		return Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup:
		return group(attr.Key, groupObject(attr.Value.Group()))
	case slog.KindLogValuer:
		return convertAttrToField(slog.Attr{
			Key:   attr.Key,
			Value: attr.Value.Resolve(),
		})
	default:
		return Any(attr.Key, attr.Value.Any())
	}
}

// groupObject holds all the Attrs saved in a slog.GroupValue.
type groupObject []slog.Attr

func (gs groupObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, attr := range gs {
		zapcore.Field(convertAttrToField(attr)).AddTo(enc)
	}
	return nil
}
