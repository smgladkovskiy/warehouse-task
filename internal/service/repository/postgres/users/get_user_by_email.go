package users

import (
	"context"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

func (r *Repository) GetByEmail(ctx context.Context, email vObject.Email) (*entities.User, error) {
	//TODO implement me
	panic("implement me")
}
