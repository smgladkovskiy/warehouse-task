// Package db provides access to databases. It uses the gorm. The package provides
// the ability to work with sync/async replicas using the gorm plugin.
package db

import (
	"database/sql"

	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/db/multisql"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
)

// DefaultLoggerName is a default logger name.
const DefaultLoggerName = "db"

// Instance contains objects for interacting with databases.
type Instance struct {
	Gorm   *gorm.DB
	SQL    *sql.DB
	Config Config

	manager    *multisql.Manager
	sourcesSQL []*sql.DB
	syncsSQL   []*sql.DB
	asyncsSQL  []*sql.DB
	gormSQL    []*sql.DB

	tracerProvider       trace.TracerProvider
	maxTracingQuerySize  int
	log                  log.Logger
	gormConfig           *gorm.Config
	applicationName      string
	preferSimpleProtocol bool
}

// Config contains information about the connection.
type Config struct {
	Adapter  string
	Addr     string
	Database string
	Multi    ConfigMulti
}

// ConfigMulti contains information about sources for multi connections.
type ConfigMulti struct {
	Sources       []string
	ReplicasSync  []string
	ReplicasAsync []string
}

// NewInstance initializes db connection.
func NewInstance(dsn string, opts ...Option) (*Instance, error) {
	inst := &Instance{
		log:        log.Named(DefaultLoggerName),
		gormConfig: &gorm.Config{},
	}

	for _, opt := range opts {
		opt(inst)
	}

	// TODO database initialisation here
	// ...

	return inst, nil
}
