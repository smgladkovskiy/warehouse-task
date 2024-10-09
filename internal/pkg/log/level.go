package log

import (
	"sync"

	"go.uber.org/zap"
)

// LoggerWithLevelRegistry represents logger with a level registry.
type LoggerWithLevelRegistry interface {
	LevelRegistry() LevelRegistry
}

// An AtomicLevel is an atomically changeable, dynamic logging level. It lets
// you safely change the log level of a tree of loggers (the root logger and
// any children created by adding context) at runtime.
type AtomicLevel = zap.AtomicLevel

// LevelRegistry is a level registry for named loggers.
type LevelRegistry struct {
	m *sync.Map
}

// NewLevelRegistry creates a new level registry.
func NewLevelRegistry() LevelRegistry {
	return LevelRegistry{m: &sync.Map{}}
}

// Get gets a level by the logger name.
func (r LevelRegistry) Get(name string) (AtomicLevel, bool) {
	if v, ok := r.m.Load(name); ok {
		return v.(AtomicLevel), true //nolint:forcetypeassert // always AtomicLevel
	}
	return zap.AtomicLevel{}, false
}

// Set sets the named logger level.
func (r LevelRegistry) Set(name string, level AtomicLevel) {
	r.m.Store(name, level)
}

// GetOrSet returns the existing level for the logger name if present.
// Otherwise, it stores and returns the given level. The result flag is true if
// the logger was loaded, false if stored.
func (r LevelRegistry) GetOrSet(name string, level AtomicLevel) (AtomicLevel, bool) {
	v, loaded := r.m.LoadOrStore(name, level)
	return v.(AtomicLevel), loaded //nolint:forcetypeassert // always AtomicLevel
}
