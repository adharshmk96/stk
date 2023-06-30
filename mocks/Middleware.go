// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	gsk "github.com/adharshmk96/stk/gsk"
	mock "github.com/stretchr/testify/mock"
)

// Middleware is an autogenerated mock type for the Middleware type
type Middleware struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *Middleware) Execute(_a0 gsk.HandlerFunc) gsk.HandlerFunc {
	ret := _m.Called(_a0)

	var r0 gsk.HandlerFunc
	if rf, ok := ret.Get(0).(func(gsk.HandlerFunc) gsk.HandlerFunc); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gsk.HandlerFunc)
		}
	}

	return r0
}

type mockConstructorTestingTNewMiddleware interface {
	mock.TestingT
	Cleanup(func())
}

// NewMiddleware creates a new instance of Middleware. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMiddleware(t mockConstructorTestingTNewMiddleware) *Middleware {
	mock := &Middleware{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}