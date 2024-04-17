// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	entities "github.com/jsfelipearaujo/lambda-login/src/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockDatabase is an autogenerated mock type for the Database type
type MockDatabase struct {
	mock.Mock
}

// GetUserByCPF provides a mock function with given fields: cpf
func (_m *MockDatabase) GetUserByCPF(cpf string) (entities.User, error) {
	ret := _m.Called(cpf)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByCPF")
	}

	var r0 entities.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (entities.User, error)); ok {
		return rf(cpf)
	}
	if rf, ok := ret.Get(0).(func(string) entities.User); ok {
		r0 = rf(cpf)
	} else {
		r0 = ret.Get(0).(entities.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(cpf)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockDatabase creates a new instance of MockDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabase {
	mock := &MockDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
