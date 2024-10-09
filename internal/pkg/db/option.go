package db

import (
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
)

// Option specifies configuration options of Instance.
type Option func(*Instance)

// WithLogger is an option to set the logger for the Instance.
func WithLogger(logger log.Logger) Option {
	return func(i *Instance) {
		i.log = logger
	}
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If the option is not specified, tracing is disabled.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return func(i *Instance) {
		i.tracerProvider = provider
	}
}

func limitQuerySize(query string, limit int) string {
	maxQuerySize := len(query)
	if limit != 0 && maxQuerySize > limit {
		maxQuerySize = limit
	}

	return query[:maxQuerySize]
}

// WithGormConfig specifies an optional GORM config.
func WithGormConfig(cfg *gorm.Config) Option {
	return func(i *Instance) {
		i.gormConfig = cfg
	}
}

// PreferSimpleProtocol is applied for Postgres connections. It disables implicit
// prepared statement usage. By default pgx automatically uses the extended
// protocol. This can improve performance due to being able to use the binary
// format. It also does not rely on client side parameter sanitization. However,
// it does incur two round-trips per query (unless using a prepared statement)
// and may be incompatible proxies such as PGBouncer (in Transaction mode).
//
// From pgx documentation:
// In general QueryExecModeSimpleProtocol should only be used if connecting to a
// proxy server, connection pool server, or non-PostgreSQL server that does not
// support the extended protocol.
func PreferSimpleProtocol() Option {
	return func(i *Instance) {
		i.preferSimpleProtocol = true
	}
}
