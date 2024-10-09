package log

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log/internal/testdata"
)

func TestFieldConstructors(t *testing.T) {
	ints := []int{5, 6}
	intsP := &ints

	sampleError := errors.New("sample error")
	protoMessage := &testdata.Request{Name: "Alice", Age: 10}

	tests := []struct {
		name   string
		field  Field
		expect Field
	}{
		{"Binary", Field{Key: "k", Type: zapcore.BinaryType, Interface: []byte("ab12")}, Binary("k", []byte("ab12"))},
		{"Bool", Field{Key: "k", Type: zapcore.BoolType, Integer: 1}, Bool("k", true)},
		{"Bool", Field{Key: "k", Type: zapcore.BoolType, Integer: 1}, Bool("k", true)},
		{"Duration", Field{Key: "k", Type: zapcore.DurationType, Integer: 1}, Duration("k", 1)},
		{"Int", Field{Key: "k", Type: zapcore.Int64Type, Integer: 1}, Int("k", 1)},
		{"Int64", Field{Key: "k", Type: zapcore.Int64Type, Integer: 1}, Int64("k", 1)},
		{"String", Field{Key: "k", Type: zapcore.StringType, String: "foo"}, String("k", "foo")},
		{"Time", Field{Key: "k", Type: zapcore.TimeType, Integer: 0, Interface: time.UTC}, Time("k", time.Unix(0, 0).In(time.UTC))},
		{"Time", Field{Key: "k", Type: zapcore.TimeType, Integer: 1000, Interface: time.UTC}, Time("k", time.Unix(0, 1000).In(time.UTC))},
		{"Time", Field{Key: "k", Type: zapcore.TimeType, Integer: math.MinInt64, Interface: time.UTC}, Time("k",
			time.Unix(0, math.MinInt64).In(time.UTC))},
		{"Time", Field{Key: "k", Type: zapcore.TimeType, Integer: math.MaxInt64, Interface: time.UTC}, Time("k",
			time.Unix(0, math.MaxInt64).In(time.UTC))},
		{"Time", Field{Key: "k", Type: zapcore.TimeFullType, Interface: time.Time{}}, Time("k", time.Time{})},
		{"Time", Field{Key: "k", Type: zapcore.TimeFullType, Interface: time.Unix(math.MaxInt64, 0)}, Time("k", time.Unix(math.MaxInt64, 0))},
		{"Reflect", Field{Key: "k", Type: zapcore.ReflectType, Interface: ints}, Reflect("k", ints)},
		{"Reflect", Field{Key: "k", Type: zapcore.ReflectType}, Reflect("k", nil)},
		{"Object", Field{Key: "[]int", Type: zapcore.ReflectType, Interface: ints}, Object(ints)},
		{"Object", Field{Key: "<nil>", Type: zapcore.ReflectType}, Object(nil)},
		{"Object", Field{Key: "[]int", Type: zapcore.ReflectType, Interface: &ints}, Object(&ints)},                                      // pointer case
		{"Object", Field{Key: "[]int", Type: zapcore.ReflectType, Interface: &intsP}, Object(&intsP)},                                    // pointer to pointer case
		{"Object", Field{Key: "struct {}", Type: zapcore.ReflectType, Interface: struct{}{}}, Object(struct{}{})},                        // anon struct
		{"Object", Field{Key: "struct {}", Type: zapcore.ReflectType, Interface: &struct{}{}}, Object(&struct{}{})},                      // pointer to anon struct
		{"Object", Field{Key: "struct { x int }", Type: zapcore.ReflectType, Interface: struct{ x int }{}}, Object(struct{ x int }{})},   // anon struct with field
		{"Object", Field{Key: "struct { x int }", Type: zapcore.ReflectType, Interface: &struct{ x int }{}}, Object(&struct{ x int }{})}, // pointer to anon struct with field
		{"Namespace", Namespace("k"), Field{Key: "k", Type: zapcore.NamespaceType}},
		{"ByteString", Field{Key: "k", Type: zapcore.ByteStringType, Interface: []byte("ab12")}, ByteString("k", []byte("ab12"))},
		{"Err", Field{Key: "error", Type: zapcore.ErrorType, Interface: sampleError}, Err(sampleError)},
		{"NamedErr", Field{Key: "k", Type: zapcore.ErrorType, Interface: sampleError}, NamedErr("k", sampleError)},
		{"ErrWithNil", Field{Key: "error", Type: zapcore.ReflectType, Interface: nil}, Err(nil)},
		{"NamedErrWithNil", Field{Key: "k", Type: zapcore.ReflectType, Interface: nil}, NamedErr("k", nil)},
		{"AnyWithError", Any("k", sampleError), NamedErr("k", sampleError)},
		{"Proto", Any("k", protoMessage), Proto("k", protoMessage)},
		{"DurationMilli", DurationMilli("k", time.Millisecond*15), Field{Key: "k_ms", Type: zapcore.Int64Type, Integer: (time.Millisecond * 15).Milliseconds()}},
		{"DurationMicro", DurationMicro("k", time.Microsecond*10), Field{Key: "k_us", Type: zapcore.Int64Type, Integer: (time.Microsecond * 10).Microseconds()}},
	}

	for _, tt := range tests {
		if !assert.Equal(t, tt.expect, tt.field, "Unexpected output from convenience field constructor %s.", tt.name) {
			t.Logf("type expected: %T\nGot: %T", tt.expect.Interface, tt.field.Interface)
		}
	}
}
