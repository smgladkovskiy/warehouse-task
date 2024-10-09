package log

import "go.uber.org/zap/zapcore"

// CheckWriteAction indicates what action to take after a log entry is processed.
type CheckWriteAction zapcore.CheckWriteAction
