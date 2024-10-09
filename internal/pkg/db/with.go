package db

import (
	"database/sql"

	"gorm.io/gorm"
)

var _ Databaser = (*WithDB)(nil)

// Databaser is an interface for set and get *db.Instance.
// It is used by the application to check if a type embeds WithDB.
type Databaser interface {
	SetDB(*Instance)
	DB() *gorm.DB
	SQL() *sql.DB
	MultiDatabaser
	MultiSQLDatabaser
}

// MultiDatabaser is an interface to get write, sync, async DBs.
type MultiDatabaser interface {
	WriteDB() *gorm.DB
	SyncDB() *gorm.DB
	AsyncDB() *gorm.DB
}

// MultiSQLDatabaser is an interface to get write, sync, async *sql.DB.
type MultiSQLDatabaser interface {
	WriteSQL() *sql.DB
	SyncSQL() *sql.DB
	AsyncSQL() *sql.DB
}

// WithDB provides provides integration with DB.
type WithDB struct {
	db *Instance
}

// SetDB sets DB object.
func (c *WithDB) SetDB(db *Instance) {
	c.db = db
}

// DB returns DB object.
func (c WithDB) DB() *gorm.DB {
	return c.db.DB()
}

// WriteDB returns DB object for write access.
func (c WithDB) WriteDB() *gorm.DB {
	return c.db.WriteDB()
}

// SyncDB returns DB object for sync replica access.
func (c WithDB) SyncDB() *gorm.DB {
	return c.db.SyncDB()
}

// AsyncDB returns DB object for async replica access.
func (c WithDB) AsyncDB() *gorm.DB {
	return c.db.AsyncDB()
}

// SQL returns SQL object.
func (c WithDB) SQL() *sql.DB {
	return c.db.SQL
}

// WriteSQL returns *sql.DB for write access.
func (c WithDB) WriteSQL() *sql.DB {
	return c.db.WriteSQL()
}

// SyncSQL returns *sql.DB for sync replica access.
func (c WithDB) SyncSQL() *sql.DB {
	return c.db.SyncSQL()
}

// AsyncSQL returns *sql.DB for async replica access.
func (c WithDB) AsyncSQL() *sql.DB {
	return c.db.AsyncSQL()
}
