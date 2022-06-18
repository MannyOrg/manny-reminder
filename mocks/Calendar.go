// Code generated by mockery v2.13.0-beta.1. DO NOT EDIT.

package mocks

import (
	context "context"
	models "manny-reminder/internal/models"

	mock "github.com/stretchr/testify/mock"

	oauth2 "golang.org/x/oauth2"
)

// Calendar is an autogenerated mock type for the Calendar type
type Calendar struct {
	mock.Mock
}

// GetEventsForUser provides a mock function with given fields: ctx, tok, nextPageToken, size
func (_m *Calendar) GetEventsForUser(ctx context.Context, tok oauth2.Token, nextPageToken string, size int) (*models.Events, string, error) {
	ret := _m.Called(ctx, tok, nextPageToken, size)

	var r0 *models.Events
	if rf, ok := ret.Get(0).(func(context.Context, oauth2.Token, string, int) *models.Events); ok {
		r0 = rf(ctx, tok, nextPageToken, size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Events)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, oauth2.Token, string, int) string); ok {
		r1 = rf(ctx, tok, nextPageToken, size)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, oauth2.Token, string, int) error); ok {
		r2 = rf(ctx, tok, nextPageToken, size)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type NewCalendarT interface {
	mock.TestingT
	Cleanup(func())
}

// NewCalendar creates a new instance of Calendar. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCalendar(t NewCalendarT) *Calendar {
	mock := &Calendar{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}