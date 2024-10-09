package tx

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
)

type trxManagerLog struct {
	log.Logger
}

func NewTrxLogger(l log.Logger) *trxManagerLog {
	return &trxManagerLog{Logger: l.Named("TransactionManager")}
}

func (tl *trxManagerLog) Warning(ctx context.Context, msg string) {
	tl.Warn(ctx, msg, log.Any("ctx", ctx))
}
