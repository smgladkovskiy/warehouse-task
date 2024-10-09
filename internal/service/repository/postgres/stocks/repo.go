package stocks

import (
	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/db"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	trx "github.com/smgladkovskiy/warehouse-task/internal/pkg/tx"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	getstocks "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_stocks"
)

type Repository struct {
	now.WithNowGenerator
	uuid.WithUUIDGenerator
	trx.WithTransactionDB
}

var _ getstocks.StocksGetter = (*Repository)(nil)

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
