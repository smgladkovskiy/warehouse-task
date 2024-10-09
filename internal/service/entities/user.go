package entities

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/smgladkovskiy/warehouse-task/internal/pkg/now"
	passCrypto "github.com/smgladkovskiy/warehouse-task/internal/pkg/pass_crypto"
	"github.com/smgladkovskiy/warehouse-task/internal/pkg/uuid"
	vObject "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
)

type User struct {
	now.WithNowGenerator
	uuid.WithUUIDGenerator
	passCrypto.WithPasswordHasher

	ID            vObject.UserID
	Email         vObject.Email
	FirstName     vObject.FirstName
	LastName      vObject.LastName
	BirthDate     vObject.Birthdate
	MaritalStatus vObject.MaritalStatus
	PasswordHash  vObject.PasswordHash
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time

	Orders []Order
}

var (
	ErrUserRecNotFound   = errors.New("user record not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

func (u User) FullName() string {
	fullName := strings.Builder{}
	if u.FirstName != "" {
		fullName.WriteString(string(u.FirstName))
	}

	if u.LastName != "" {
		fullName.WriteString(" ")
		fullName.WriteString(string(u.LastName))
	}

	return fullName.String()
}

func NewUser(email, firstName, lastName, maritalStatus string, birthdate time.Time, opts ...Option[*User]) (*User, error) {
	var (
		u   User
		err error
	)

	u.Email, err = vObject.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("[NewUser - NewEmail] %w", err)
	}

	u.FirstName, err = vObject.NewFirstName(firstName)
	if err != nil {
		return nil, fmt.Errorf("[NewUser - NewFirstName] %w", err)
	}

	u.LastName, err = vObject.NewLastName(lastName)
	if err != nil {
		return nil, fmt.Errorf("[NewUser - NewLastName] %w", err)
	}

	u.MaritalStatus, err = vObject.NewMaritalStatus(maritalStatus)
	if err != nil {
		return nil, fmt.Errorf("[NewUser - NewMaritalStatus] %w", err)
	}

	u.BirthDate, err = vObject.NewBirthDate(birthdate)
	if err != nil {
		return nil, fmt.Errorf("[NewUser - NewAge] %w", err)
	}

	for _, opt := range opts {
		if err = opt(&u); err != nil {
			return nil, err
		}
	}

	tn := u.Now()
	u.ID = vObject.NewUserIDFromUUIDUnsafe(u.UUID())
	u.CreatedAt = tn
	u.UpdatedAt = tn

	return &u, nil
}
