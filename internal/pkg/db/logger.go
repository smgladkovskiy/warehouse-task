package db

import (
	"context"
	"strconv"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
)

// DefaultGormLoggerName is a default GORM logger name.
const DefaultGormLoggerName = "gorm"

var _ logger.Interface = (*gormLogger)(nil)

type gormLogger struct {
	log log.Logger
}

func newGormLogger(logger log.Logger) *gormLogger {
	return &gormLogger{
		log: logger.Named(DefaultGormLoggerName),
	}
}

// mapGormLogLevelToZap converts the GORM log level value to the Zap log level.
// GORM logger.Info means DEBUG for Zap logger, Warn means INFO, etc.
func mapGormLogLevelToZap(logLevel logger.LogLevel) string {
	switch logLevel {
	case logger.Info:
		return "DEBUG"

	case logger.Warn:
		return "INFO"

	case logger.Error:
		return "WARN"

	case logger.Silent:
		return "ERROR"
	}

	return "INFO"
}

func (gl *gormLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	level := mapGormLogLevelToZap(logLevel)

	l, err := log.NewLogger(log.ConfigLogger{
		Level: level,
	})
	if err != nil {
		gl.log.Error(context.Background(), "Can't create a new GORM logger",
			log.Err(err))
		return gl
	}

	return newGormLogger(l)
}

func (gl *gormLogger) Info(ctx context.Context, format string, args ...interface{}) {
	gl.log.Debugf(ctx, format, args...)
}

func (gl *gormLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	gl.log.Warnf(ctx, format, args...)
}

func (gl *gormLogger) Error(ctx context.Context, format string, args ...interface{}) {
	gl.log.Errorf(ctx, format, args...)
}

func (gl *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), _ error) {
	sql, rows := fc()
	var rowsStr string
	if rows == -1 {
		rowsStr = "-"
	} else {
		rowsStr = strconv.Itoa(int(rows))
	}

	fileLine := utils.FileWithLineNum()
	elapsed := time.Since(begin)
	duration := float64(elapsed.Nanoseconds()) / 1e6

	gl.log.Debugf(ctx, "[GORM] %s [%.3fms] [rows:%s] %s", fileLine, duration, rowsStr, sql)
}
