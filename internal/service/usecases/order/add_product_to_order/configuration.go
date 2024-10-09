package addproducttoorder

import (
	"fmt"

	upsertOrder "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order/upsert"
	upsertOrderProduct "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order_product/upsert"
	getOrderByID "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_order"
	getStocks "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_stocks"
	getProduct "github.com/smgladkovskiy/warehouse-task/internal/service/queries/product/get_product"
	usecase "github.com/smgladkovskiy/warehouse-task/internal/service/usecases"
)

func WithGetOrderQuery(handler *getOrderByID.QueryHandler) usecase.Configuration[*UseCase] {
	return func(uc *UseCase) error {
		if handler == nil {
			return fmt.Errorf("%w %s", usecase.ErrEmptyStructParam, "getOrderByID")
		}

		uc.getOrderQuery = handler

		return nil
	}
}

func WithGetProductQuery(handler *getProduct.QueryHandler) usecase.Configuration[*UseCase] {
	return func(uc *UseCase) error {
		if handler == nil {
			return fmt.Errorf("%w %s", usecase.ErrEmptyStructParam, "getProduct")
		}

		uc.getProductQuery = handler

		return nil
	}
}

func WithGetStocksQuery(handler *getStocks.QueryHandler) usecase.Configuration[*UseCase] {
	return func(uc *UseCase) error {
		if handler == nil {
			return fmt.Errorf("%w %s", usecase.ErrEmptyStructParam, "getStocks")
		}

		uc.getStocksQuery = handler

		return nil
	}
}

func WithUpsertOrderCommand(handler *upsertOrder.CommandHandler) usecase.Configuration[*UseCase] {
	return func(uc *UseCase) error {
		if handler == nil {
			return fmt.Errorf("%w %s", usecase.ErrEmptyStructParam, "upsertOrder")
		}

		uc.upsertOrderCmd = handler

		return nil
	}
}

func WithUpsertOrderProductCommand(handler *upsertOrderProduct.CommandHandler) usecase.Configuration[*UseCase] {
	return func(uc *UseCase) error {
		if handler == nil {
			return fmt.Errorf("%w %s", usecase.ErrEmptyStructParam, "upsertOrderProduct")
		}

		uc.upsertOrderProductCmd = handler

		return nil
	}
}
