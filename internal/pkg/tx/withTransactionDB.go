package tx

import (
	"context"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/avito-tech/go-transaction-manager/trm/settings"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/db"
	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
)

//var _ base.Databaser = (*WithTransactionDB)(nil)

// WithTransactionDB добавляет поддержку транзакционности с управлением за рамками репозитория.
// На борту уже имеет base.WithDB, поэтому использовать его смысла нет.
type WithTransactionDB struct {
	db.WithDB
	trx *trmgorm.CtxGetter
}

// SetTransactionDB sets transaction object.
func (t *WithTransactionDB) SetTransactionDB(db *db.Instance, trx *trmgorm.CtxGetter) {
	if db == nil {
		panic("database instance is nil")
	}

	if trx == nil {
		panic("transaction CtxGetter is nil")
	}

	t.SetDB(db)

	t.trx = trx
}

func (t *WithTransactionDB) WriteDBTrx(ctx context.Context) *gorm.DB {
	return t.trx.TrOrDB(ctx, settings.DefaultCtxKey, t.WriteDB())
}

func (t *WithTransactionDB) GetQueryDB(ctx context.Context, qos queryOptions.QueryOptionable) (db *gorm.DB) {
	defer func() {
		if r := recover(); r != nil {
			db = t.AsyncDB().WithContext(ctx)
		}
	}()

	if qos == nil || qos == (*queryOptions.BasicQueryOptions)(nil) {
		return t.AsyncDB().WithContext(ctx)
	}

	if forUpdate := qos.IsForUpdate(); forUpdate {
		return t.WriteDBTrx(ctx).Clauses(clause.Locking{Strength: "UPDATE"})
	} else {
		return t.AsyncDB().WithContext(ctx)
	}
}
