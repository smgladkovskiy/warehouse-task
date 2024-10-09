package db

//go:generate go run go.uber.org/mock/mockgen@v0.4.0 -destination internal/mocks/driver.go -package mocks database/sql/driver ExecerContext,QueryerContext,StmtExecContext,StmtQueryContext
