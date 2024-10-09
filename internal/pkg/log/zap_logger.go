package log

import (
	"context"
	"errors"
	"fmt"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field keys used by the current package and zap.
const (
	timestampKey  = "timestamp"
	traceKey      = "trace_id"
	spanKey       = "span_id"
	levelKey      = "level"
	nameKey       = "logger"
	callerKey     = "caller"
	messageKey    = "msg"
	errorVerbose  = "errorVerbose"
	stacktraceKey = "stacktrace"
	componentKey  = "component"

	// TraceFieldKey is a key for trace field.
	TraceFieldKey = "_trace"

	// TraceMetricLevel is a value for metric label.
	TraceMetricLevel = "trace"
)

var _ Logger = (*zapLogger)(nil)

type zapLogger struct {
	config zap.Config
	logger *zap.Logger
	level  zap.AtomicLevel
	levels LevelRegistry
	name   string
}

func initZap(config ConfigLogger) (cfg zap.Config, logger *zap.Logger, err error) {
	var options []zap.Option

	cfg = zap.NewProductionConfig()

	cfg.EncoderConfig.NameKey = componentKey
	cfg.EncoderConfig.TimeKey = timestampKey
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	cfg.Level = zap.NewAtomicLevel()
	if err = cfg.Level.UnmarshalText([]byte(config.Level)); err != nil {
		return
	}

	// check encoding type
	if config.EncodingType == "console" {
		cfg.Encoding = config.EncodingType
	}

	if !config.EnableStacktrace {
		cfg.DisableStacktrace = true
		options = append(options, zap.AddStacktrace(zapcore.FatalLevel))
	}

	if config.EnableCaller {
		options = append([]zap.Option{zap.AddCallerSkip(2)}, options...)
	} else {
		cfg.DisableCaller = true
	}

	logger, err = cfg.Build(options...)

	return
}

// SetSentry sets sentry server to send errors.
func (s *zapLogger) SetSentry(client *sentry.Client) error {
	cfg := zapsentry.Configuration{
		Level: zapcore.ErrorLevel, // when to send message to sentry
		Tags: map[string]string{
			"component": "logs",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))
	// in case of err it will return noop core. so we can safely attach it
	if err != nil {
		return fmt.Errorf("creation of zap core for sentry is failed: %w", err)
	}
	s.logger = zapsentry.AttachCoreToLogger(core, s.logger)
	return nil
}

// LevelRegistry changes log level.
func (s *zapLogger) LevelRegistry() LevelRegistry {
	return s.levels
}

// SetLevel changes log level.
func (s *zapLogger) SetLevel(lvl string) error {
	value := new(zapcore.Level)
	if err := value.UnmarshalText([]byte(lvl)); err != nil {
		return err
	}

	s.level.SetLevel(*value)
	return nil
}

// OnFatal clones logger with different behaviour on fatal logs.
func (s *zapLogger) OnFatal(action CheckWriteAction) Logger {
	l := s.clone()
	l.logger = l.logger.WithOptions(zap.WithFatalHook(zapcore.CheckWriteAction(action)))

	return l
}

// Sync flushes any buffered log entries. Applications should take care to call
// Sync before exiting.
func (s *zapLogger) Sync() {
	// know errors
	// sync /dev/stdout: inappropriate ioctl for device
	// sync /dev/stdout: bad file descriptor
	// sync /dev/stderr: invalid argument
	knowErrors := []string{"inappropriate ioctl for device", "bad file descriptor", "invalid argument"}
	if err := s.logger.Sync(); err != nil {
		if baseError := errors.Unwrap(err); baseError != nil && !stringsContains(knowErrors, baseError.Error()) {
			s.Errorf(context.Background(), "sync log: %s", err)
		}
	}
}

func (s *zapLogger) clone() *zapLogger {
	cp := *s
	return &cp
}

func (s *zapLogger) Named(name string) Logger {
	if name == "" {
		return s
	}

	l := s.clone()

	if s.name == "" {
		l.name = name
	} else {
		l.name = l.name + "." + name
	}

	l.logger = l.logger.Named(name)
	l.level, _ = l.levels.GetOrSet(l.name, zap.NewAtomicLevelAt(s.level.Level()))

	return l
}

func (s *zapLogger) With(fields ...Field) Logger {
	l := s.clone()
	l.logger = l.logger.With(s.zapFieldHandle(context.Background(), fields)...)

	return l
}

func (s *zapLogger) Set(fields ...Field) {
	s.logger = s.logger.With(s.zapFieldHandle(context.Background(), fields)...)
}

func (s *zapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	if s.level.Enabled(zapcore.DebugLevel) {
		s.logger.Debug(msg, s.zapFieldHandle(ctx, fields)...)
	}
}

func (s *zapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	if s.level.Enabled(zapcore.InfoLevel) {
		s.logger.Info(msg, s.zapFieldHandle(ctx, fields)...)
	}
}

func (s *zapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	if s.level.Enabled(zapcore.WarnLevel) {
		s.logger.Warn(msg, s.zapFieldHandle(ctx, fields)...)
	}
}

func (s *zapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	if s.level.Enabled(zapcore.ErrorLevel) {
		s.logger.Error(msg, s.zapFieldHandle(ctx, fields)...)
	}
}

func (s *zapLogger) Panic(ctx context.Context, msg string, fields ...Field) {
	s.logger.Panic(msg, s.zapFieldHandle(ctx, fields)...)
}

func (s *zapLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	s.logger.Fatal(msg, s.zapFieldHandle(ctx, fields)...)
}

func (s *zapLogger) Debugf(ctx context.Context, format string, a ...interface{}) {
	if s.level.Enabled(zapcore.DebugLevel) {
		s.logger.Debug(fmt.Sprintf(format, a...), ctxToFields(ctx)...)
	}
}

func (s *zapLogger) Infof(ctx context.Context, format string, a ...interface{}) {
	if s.level.Enabled(zapcore.InfoLevel) {
		s.logger.Info(fmt.Sprintf(format, a...), ctxToFields(ctx)...)
	}
}

func (s *zapLogger) Warnf(ctx context.Context, format string, a ...interface{}) {
	if s.level.Enabled(zapcore.WarnLevel) {
		s.logger.Warn(fmt.Sprintf(format, a...), ctxToFields(ctx)...)
	}
}

func (s *zapLogger) Errorf(ctx context.Context, format string, a ...interface{}) {
	if s.level.Enabled(zapcore.ErrorLevel) {
		s.logger.Error(fmt.Sprintf(format, a...), ctxToFields(ctx)...)
	}
}

func (s *zapLogger) Panicf(ctx context.Context, format string, a ...interface{}) {
	s.logger.Panic(fmt.Sprintf(format, a...), ctxToFields(ctx)...)
}

func (s *zapLogger) Fatalf(ctx context.Context, format string, a ...interface{}) {
	s.logger.Fatal(fmt.Sprintf(format, a...), ctxToFields(ctx)...)
}

func (s *zapLogger) Sugar() SugarLogger {
	return &zapSugaredLogger{base: s}
}

func (s *zapLogger) warnInternalKey(ctx context.Context, keys []string) {
	s.Warnf(ctx, "Attempting to set a log field with a keys %v that is used internally, skipping.", keys)
}

// zapFieldHandle excludes fields with internal keys, adds trace fields, and
// converts the field to zap.Field.
func (s *zapLogger) zapFieldHandle(ctx context.Context, f []Field) []zap.Field {
	f, keys := filterInternalKeys(f)
	if len(keys) != 0 {
		s.warnInternalKey(ctx, keys)
	}

	c := ctxToFields(ctx)
	z := make([]zap.Field, 0, len(f)+len(c))
	z = append(z, c...)

	for _, i := range f {
		z = append(z, zap.Field(i))
	}

	return z
}

func isInternalKey(key string) bool {
	switch key {
	case timestampKey,
		traceKey,
		spanKey,
		levelKey,
		nameKey,
		callerKey,
		messageKey,
		errorVerbose,
		stacktraceKey:
		return true
	}
	return false
}

// filterInternalKeys modifies the provided fields to exclude those containing
// internal keys. It returns keys that were excluded.
func filterInternalKeys(fields []Field) ([]Field, []string) {
	exCount := 0
	for _, f := range fields {
		if isInternalKey(f.Key) {
			exCount++
		}
	}

	if exCount == 0 {
		return fields, nil
	}

	fs := make([]Field, 0, len(fields)-exCount)
	excludedKeys := make([]string, 0, exCount)

	for _, f := range fields {
		if isInternalKey(f.Key) {
			excludedKeys = append(excludedKeys, f.Key)
			continue
		}

		fs = append(fs, f)
	}

	return fs, excludedKeys
}

// stringsContains reports whether strings slice contains value.
func stringsContains(arr []string, value string) bool {
	for _, e := range arr {
		if value == e {
			return true
		}
	}

	return false
}

func ctxToFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		return []zap.Field{
			zap.String(traceKey, sc.TraceID().String()),
			zap.String(spanKey, sc.SpanID().String()),
		}
	}

	return nil
}
