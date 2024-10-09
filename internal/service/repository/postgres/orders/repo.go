package orders

import (
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/db"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	trx "github.com/smgladkovskiy/warehouse-task/internal/pkg/tx"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	upsertOrder "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order/upsert"
	getOrder "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_order"
)

type Repository struct {
	now.WithNowGenerator
	uuid.WithUUIDGenerator
	trx.WithTransactionDB
}

var (
	_ getOrder.OrderGetter      = (*Repository)(nil)
	_ upsertOrder.OrderUpserter = (*Repository)(nil)
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
