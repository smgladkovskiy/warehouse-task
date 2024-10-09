package getuserbyemail

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	vObjects "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

//go:generate mockgen -source=handler.go -destination=user_getter_mock.go -package=getuserbyemail -mock_names UserGetter=GetUserMock
type UserGetter interface {
	GetByEmail(ctx context.Context, email vObjects.Email) (*entities.User, error)
}

type QueryHandler struct {
	repo UserGetter
}

func NewQueryHandler(repo UserGetter) *QueryHandler {
	if repo == nil {
		panic("UserGetter repo is nil")
	}

	return &QueryHandler{repo: repo}
}

func (h *QueryHandler) Handle(ctx context.Context, q Query) (*entities.User, error) {
	return h.repo.GetByEmail(ctx, q.email)
}
