package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLevelRegistry_GetOrSet(t *testing.T) {
	t.Run("stored", func(t *testing.T) {
		registry := NewLevelRegistry()
		lvl := zap.NewAtomicLevel()
		got, loaded := registry.GetOrSet("foo", lvl)
		assert.False(t, loaded)
		assert.Equal(t, lvl, got)
	})
	t.Run("loaded", func(t *testing.T) {
		lvl := zap.NewAtomicLevel()
		registry := provideLevelRegistry(t, "foo", lvl)
		got, loaded := registry.GetOrSet("foo", zap.NewAtomicLevelAt(zap.ErrorLevel))
		assert.True(t, loaded)
		assert.Equal(t, lvl, got)
	})
}

func provideLevelRegistry(t *testing.T, name string, lvl AtomicLevel) LevelRegistry {
	t.Helper()
	r := NewLevelRegistry()
	r.Set(name, lvl)
	return r
}
