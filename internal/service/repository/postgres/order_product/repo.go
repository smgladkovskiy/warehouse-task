package orderproduct

import (
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/db"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	trx "github.com/smgladkovskiy/warehouse-task/internal/pkg/tx"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	upsertorderproduct "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order_product/upsert"
)

type Repository struct {
	now.WithNowGenerator
	uuid.WithUUIDGenerator
	trx.WithTransactionDB
}

var (
	_ upsertorderproduct.OrderProductUpserter = (*Repository)(nil)
)

func NewRepository(db *db.Instance, trx *trmgorm.CtxGetter) *Repository {
	if db == nil {
		panic("database instance is nil")
	}

	if trx == nil {
		panic("transaction CtxGetter is nil")
	}

	r := Repository{}

	r.SetTransactionDB(db, trx)

	return &r
}
