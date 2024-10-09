package ioc

import (
	"github.com/avito-tech/go-transaction-manager/trm"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/application"
	upsertOrder "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order/upsert"
	upsertOrderProduct "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order_product/upsert"
	createUser "github.com/smgladkovskiy/warehouse-task/internal/service/commands/user/create"
	getOrderByID "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_order"
	getStocks "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_stocks"
	getProduct "github.com/smgladkovskiy/warehouse-task/internal/service/queries/product/get_product"
	getUserByEmail "github.com/smgladkovskiy/warehouse-task/internal/service/queries/user/get_by_email"
	orderProducts "github.com/smgladkovskiy/warehouse-task/internal/service/repository/postgres/order_product"
	"github.com/smgladkovskiy/warehouse-task/internal/service/repository/postgres/orders"
	"github.com/smgladkovskiy/warehouse-task/internal/service/repository/postgres/products"
	"github.com/smgladkovskiy/warehouse-task/internal/service/repository/postgres/stocks"
	"github.com/smgladkovskiy/warehouse-task/internal/service/repository/postgres/users"
)

type Implementationable interface {
	OrderGetter() getOrderByID.OrderGetter
	StocksGetter() getStocks.StocksGetter
	ProductGetter() getProduct.ProductGetter
	UserGetter() getUserByEmail.UserGetter

	OrderUpserter() upsertOrder.OrderUpserter
	OrderProductUpserter() upsertOrderProduct.OrderProductUpserter
	UserCreator() createUser.UserCreator
	TransactionManager() trm.Manager
}

type Implementations struct {
	txManager        trm.Manager
	orderRepo        *orders.Repository
	stockRepo        *stocks.Repository
	productRepo      *products.Repository
	userRepo         *users.Repository
	orderProductRepo *orderProducts.Repository
}

var _ Implementationable = (*Implementations)(nil)

func NewImplementations(app *application.App) *Implementations {
	return &Implementations{
		orderRepo:   orders.NewRepository(app.DB, app.TrxGetter),
		stockRepo:   stocks.NewRepository(app.DB, app.TrxGetter),
		productRepo: products.NewRepository(app.DB, app.TrxGetter),
		userRepo:    users.NewRepository(app.DB, app.TrxGetter),
		txManager:   app.TxManager,
	}
}

func (i *Implementations) OrderGetter() getOrderByID.OrderGetter {
	return i.orderRepo
}

func (i *Implementations) StocksGetter() getStocks.StocksGetter {
	return i.stockRepo
}

func (i *Implementations) ProductGetter() getProduct.ProductGetter {
	return i.productRepo
}

func (i *Implementations) UserGetter() getUserByEmail.UserGetter {
	return i.userRepo
}

func (i *Implementations) OrderUpserter() upsertOrder.OrderUpserter {
	return i.orderRepo
}

func (i *Implementations) OrderProductUpserter() upsertOrderProduct.OrderProductUpserter {
	return i.orderProductRepo
}

func (i *Implementations) UserCreator() createUser.UserCreator {
	return i.userRepo
}

func (i *Implementations) TransactionManager() trm.Manager {
	return i.txManager
}
