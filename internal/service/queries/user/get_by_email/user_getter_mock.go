// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go
//
// Generated by this command:
//
//	mockgen -source=handler.go -destination=user_getter_mock.go -package=getuserbyemail -mock_names UserGetter=GetUserMock
//

// Package getuserbyemail is a generated GoMock package.
package getuserbyemail

import (
	context "context"
	reflect "reflect"

	entities "github.com/smgladkovskiy/warehouse-task/internal/service/entities"
	valueobjects "github.com/smgladkovskiy/warehouse-task/internal/service/entities/value_objects"
	gomock "go.uber.org/mock/gomock"
)

// GetUserMock is a mock of UserGetter interface.
type GetUserMock struct {
	ctrl     *gomock.Controller
	recorder *GetUserMockMockRecorder
}

// GetUserMockMockRecorder is the mock recorder for GetUserMock.
type GetUserMockMockRecorder struct {
	mock *GetUserMock
}

// NewGetUserMock creates a new mock instance.
func NewGetUserMock(ctrl *gomock.Controller) *GetUserMock {
	mock := &GetUserMock{ctrl: ctrl}
	mock.recorder = &GetUserMockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *GetUserMock) EXPECT() *GetUserMockMockRecorder {
	return m.recorder
}

// GetByEmail mocks base method.
func (m *GetUserMock) GetByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", ctx, email)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *GetUserMockMockRecorder) GetByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*GetUserMock)(nil).GetByEmail), ctx, email)
}
