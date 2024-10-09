package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestFilterInternalKeys(t *testing.T) {
	tests := []struct {
		fields       []Field
		expected     []Field
		excludedKeys []string
	}{
		{
			fields:       []Field{{Key: "abcd"}, {Key: "abcd"}},
			expected:     []Field{{Key: "abcd"}, {Key: "abcd"}},
			excludedKeys: nil,
		},
		{
			fields:       []Field{{Key: messageKey}, {Key: "abcd"}},
			expected:     []Field{{Key: "abcd"}},
			excludedKeys: []string{messageKey},
		},
		{
			fields:       []Field{{Key: messageKey}, {Key: messageKey}},
			expected:     []Field{},
			excludedKeys: []string{messageKey, messageKey},
		},
		{
			fields:       []Field{{Key: messageKey}, {Key: "abcd"}, {Key: timestampKey}},
			expected:     []Field{{Key: "abcd"}},
			excludedKeys: []string{messageKey, timestampKey},
		},
		{
			fields:       []Field{{Key: "abcd"}, {Key: levelKey}, {Key: "bcde"}},
			expected:     []Field{{Key: "abcd"}, {Key: "bcde"}},
			excludedKeys: []string{levelKey},
		},
		{
			fields:       []Field{{Key: "abcd"}, {Key: levelKey}, {Key: "bcde"}, {Key: spanKey}},
			expected:     []Field{{Key: "abcd"}, {Key: "bcde"}},
			excludedKeys: []string{levelKey, spanKey},
		},
	}

	for _, tc := range tests {
		fields, keys := filterInternalKeys(tc.fields)
		assert.Equal(t, tc.expected, fields)
		assert.Equal(t, tc.excludedKeys, keys)
	}
}

func TestCtxToFields(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// nil
	assert.Empty(ctxToFields(nil)) //nolint:staticcheck // checking nil case.

	// empty
	assert.Empty(ctxToFields(ctx))

	// with trace
	exp := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exp))
	defer func() {
		_ = tp.Shutdown(ctx)
	}()

	ctx, span := tp.Tracer("").Start(ctx, "test span")
	span.End()

	fields := ctxToFields(ctx)

	assert.Len(fields, 2)
	assert.NotEmpty(fields[0])
	assert.NotEmpty(fields[1])
}

func BenchmarkCtxToFields(b *testing.B) {
	b.Run("nil", func(b *testing.B) {
		var ctx context.Context

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctxToFields(ctx)
		}
	})

	b.Run("empty", func(b *testing.B) {
		ctx := context.Background()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctxToFields(ctx)
		}
	})

	b.Run("with trace", func(b *testing.B) {
		ctx := context.Background()

		exp := tracetest.NewInMemoryExporter()
		tp := trace.NewTracerProvider(trace.WithSyncer(exp))
		defer func() {
			_ = tp.Shutdown(ctx)
		}()

		ctx, span := tp.Tracer("").Start(ctx, "test span")
		defer span.End()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctxToFields(ctx)
		}
	})
}
