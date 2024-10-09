package log

import "go.uber.org/zap/zapcore"

// groupMarshaler allows user-defined types to efficiently add themselves to the
// logging context, and to selectively omit information which shouldn't be
// included in logs (e.g., passwords).
//
// Note: groupMarshaler is only used when zap.Object is used or when
// passed directly to zap.Any. It is not used when reflection-based
// encoding is used.
type groupMarshaler = zapcore.ObjectMarshaler
