package addproducttoorder

import (
	"context"
	"errors"
	"fmt"

	baseUUID "github.com/google/uuid"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/checker"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/log"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/tx"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	upsertOrder "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order/upsert"
	upsertOrderProduct "github.com/smgladkovskiy/warehouse-task/internal/service/commands/order_product/upsert"
	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	getOrderByID "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_order"
	getStocks "github.com/smgladkovskiy/warehouse-task/internal/service/queries/order/get_stocks"
	getProduct "github.com/smgladkovskiy/warehouse-task/internal/service/queries/product/get_product"
	usecase "github.com/smgladkovskiy/warehouse-task/internal/service/usecases"
)

type UseCase struct {
	uuid.WithUUIDGenerator
	now.WithNowGenerator
	checker.WithCheck
	tx.WithTransactionManager
	log.WithLogger

	// Query handlers
	getOrderQuery   *getOrderByID.QueryHandler
	getProductQuery *getProduct.QueryHandler
	getStocksQuery  *getStocks.QueryHandler

	// Command handlers
	upsertOrderCmd        *upsertOrder.CommandHandler
	upsertOrderProductCmd *upsertOrderProduct.CommandHandler
}

func NewUseCase(cfgs ...usecase.Configuration[*UseCase]) (*UseCase, error) {
	uc := &UseCase{}

	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		if cfg == nil {
			return nil, checker.ErrInitError
		}

		err := cfg(uc)
		if err != nil {
			return nil, err
		}
	}

	if err := uc.Check(*uc); err != nil {
		return nil, err
	}

	return uc, nil
}

func (uc *UseCase) Run(ctx context.Context, req Requestable) error {
	l := uc.Logger().With(
		log.String("orderUUID", req.GetOrderID().String()),
		log.String("userUUID", req.GetUserID().String()),
		log.String("productUUID", req.GetProductID().String()),
		log.Uint64("quantity", req.GetQuantity()),
	)

	l.Debug(ctx, "START usecase")

	if err := uc.TransactionDo(ctx, uc.transaction(l, req)); err != nil {
		l.Error(ctx, "STOP usecase! transaction error", log.Err(err))

		return fmt.Errorf("[addProductToOrder - uc.TransactionDo error]: %w", err)
	}

	l.Debug(ctx, "END usecase")

	return nil
}

func (uc *UseCase) transaction(l log.Logger, req Requestable) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// 1. Получаем заказ по ID (если есть ID и запись в БД), либо создаём новый
		order, err := uc.getOrder(ctx, req)
		if err != nil {
			return fmt.Errorf("[addProductToOrder - uc.getOrder error]: %w", err)
		}

		l = l.With(log.String("orderID", order.ID.String()))

		// 2. Получаем товар по id
		productQuery, err := getProduct.NewQueryByID(req.GetProductID())
		if err != nil {
			return fmt.Errorf("[addProductToOrder - getProduct.NewQueryByProductIDUnsafe error]: %w", err)
		}

		product, err := uc.getProductQuery.Handle(ctx, *productQuery)
		if err != nil {
			return fmt.Errorf("[addProductToOrder - uc.getProductQuery.Handle error]: %w", err)
		}

		// 3. Получаем количество товара на складе
		productStocks, err := uc.getStocksQuery.Handle(ctx, getStocks.NewQueryByProductIDUnsafe(product.ID))
		if err != nil {
			return fmt.Errorf("[addProductToOrder - uc.getStocksQuery.Handle error]: %w", err)
		}

		l = l.With(log.Uint64("productAvailableQuantity", productStocks.GetAvailableQuantity().Uint64()))

		// 4. Изменяем количество товара в заказе с проверкой на доступность указанного количества товара на складе
		if err = order.ChangeOrderProducts(productStocks, *product, req.GetQuantity()); err != nil {
			return fmt.Errorf("[addProductToOrder - order.ChangeProductAmount error]: %w", err)
		}

		// 5. Сохраняем заказ
		if err = uc.upsertOrderCmd.Handle(ctx, upsertOrder.NewCommandUnsafe(order)); err != nil {
			return fmt.Errorf("[addProductToOrder - uc.upsertOrderCmd.Run error]: %w", err)
		}

		// 6. Сохраняем товар в заказе
		orderProduct := order.GetOrderProductByProductIDUnsafe(product.ID)
		if err = uc.upsertOrderProductCmd.Handle(ctx, upsertOrderProduct.NewCommandUnsafe(orderProduct)); err != nil {
			return fmt.Errorf("[addProductToOrder - uc.upsertOrderProductCmd.Run error]: %w", err)
		}

		l = l.With(log.Uint64("orderProductQuantity", orderProduct.Quantity.Uint64()))

		return nil
	}
}

func (uc *UseCase) getOrder(ctx context.Context, req Requestable) (*entities.Order, error) {
	if req.GetOrderID() != baseUUID.Nil {
		// подавляем ошибку, так как orderID может быть пустым
		query, _ := getOrderByID.NewQueryForUpdate(req.GetOrderID())
		order, err := uc.getOrderQuery.Handle(ctx, *query)

		if errors.Is(err, entities.ErrOrderRecNotFound) {
			return uc.createOrder(ctx, req)
		}

		if err != nil {
			return nil, err
		}

		return order, nil
	}

	return uc.createOrder(ctx, req)
}

func (uc *UseCase) createOrder(ctx context.Context, req Requestable) (*entities.Order, error) {
	cmd, err := upsertOrder.NewCommand(
		req.GetUserID(),
		entities.WithUUIDFunc[*entities.Order](uc.GetUUIDGen()),
		entities.WithNowFunc[*entities.Order](uc.GetNowGen()),
	)
	if err != nil {
		return nil, fmt.Errorf("[addProductToOrder - upsertOrder.NewCommand error]: %w", err)
	}

	if err = uc.upsertOrderCmd.Handle(ctx, *cmd); err != nil {
		return nil, fmt.Errorf("[addProductToOrder - uc.createOrderCmd.Run error]: %w", err)
	}

	return cmd.GetOrder(), nil
}
