package userregistration

import "time"

type testRequest struct {
	email         string
	firstName     string
	lastName      string
	birthdate     time.Time
	maritalStatus string
	password      string
}

var _ Requestable = (*testRequest)(nil)

func (t testRequest) GetEmail() string {
	return t.email
}

func (t testRequest) GetFirstName() string {
	return t.firstName
}

func (t testRequest) GetLastName() string {
	return t.lastName
}

func (t testRequest) GetBirthDate() time.Time {
	return t.birthdate
}

func (t testRequest) GetMaritalStatus() string {
	return t.maritalStatus
}

func (t testRequest) GetPassword() string {
	return t.password
}
