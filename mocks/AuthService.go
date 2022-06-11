// Code generated by mockery v2.13.0-beta.1. DO NOT EDIT.

package mocks

import (
	models "manny-reminder/internal/models"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// AuthService is an autogenerated mock type for the AuthService type
type AuthService struct {
	mock.Mock
}

// GetClient provides a mock function with given fields: user
func (_m *AuthService) GetClient(user string) *http.Client {
	ret := _m.Called(user)

	var r0 *http.Client
	if rf, ok := ret.Get(0).(func(string) *http.Client); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Client)
		}
	}

	return r0
}

// GetTokenFromWeb provides a mock function with given fields:
func (_m *AuthService) GetTokenFromWeb() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetUser provides a mock function with given fields: id
func (_m *AuthService) GetUser(id string) (*models.User, error) {
	ret := _m.Called(id)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsers provides a mock function with given fields:
func (_m *AuthService) GetUsers() ([]models.User, error) {
	ret := _m.Called()

	var r0 []models.User
	if rf, ok := ret.Get(0).(func() []models.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveUser provides a mock function with given fields: authCode
func (_m *AuthService) SaveUser(authCode string) error {
	ret := _m.Called(authCode)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(authCode)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewAuthServiceT interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthService creates a new instance of AuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthService(t NewAuthServiceT) *AuthService {
	mock := &AuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
