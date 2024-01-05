// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/raffaele-pilloni/axxon-test/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// TaskRepositoryInterface is an autogenerated mock type for the TaskRepositoryInterface type
type TaskRepositoryInterface struct {
	mock.Mock
}

// FindTaskByID provides a mock function with given fields: ctx, taskID
func (_m *TaskRepositoryInterface) FindTaskByID(ctx context.Context, taskID int) (*entity.Task, error) {
	ret := _m.Called(ctx, taskID)

	var r0 *entity.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*entity.Task, error)); ok {
		return rf(ctx, taskID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *entity.Task); ok {
		r0 = rf(ctx, taskID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, taskID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTaskRepositoryInterface creates a new instance of TaskRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskRepositoryInterface {
	mock := &TaskRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}