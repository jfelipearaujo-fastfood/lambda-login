// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	entities "github.com/jsfelipearaujo/lambda-login/src/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockToken is an autogenerated mock type for the Token type
type MockToken struct {
	mock.Mock
}

// CreateJwtToken provides a mock function with given fields: user
func (_m *MockToken) CreateJwtToken(user entities.User) (string, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for CreateJwtToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(entities.User) (string, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(entities.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(entities.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockToken creates a new instance of MockToken. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockToken(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockToken {
	mock := &MockToken{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
