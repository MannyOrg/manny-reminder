package events

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"log"
	"manny-reminder/pkg/models"
	"net/http"
	"testing"
)

const test_error_msg = "test error occured"

type EventsRepositoryMock struct {
}

type AuthServiceMock struct {
	Users      []models.User
	User       *models.User
	ThrowError bool
}

func (a AuthServiceMock) SaveUser(authCode string) error {
	//TODO implement me
	panic("implement me")
}

func (a AuthServiceMock) GetUsers() ([]models.User, error) {
	if a.ThrowError {
		return nil, errors.New(test_error_msg)
	}
	return a.Users, nil
}

func (a AuthServiceMock) GetTokenFromWeb() string {
	//TODO implement me
	panic("implement me")
}

func (a AuthServiceMock) GetClient(userId string) *http.Client {
	//TODO implement me
	panic("implement me")
}

func (a AuthServiceMock) GetUser(userId string) (*models.User, error) {
	return a.User, nil
}

type CalendarMock struct {
	Events map[string]*calendar.Events
}

func (c CalendarMock) GetEventsForUser(ctx context.Context, token oauth2.Token, nextPageToken string, size int) (*calendar.Events, error) {
	tokenStr, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	return c.Events[string(tokenStr)], nil
}

func TestService_GetUsersEvents_WhenNoUsers_Empty(t *testing.T) {
	er := EventsRepositoryMock{}
	as := AuthServiceMock{}
	c := CalendarMock{}
	es := NewService(er, log.Default(), as, c)

	events, err := es.GetUsersEvents("", 10)

	assert.Nil(t, err)
	assert.Empty(t, events)
}

func TestService_GetUsersEvents_WhenAsErr_Err(t *testing.T) {
	er := EventsRepositoryMock{}
	as := AuthServiceMock{}
	c := CalendarMock{}

	as.ThrowError = true

	es := NewService(er, log.Default(), as, c)

	events, err := es.GetUsersEvents("", 10)

	assert.Error(t, err)
	assert.Exactly(t, err.Error(), test_error_msg)
	assert.Nil(t, events)
}

func TestService_GetUsersEvents_WhenUsersAndNoEvents(t *testing.T) {
	er := EventsRepositoryMock{}
	as := AuthServiceMock{}
	c := CalendarMock{}

	userToken := "{\"access_token\":\"test\",\"token_type\":\"Bearer\",\"refresh_token\":\"test\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	uuid1, _ := uuid.NewUUID()
	uuid2, _ := uuid.NewUUID()
	as.Users = []models.User{
		{Id: &uuid1, Token: &userToken},
		{Id: &uuid2, Token: &userToken},
	}

	c.Events = make(map[string]*calendar.Events)

	es := NewService(er, log.Default(), as, c)

	events, err := es.GetUsersEvents("", 10)

	assert.Nil(t, err)
	assert.NotEmpty(t, events)
	assert.Exactly(t, 2, len(events))
	for _, event := range events {
		assert.Empty(t, event.Items)
		assert.Exactly(t, event.NextPageToken, "")
	}
}

func TestService_GetUsersEvents_WhenUsersAndEvents(t *testing.T) {
	er := EventsRepositoryMock{}
	as := AuthServiceMock{}
	c := CalendarMock{}

	userToken1 := "{\"access_token\":\"test1\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	userToken2 := "{\"access_token\":\"test2\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	userToken3 := "{\"access_token\":\"test3\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	uuid3 := uuid.New()
	as.Users = []models.User{
		{Id: &uuid1, Token: &userToken1},
		{Id: &uuid2, Token: &userToken2},
		{Id: &uuid3, Token: &userToken3},
	}

	c.Events = make(map[string]*calendar.Events)
	user1Events := []*calendar.Event{
		{Summary: "Meeting1", Start: &calendar.EventDateTime{}, End: &calendar.EventDateTime{}, Organizer: &calendar.EventOrganizer{Email: "test@test.com"}},
		{Summary: "Meeting2", Start: &calendar.EventDateTime{}, End: &calendar.EventDateTime{}, Organizer: &calendar.EventOrganizer{Email: "test@test.com"}},
		{Summary: "Meeting3", Start: &calendar.EventDateTime{}, End: &calendar.EventDateTime{}, Organizer: &calendar.EventOrganizer{Email: "test@test.com"}},
	}
	c.Events[userToken1] = &calendar.Events{
		Items: user1Events,
	}
	user2Events := []*calendar.Event{
		{Summary: "Meeting4", Start: &calendar.EventDateTime{}, End: &calendar.EventDateTime{}, Organizer: &calendar.EventOrganizer{Email: "test@test.com"}, Attendees: []*calendar.EventAttendee{&calendar.EventAttendee{Email: "test@email.com"}}},
		{Summary: "Meeting5", Start: &calendar.EventDateTime{}, End: &calendar.EventDateTime{}, Organizer: &calendar.EventOrganizer{Email: "test@test.com"}},
	}
	c.Events[userToken2] = &calendar.Events{
		Items: user2Events,
	}
	user3Events := []*calendar.Event{}
	c.Events[userToken3] = &calendar.Events{
		Items: user3Events,
	}

	es := NewService(er, log.Default(), as, c)

	events, err := es.GetUsersEvents("", 10)

	assert.Nil(t, err)
	assert.NotEmpty(t, events)
	assert.Exactly(t, 3, len(events))

	assert.NotEmpty(t, events[uuid1.String()].Items)
	assert.Exactly(t, 3, len(events[uuid1.String()].Items))
	assert.Exactly(t, events[uuid1.String()].NextPageToken, "")

	assert.NotEmpty(t, events[uuid2.String()].Items)
	assert.Exactly(t, 2, len(events[uuid2.String()].Items))
	assert.Exactly(t, events[uuid2.String()].NextPageToken, "")

	assert.Empty(t, events[uuid3.String()].Items)
	assert.Exactly(t, events[uuid3.String()].NextPageToken, "")
}

func TestService_GetUsersEvents_UserInvalidToken(t *testing.T) {
	er := EventsRepositoryMock{}
	as := AuthServiceMock{}
	c := CalendarMock{}

	userToken := "invalid-token-obs"
	uuid1, _ := uuid.NewUUID()
	uuid2, _ := uuid.NewUUID()
	as.Users = []models.User{
		{Id: &uuid1, Token: &userToken},
		{Id: &uuid2, Token: &userToken},
	}

	c.Events = make(map[string]*calendar.Events)

	es := NewService(er, log.Default(), as, c)

	events, err := es.GetUsersEvents("", 10)

	assert.NotNil(t, err)
	assert.Empty(t, events)
}

func TestService_GetUserEvents_UserDoesNotExist(t *testing.T) {
	er := EventsRepositoryMock{}
	as := AuthServiceMock{}
	c := CalendarMock{}

	as.User = nil

	es := NewService(er, log.Default(), as, c)

	events, err := es.GetUserEvents(uuid.New().String(), "", 10)

	assert.Nil(t, err)
	assert.Empty(t, events)
}

func TestService_GetUserEvents_NoUserEvents(t *testing.T) {
	er := EventsRepositoryMock{}
	as := AuthServiceMock{}
	c := CalendarMock{}

	userToken := "{\"access_token\":\"test\",\"token_type\":\"Bearer\",\"refresh_token\":\"test\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	uuid1, _ := uuid.NewUUID()

	as.User = &models.User{Id: &uuid1, Token: &userToken}

	es := NewService(er, log.Default(), as, c)

	events, err := es.GetUserEvents(uuid.New().String(), "", 10)

	assert.Nil(t, err)
	assert.Empty(t, events)
}

func TestService_GetUserEvents_UserEvents(t *testing.T) {
	er := EventsRepositoryMock{}
	as := AuthServiceMock{}
	c := CalendarMock{}

	userToken := "{\"access_token\":\"test\",\"token_type\":\"Bearer\",\"refresh_token\":\"test\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	uuid1, _ := uuid.NewUUID()

	as.User = &models.User{Id: &uuid1, Token: &userToken}

	c.Events = make(map[string]*calendar.Events)
	user1Events := []*calendar.Event{
		{Summary: "Meeting1", Start: &calendar.EventDateTime{}, End: &calendar.EventDateTime{}, Organizer: &calendar.EventOrganizer{Email: "test@test.com"}},
		{Summary: "Meeting2", Start: &calendar.EventDateTime{}, End: &calendar.EventDateTime{}, Organizer: &calendar.EventOrganizer{Email: "test@test.com"}},
		{Summary: "Meeting3", Start: &calendar.EventDateTime{}, End: &calendar.EventDateTime{}, Organizer: &calendar.EventOrganizer{Email: "test@test.com"}},
	}
	c.Events[userToken] = &calendar.Events{
		Items: user1Events,
	}

	es := NewService(er, log.Default(), as, c)

	events, err := es.GetUserEvents(uuid.New().String(), "", 10)

	assert.Nil(t, err)
	assert.NotEmpty(t, events)
	assert.Exactly(t, 3, len(events.Items))
}
