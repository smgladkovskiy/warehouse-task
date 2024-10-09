package products

import (
	"context"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/db"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	trx "github.com/smgladkovskiy/warehouse-task/internal/pkg/tx"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	queryOptions "github.com/smgladkovskiy/warehouse-task/internal/service/entities/query_options"
	getProduct "github.com/smgladkovskiy/warehouse-task/internal/service/queries/product/get_product"
)

type Repository struct {
	now.WithNowGenerator
	uuid.WithUUIDGenerator
	trx.WithTransactionDB
}

var _ getProduct.ProductGetter = (*Repository)(nil)

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

func (r *Repository) GetProduct(ctx context.Context, qos queryOptions.ProductQueryOptionable) (*entities.Product, error) {
	//TODO implement me
	panic("implement me")
}
