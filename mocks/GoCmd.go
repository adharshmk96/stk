// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// GoCmd is an autogenerated mock type for the GoCmd type
type GoCmd struct {
	mock.Mock
}

// IsMod provides a mock function with given fields:
func (_m *GoCmd) IsMod() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ModInit provides a mock function with given fields: _a0
func (_m *GoCmd) ModInit(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ModPackageName provides a mock function with given fields:
func (_m *GoCmd) ModPackageName() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModTidy provides a mock function with given fields:
func (_m *GoCmd) ModTidy() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunCmd provides a mock function with given fields: args
func (_m *GoCmd) RunCmd(args ...string) (string, error) {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(...string) (string, error)); ok {
		return rf(args...)
	}
	if rf, ok := ret.Get(0).(func(...string) string); ok {
		r0 = rf(args...)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(...string) error); ok {
		r1 = rf(args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewGoCmd creates a new instance of GoCmd. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGoCmd(t interface {
	mock.TestingT
	Cleanup(func())
}) *GoCmd {
	mock := &GoCmd{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
