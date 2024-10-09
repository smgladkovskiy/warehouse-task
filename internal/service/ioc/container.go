package ioc

import (
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
	upsertOrder "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order/upsert"
	upsertOrderProduct "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order_product/upsert"
	createUser "github.com/smgladkovskiy/warehouse-task/internal/service/commands/user/create"
	getOrder "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_order"
	getStocks "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_stocks"
	getProduct "github.com/smgladkovskiy/warehouse-task/internal/service/queries/product/get_product"
	getUserByEmail "github.com/smgladkovskiy/warehouse-task/internal/service/queries/user/get_by_email"
	usecase "github.com/smgladkovskiy/warehouse-task/internal/service/usecases"
	addProductToOrder "github.com/smgladkovskiy/warehouse-task/internal/service/usecases/order/add_product_to_order"
	userRegistration "github.com/smgladkovskiy/warehouse-task/internal/service/usecases/user/registration"
)

type Container struct {
	Queries  Queries
	Commands Commands
	UseCases UseCases
}

type Queries struct {
	// order
	GetOrder  *getOrder.QueryHandler
	GetStocks *getStocks.QueryHandler

	// product
	GetProduct *getProduct.QueryHandler

	// user
	GetUserByEmail *getUserByEmail.QueryHandler
}

type Commands struct {
	// order
	UpsertOrder *upsertOrder.CommandHandler

	// order product
	UpsertOrderProduct *upsertOrderProduct.CommandHandler

	// user
	CreateUser *createUser.CommandHandler
}

type UseCases struct {
	// order
	AddProductToOrder *addProductToOrder.UseCase

	// user
	UserRegistration *userRegistration.UseCase
}

func NewContainer(realisations Implementationable) (*Container, error) {
	c := Container{
		Queries: Queries{
			GetOrder:       getOrder.NewQueryHandler(realisations.OrderGetter()),
			GetStocks:      getStocks.NewQueryHandler(realisations.StocksGetter()),
			GetProduct:     getProduct.NewQueryHandler(realisations.ProductGetter()),
			GetUserByEmail: getUserByEmail.NewQueryHandler(realisations.UserGetter()),
		},
		Commands: Commands{
			UpsertOrder:        upsertOrder.NewCommandHandler(realisations.OrderUpserter()),
			UpsertOrderProduct: upsertOrderProduct.NewCommandHandler(realisations.OrderProductUpserter()),
			CreateUser:         createUser.NewCommandHandler(realisations.UserCreator()),
		},
	}

	var err error

	c.UseCases.AddProductToOrder, err = addProductToOrder.NewUseCase(
		addProductToOrder.WithGetOrderQuery(c.Queries.GetOrder),
		addProductToOrder.WithGetProductQuery(c.Queries.GetProduct),
		addProductToOrder.WithGetStocksQuery(c.Queries.GetStocks),
		addProductToOrder.WithUpsertOrderCommand(c.Commands.UpsertOrder),
		addProductToOrder.WithUpsertOrderProductCommand(c.Commands.UpsertOrderProduct),
		usecase.WithTransactionManager[*addProductToOrder.UseCase](realisations.TransactionManager()),
		usecase.WithLogger[*addProductToOrder.UseCase](log.Named("usecase.addProductToOrder")),
	)

	c.UseCases.UserRegistration, err = userRegistration.NewUseCase(
		userRegistration.WithGetUserByEmailQuery(c.Queries.GetUserByEmail),
		userRegistration.WithCreateUserCommand(c.Commands.CreateUser),
		usecase.WithLogger[*userRegistration.UseCase](log.Named("usecase.userRegistration")),
	)

	return &c, err
}
