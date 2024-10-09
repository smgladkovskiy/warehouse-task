package log

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
)

// A Field is a marshaling operation used to add a key-value pair to a logger's
// context. Most fields are lazily marshaled, so it's inexpensive to add fields
// to disabled debug-level log statements.
type Field zap.Field

var (
	_minTimeInt64 = time.Unix(0, math.MinInt64)
	_maxTimeInt64 = time.Unix(0, math.MaxInt64)
)

// nilField returns a field which will marshal explicitly as nil. See motivation
// in https://github.com/uber-go/zap/issues/753 . If we ever make breaking
// changes and add zapcore.NilType and zapcore.ObjectEncoder.AddNil, the
// implementation here should be changed to reflect that.
func nilField(key string) Field { return Reflect(key, nil) }

// Skip constructs a no-op field, which is often useful when handling invalid
// inputs in other Field constructors.
func Skip() Field {
	return Field{Type: zapcore.SkipType}
}

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) Field {
	return Field{Key: key, Type: zapcore.BinaryType, Interface: val}
}

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) Field {
	var ival int64
	if val {
		ival = 1
	}
	return Field{Key: key, Type: zapcore.BoolType, Integer: ival}
}

// Float64 constructs a field that carries a float64. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float64(key string, val float64) Field {
	return Field{Key: key, Type: zapcore.Float64Type, Integer: int64(math.Float64bits(val))}
}

// Int constructs a field with the given key and value.
func Int(key string, val int) Field {
	return Int64(key, int64(val))
}

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) Field {
	return Field{Key: key, Type: zapcore.Int64Type, Integer: val}
}

// Uint64 constructs a field with the given key and value.
func Uint64(key string, val uint64) Field {
	return Field{Key: key, Type: zapcore.Uint64Type, Integer: int64(val)}
}

// String constructs a field with the given key and value.
func String(key string, val string) Field {
	return Field{Key: key, Type: zapcore.StringType, String: val}
}

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) Field {
	return Field{Key: key, Type: zapcore.ReflectType, Interface: val}
}

// Object is a variant of Reflect where key is fixed to val's type name (without pointer prefix if any).
func Object(val interface{}) Field {
	typeName := "<nil>"
	if valType := reflect.TypeOf(val); valType != nil {
		typeName = valType.String()
	}
	return Reflect(strings.TrimLeft(typeName, "*"), val)
}

// group constructs a field with the given key and groupMarshaler. It
// provides a flexible, but still type-safe and efficient, way to add map- or
// struct-like user-defined types to the logging context. The struct's
// MarshalLogObject method is called lazily.
func group(key string, val groupMarshaler) Field {
	return Field{Key: key, Type: zapcore.ObjectMarshalerType, Interface: val}
}

// Namespace creates a named, isolated scope within the logger's context. All
// subsequent fields will be added to the new namespace.
//
// This helps prevent key collisions when injecting loggers into sub-components
// or third-party libraries.
func Namespace(key string) Field {
	return Field{Key: key, Type: zapcore.NamespaceType}
}

// Stringer constructs a field with the given key and the output of the value's
// String method. The Stringer's String method is called lazily.
func Stringer(key string, val fmt.Stringer) Field {
	return Field{Key: key, Type: zapcore.StringerType, Interface: val}
}

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) Field {
	if val.Before(_minTimeInt64) || val.After(_maxTimeInt64) {
		return Field{Key: key, Type: zapcore.TimeFullType, Interface: val}
	}
	return Field{Key: key, Type: zapcore.TimeType, Integer: val.UnixNano(), Interface: val.Location()}
}

// Timep constructs a field that carries a *time.Time. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Timep(key string, val *time.Time) Field {
	if val == nil {
		return nilField(key)
	}
	return Time(key, *val)
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: zapcore.DurationType, Integer: int64(val)}
}

// Proto constructs a field with the given key and value.
// Use method to log proto messages.
func Proto(key string, val proto.Message) Field {
	return Reflect(key, val)
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value interface{}) Field {
	switch val := value.(type) {
	case bool:
		return Bool(key, val)
	case float64:
		return Float64(key, val)
	case int:
		return Int(key, val)
	case string:
		return String(key, val)
	case []byte:
		return Binary(key, val)
	case time.Time:
		return Time(key, val)
	case *time.Time:
		return Timep(key, val)
	case time.Duration:
		return Duration(key, val)
	case proto.Message:
		return Proto(key, val)
	case fmt.Stringer:
		return Stringer(key, val)
	case error:
		return NamedErr(key, val)
	default:
		return Reflect(key, val)
	}
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) Field {
	return Field{Key: key, Type: zapcore.ByteStringType, Interface: val}
}

// NamedErr constructs a field that carries error value with given key.
func NamedErr(key string, val error) Field {
	fieldType := zapcore.ErrorType
	if val == nil {
		// Calling logging function on 'NamedErr' where 'val' is nil results in panic:
		//
		//   interface conversion: interface is nil, not error
		//
		// This is because our field's type is 'ErrorType'. Zap logger attempts to extract
		// error value from 'Interface' field of type 'interface{}', but that field is nil.
		//
		// Switching to 'ReflectType' makes Zap interpret nil values without errors.
		fieldType = zapcore.ReflectType
	}
	return Field{Key: key, Type: fieldType, Interface: val}
}

// Err is a variant of NamedErr where key value is fixed to "error".
func Err(val error) Field {
	return NamedErr("error", val)
}

// Array constructs a field with the given key and ArrayMarshaler. It provides
// a flexible, but still type-safe and efficient, way to add array-like types
// to the logging context. The struct's MarshalLogArray method is called lazily.
func Array(key string, val zapcore.ArrayMarshaler) Field {
	return Field{Key: key, Type: zapcore.ArrayMarshalerType, Interface: val}
}

// DurationMilli constructs a field with the given key and value. The encoder
// controls how the millisecond duration is serialized.
func DurationMilli(key string, val time.Duration) Field {
	return Field{Key: key + "_ms", Type: zapcore.Int64Type, Integer: val.Milliseconds()}
}

// DurationMicro constructs a field with the given key and value. The encoder
// controls how the microsecond duration is serialized.
func DurationMicro(key string, val time.Duration) Field {
	return Field{Key: key + "_us", Type: zapcore.Int64Type, Integer: val.Microseconds()}
}
