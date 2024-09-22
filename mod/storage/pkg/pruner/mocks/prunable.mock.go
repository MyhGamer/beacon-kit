// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Prunable is an autogenerated mock type for the Prunable type
type Prunable struct {
	mock.Mock
}

type Prunable_Expecter struct {
	mock *mock.Mock
}

func (_m *Prunable) EXPECT() *Prunable_Expecter {
	return &Prunable_Expecter{mock: &_m.Mock}
}

// Prune provides a mock function with given fields: start, end
func (_m *Prunable) Prune(start uint64, end uint64) error {
	ret := _m.Called(start, end)

	if len(ret) == 0 {
		panic("no return value specified for Prune")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, uint64) error); ok {
		r0 = rf(start, end)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Prunable_Prune_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Prune'
type Prunable_Prune_Call struct {
	*mock.Call
}

// Prune is a helper method to define mock.On call
//   - start uint64
//   - end uint64
func (_e *Prunable_Expecter) Prune(start interface{}, end interface{}) *Prunable_Prune_Call {
	return &Prunable_Prune_Call{Call: _e.mock.On("Prune", start, end)}
}

func (_c *Prunable_Prune_Call) Run(run func(start uint64, end uint64)) *Prunable_Prune_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint64), args[1].(uint64))
	})
	return _c
}

func (_c *Prunable_Prune_Call) Return(_a0 error) *Prunable_Prune_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Prunable_Prune_Call) RunAndReturn(run func(uint64, uint64) error) *Prunable_Prune_Call {
	_c.Call.Return(run)
	return _c
}

// NewPrunable creates a new instance of Prunable. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPrunable(t interface {
	mock.TestingT
	Cleanup(func())
}) *Prunable {
	mock := &Prunable{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
