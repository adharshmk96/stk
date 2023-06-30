// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	gsk "github.com/adharshmk96/stk/gsk"

	logrus "github.com/sirupsen/logrus"

	mock "github.com/stretchr/testify/mock"
)

// Context is an autogenerated mock type for the Context type
type Context struct {
	mock.Mock
}

// DecodeJSONBody provides a mock function with given fields: v
func (_m *Context) DecodeJSONBody(v interface{}) error {
	ret := _m.Called(v)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllowedOrigins provides a mock function with given fields:
func (_m *Context) GetAllowedOrigins() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetCookie provides a mock function with given fields: name
func (_m *Context) GetCookie(name string) (*http.Cookie, error) {
	ret := _m.Called(name)

	var r0 *http.Cookie
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*http.Cookie, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *http.Cookie); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Cookie)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetParam provides a mock function with given fields: key
func (_m *Context) GetParam(key string) string {
	ret := _m.Called(key)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetQueryParam provides a mock function with given fields: key
func (_m *Context) GetQueryParam(key string) string {
	ret := _m.Called(key)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetRequest provides a mock function with given fields:
func (_m *Context) GetRequest() *http.Request {
	ret := _m.Called()

	var r0 *http.Request
	if rf, ok := ret.Get(0).(func() *http.Request); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Request)
		}
	}

	return r0
}

// GetWriter provides a mock function with given fields:
func (_m *Context) GetWriter() http.ResponseWriter {
	ret := _m.Called()

	var r0 http.ResponseWriter
	if rf, ok := ret.Get(0).(func() http.ResponseWriter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.ResponseWriter)
		}
	}

	return r0
}

// JSONResponse provides a mock function with given fields: data
func (_m *Context) JSONResponse(data interface{}) {
	_m.Called(data)
}

// Logger provides a mock function with given fields:
func (_m *Context) Logger() *logrus.Logger {
	ret := _m.Called()

	var r0 *logrus.Logger
	if rf, ok := ret.Get(0).(func() *logrus.Logger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*logrus.Logger)
		}
	}

	return r0
}

// RawResponse provides a mock function with given fields: raw
func (_m *Context) RawResponse(raw []byte) {
	_m.Called(raw)
}

// SetCookie provides a mock function with given fields: cookie
func (_m *Context) SetCookie(cookie *http.Cookie) {
	_m.Called(cookie)
}

// SetHeader provides a mock function with given fields: _a0, _a1
func (_m *Context) SetHeader(_a0 string, _a1 string) {
	_m.Called(_a0, _a1)
}

// Status provides a mock function with given fields: status
func (_m *Context) Status(status int) gsk.Context {
	ret := _m.Called(status)

	var r0 gsk.Context
	if rf, ok := ret.Get(0).(func(int) gsk.Context); ok {
		r0 = rf(status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gsk.Context)
		}
	}

	return r0
}

type mockConstructorTestingTNewContext interface {
	mock.TestingT
	Cleanup(func())
}

// NewContext creates a new instance of Context. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewContext(t mockConstructorTestingTNewContext) *Context {
	mock := &Context{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
