// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	volume "github.com/docker/go-plugins-helpers/volume"
	mock "github.com/stretchr/testify/mock"
)

// VolumeService is an autogenerated mock type for the VolumeService type
type VolumeService struct {
	mock.Mock
}

// Capabilities provides a mock function with given fields:
func (_m *VolumeService) Capabilities() volume.Capability {
	ret := _m.Called()

	var r0 volume.Capability
	if rf, ok := ret.Get(0).(func() volume.Capability); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(volume.Capability)
	}

	return r0
}

// Create provides a mock function with given fields: name, opt
func (_m *VolumeService) Create(name string, opt map[string]string) error {
	ret := _m.Called(name, opt)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, map[string]string) error); ok {
		r0 = rf(name, opt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: name
func (_m *VolumeService) Get(name string) (*volume.Volume, error) {
	ret := _m.Called(name)

	var r0 *volume.Volume
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*volume.Volume, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *volume.Volume); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*volume.Volume)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields:
func (_m *VolumeService) List() ([]*volume.Volume, error) {
	ret := _m.Called()

	var r0 []*volume.Volume
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*volume.Volume, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*volume.Volume); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*volume.Volume)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Mount provides a mock function with given fields: id, name
func (_m *VolumeService) Mount(id string, name string) (string, error) {
	ret := _m.Called(id, name)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(id, name)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(id, name)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Path provides a mock function with given fields: name
func (_m *VolumeService) Path(name string) (string, error) {
	ret := _m.Called(name)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: name
func (_m *VolumeService) Remove(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unmount provides a mock function with given fields: id, name
func (_m *VolumeService) Unmount(id string, name string) error {
	ret := _m.Called(id, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(id, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewVolumeService creates a new instance of VolumeService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewVolumeService(t interface {
	mock.TestingT
	Cleanup(func())
}) *VolumeService {
	mock := &VolumeService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
