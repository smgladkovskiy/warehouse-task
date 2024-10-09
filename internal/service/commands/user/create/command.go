package createuser

import (
	"time"

	"github.com/smgladkovskiy/warehouse-task/internal/service/entities"
)

type Command struct {
	user *entities.User
}

func NewCommand(email, firstName, lastName, maritalStatus string, birthdate time.Time, opts ...entities.Option[*entities.User]) (*Command, error) {
	user, err := entities.NewUser(email, firstName, lastName, maritalStatus, birthdate, opts...)
	if err != nil {
		return nil, err
	}

	return &Command{user: user}, nil
}

func (c Command) GetUser() *entities.User {
	return c.user
}
