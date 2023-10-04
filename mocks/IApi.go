// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	entity "github.com/AlbertPuwadol/go-workshop/pkg/entity"
	mock "github.com/stretchr/testify/mock"
)

// IApi is an autogenerated mock type for the IApi type
type IApi struct {
	mock.Mock
}

// GetInfo provides a mock function with given fields: userid
func (_m *IApi) GetInfo(userid string) (*entity.Info, error) {
	ret := _m.Called(userid)

	var r0 *entity.Info
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*entity.Info, error)); ok {
		return rf(userid)
	}
	if rf, ok := ret.Get(0).(func(string) *entity.Info); ok {
		r0 = rf(userid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Info)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTimeline provides a mock function with given fields: hashtag, pagesize, cursor
func (_m *IApi) GetTimeline(hashtag string, pagesize int, cursor *string) (*entity.Timeline, error) {
	ret := _m.Called(hashtag, pagesize, cursor)

	var r0 *entity.Timeline
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int, *string) (*entity.Timeline, error)); ok {
		return rf(hashtag, pagesize, cursor)
	}
	if rf, ok := ret.Get(0).(func(string, int, *string) *entity.Timeline); ok {
		r0 = rf(hashtag, pagesize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Timeline)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int, *string) error); ok {
		r1 = rf(hashtag, pagesize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIApi creates a new instance of IApi. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIApi(t interface {
	mock.TestingT
	Cleanup(func())
}) *IApi {
	mock := &IApi{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
