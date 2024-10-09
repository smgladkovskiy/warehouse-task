package userregistration

import "time"

type Requestable interface {
	GetEmail() string
	GetFirstName() string
	GetLastName() string
	GetBirthDate() time.Time
	GetMaritalStatus() string
	GetPassword() string
}
