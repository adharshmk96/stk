// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	sqlBuilder "github.com/adharshmk96/stk/pkg/sqlBuilder"
	mock "github.com/stretchr/testify/mock"
)

// SqlQuery is an autogenerated mock type for the SqlQuery type
type SqlQuery struct {
	mock.Mock
}

// Build provides a mock function with given fields:
func (_m *SqlQuery) Build() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// DeleteFrom provides a mock function with given fields: table
func (_m *SqlQuery) DeleteFrom(table string) sqlBuilder.SqlQuery {
	ret := _m.Called(table)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(string) sqlBuilder.SqlQuery); ok {
		r0 = rf(table)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// Fields provides a mock function with given fields: fields
func (_m *SqlQuery) Fields(fields ...string) sqlBuilder.SqlQuery {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(...string) sqlBuilder.SqlQuery); ok {
		r0 = rf(fields...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// From provides a mock function with given fields: tables
func (_m *SqlQuery) From(tables ...string) sqlBuilder.SqlQuery {
	_va := make([]interface{}, len(tables))
	for _i := range tables {
		_va[_i] = tables[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(...string) sqlBuilder.SqlQuery); ok {
		r0 = rf(tables...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// InsertInto provides a mock function with given fields: table
func (_m *SqlQuery) InsertInto(table string) sqlBuilder.SqlQuery {
	ret := _m.Called(table)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(string) sqlBuilder.SqlQuery); ok {
		r0 = rf(table)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// Join provides a mock function with given fields: tables
func (_m *SqlQuery) Join(tables ...string) sqlBuilder.SqlQuery {
	_va := make([]interface{}, len(tables))
	for _i := range tables {
		_va[_i] = tables[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(...string) sqlBuilder.SqlQuery); ok {
		r0 = rf(tables...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// Limit provides a mock function with given fields: limit
func (_m *SqlQuery) Limit(limit string) sqlBuilder.SqlQuery {
	ret := _m.Called(limit)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(string) sqlBuilder.SqlQuery); ok {
		r0 = rf(limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// Offset provides a mock function with given fields: offset
func (_m *SqlQuery) Offset(offset string) sqlBuilder.SqlQuery {
	ret := _m.Called(offset)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(string) sqlBuilder.SqlQuery); ok {
		r0 = rf(offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// On provides a mock function with given fields: conditions
func (_m *SqlQuery) On(conditions ...string) sqlBuilder.SqlQuery {
	_va := make([]interface{}, len(conditions))
	for _i := range conditions {
		_va[_i] = conditions[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(...string) sqlBuilder.SqlQuery); ok {
		r0 = rf(conditions...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// OrderBy provides a mock function with given fields: columns
func (_m *SqlQuery) OrderBy(columns ...string) sqlBuilder.SqlQuery {
	_va := make([]interface{}, len(columns))
	for _i := range columns {
		_va[_i] = columns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(...string) sqlBuilder.SqlQuery); ok {
		r0 = rf(columns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// Select provides a mock function with given fields: columns
func (_m *SqlQuery) Select(columns ...string) sqlBuilder.SqlQuery {
	_va := make([]interface{}, len(columns))
	for _i := range columns {
		_va[_i] = columns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(...string) sqlBuilder.SqlQuery); ok {
		r0 = rf(columns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// Set provides a mock function with given fields: values
func (_m *SqlQuery) Set(values ...string) sqlBuilder.SqlQuery {
	_va := make([]interface{}, len(values))
	for _i := range values {
		_va[_i] = values[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(...string) sqlBuilder.SqlQuery); ok {
		r0 = rf(values...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// Update provides a mock function with given fields: table
func (_m *SqlQuery) Update(table string) sqlBuilder.SqlQuery {
	ret := _m.Called(table)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(string) sqlBuilder.SqlQuery); ok {
		r0 = rf(table)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// Values provides a mock function with given fields: values
func (_m *SqlQuery) Values(values ...string) sqlBuilder.SqlQuery {
	_va := make([]interface{}, len(values))
	for _i := range values {
		_va[_i] = values[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(...string) sqlBuilder.SqlQuery); ok {
		r0 = rf(values...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

// Where provides a mock function with given fields: conditions
func (_m *SqlQuery) Where(conditions ...string) sqlBuilder.SqlQuery {
	_va := make([]interface{}, len(conditions))
	for _i := range conditions {
		_va[_i] = conditions[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 sqlBuilder.SqlQuery
	if rf, ok := ret.Get(0).(func(...string) sqlBuilder.SqlQuery); ok {
		r0 = rf(conditions...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sqlBuilder.SqlQuery)
		}
	}

	return r0
}

type mockConstructorTestingTNewSqlQuery interface {
	mock.TestingT
	Cleanup(func())
}

// NewSqlQuery creates a new instance of SqlQuery. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSqlQuery(t mockConstructorTestingTNewSqlQuery) *SqlQuery {
	mock := &SqlQuery{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}