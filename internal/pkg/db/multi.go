package db

import (
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

const (
	maxCountDSN            = 9
	sourcesDSNPrefix       = "pg_dsn_rw"
	replicasSyncDSNPrefix  = "pg_dsn_ro_sync"
	replicasAsyncDSNPrefix = "pg_dsn_ro_async"

	sourcesResolverName = "sources"
	syncResolverName    = "sync"
	asyncResolverName   = "async"
)

// Variables describe the operation to be performed.
var (
	useWriteResolver = dbresolver.Write
	useSyncResolver  = dbresolver.Read
	useAsyncResolver = dbresolver.Use(asyncResolverName)
)

// DB returns DB object.
func (i *Instance) DB() *gorm.DB {
	return i.Gorm
}

// WriteDB returns DB object for write access.
func (i *Instance) WriteDB() *gorm.DB {
	return i.DB().Clauses(useWriteResolver)
}

// SyncDB returns DB object for sync replica access.
func (i *Instance) SyncDB() *gorm.DB {
	return i.DB().Clauses(useSyncResolver)
}

// AsyncDB returns DB object for async replica access.
func (i *Instance) AsyncDB() *gorm.DB {
	return i.DB().Clauses(useAsyncResolver)
}

// WriteSQL returns *sql.DB for write access.
func (i *Instance) WriteSQL() *sql.DB {
	return i.manager.Default()
}

// SyncSQL returns *sql.DB for sync replica access.
func (i *Instance) SyncSQL() *sql.DB {
	return i.manager.Get(syncResolverName)
}

// AsyncSQL returns *sql.DB for async replica access.
func (i *Instance) AsyncSQL() *sql.DB {
	return i.manager.Get(asyncResolverName)
}
