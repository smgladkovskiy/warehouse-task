package multisql

import (
	"database/sql"
	"math/rand"
)

// Manager is a multiple databases manager.
type Manager struct {
	db        *sql.DB
	resolvers Resolvers
}

// Resolvers is a map of resolvers.
type Resolvers map[string]Resolver

// NewManager creates a new multiple databases manager.
func NewManager(db *sql.DB, resolvers Resolvers) *Manager {
	return &Manager{db: db, resolvers: resolvers}
}

// Default returns default database.
func (db *Manager) Default() *sql.DB {
	return db.db
}

// Get resolves database by the key. Returns default database if no resolver
// found, or resolver returns nil.
func (db *Manager) Get(key string) *sql.DB {
	resolver, ok := db.resolvers[key]
	if ok {
		if res := resolver(); res != nil {
			return res
		}
	}

	return db.Default()
}

// Resolver simply returns *sql.DB.
type Resolver func() *sql.DB

// NilResolver always returns nil.
func NilResolver() *sql.DB {
	return nil
}

// SingleResolver always returns db.
func SingleResolver(db *sql.DB) Resolver {
	return func() *sql.DB {
		return db
	}
}

// RandomResolver returns resolver that randomly selects *sql.DB from the
// sources. Returns NilResorver on empty sources. Returns SingleResolver if
// sources contains only one element.
func RandomResolver(sources []*sql.DB) Resolver {
	switch len(sources) {
	case 0:
		return NilResolver
	case 1:
		return SingleResolver(sources[0])
	default:
		return func() *sql.DB {
			return sources[rand.Intn(len(sources))] //nolint:gosec // doesn't require security
		}
	}
}
