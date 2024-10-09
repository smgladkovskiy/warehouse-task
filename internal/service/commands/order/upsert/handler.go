package upsertorder

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
)

//go:generate mockgen -source=handler.go -destination=upsert_order_mock.go -package=upsertorder -mock_names OrderUpserter=UpsertOrderMock
type OrderUpserter interface {
	UpsertOrder(ctx context.Context, order *entities.Order) error
}

type CommandHandler struct {
	repo OrderUpserter
}

func NewCommandHandler(repo OrderUpserter) *CommandHandler {
	if repo == nil {
		panic("OrderUpserter repo is nil")
	}

	return &CommandHandler{repo: repo}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd Command) error {
	return h.repo.UpsertOrder(ctx, cmd.order)
}
