package createuser

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
)

//go:generate mockgen -source=handler.go -destination=user_creator_mock.go -package=createuser -mock_names UserCreator=CreateUserMock
type UserCreator interface {
	CreateUser(ctx context.Context, user *entities.User) error
}

type CommandHandler struct {
	repo UserCreator
}

func NewCommandHandler(repo UserCreator) *CommandHandler {
	if repo == nil {
		panic("UserCreator repo is nil")
	}

	return &CommandHandler{repo: repo}
}

func (h *CommandHandler) Handle(ctx context.Context, cmd Command) error {
	return h.repo.CreateUser(ctx, cmd.user)
}
