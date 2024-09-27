// Code generated by mockery v2.14.0. DO NOT EDIT.

package automock

import (
	internal "github.com/kyma-project/kyma-environment-broker/internal"
	mock "github.com/stretchr/testify/mock"
)

// ConfigurationProvider is an autogenerated mock type for the ConfigurationProvider type
type ConfigurationProvider struct {
	mock.Mock
}

// ProvideForGivenVersionAndPlan provides a mock function with given fields: kymaVersion, planName
func (_m *ConfigurationProvider) ProvideForGivenPlan(planName string) (*internal.ConfigForPlan, error) {
	ret := _m.Called(planName)

	var r0 *internal.ConfigForPlan
	if rf, ok := ret.Get(0).(func(string) *internal.ConfigForPlan); ok {
		r0 = rf(planName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.ConfigForPlan)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(planName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewConfigurationProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewConfigurationProvider creates a new instance of ConfigurationProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfigurationProvider(t mockConstructorTestingTNewConfigurationProvider) *ConfigurationProvider {
	mock := &ConfigurationProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
