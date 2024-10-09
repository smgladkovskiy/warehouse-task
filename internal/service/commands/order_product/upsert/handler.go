package upsertorderproduct

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
)

//go:generate mockgen -source=handler.go -destination=upsert_order_product_mock.go -package=upsertorderproduct -mock_names OrderProductUpserter=UpsertOrderProductMock
type OrderProductUpserter interface {
	UpsertOrderProduct(ctx context.Context, orderProduct *entities.OrderProduct) error
}

type CommandHandler struct {
	repo OrderProductUpserter
}

func NewCommandHandler(repo OrderProductUpserter) *CommandHandler {
	if repo == nil {
		panic("OrderProductUpserter repo is nil")
	}

	return &CommandHandler{repo: repo}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd Command) error {
	return h.repo.UpsertOrderProduct(ctx, cmd.orderProduct)
}
