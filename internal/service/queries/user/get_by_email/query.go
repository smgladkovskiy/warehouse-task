package getuserbyemail

import (
	"fmt"

	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type Query struct {
	email vObject.Email
}

func NewQuery(email string) (*Query, error) {
	var (
		q   Query
		err error
	)

	q.email, err = vObject.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("[get_user_by_email.NewQuery] %w", err)
	}

	return &q, nil
}
