package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"log"
	"manny-reminder/internal/models"
	"manny-reminder/mocks"
	"strconv"
	"testing"
	"time"
)

const test_error_msg = "test error occured"

func TestService_GetUsersEvents_WhenNoUsers_Empty(t *testing.T) {
	_, as, _, es := initService(t)

	mockAuthServiceGetUsers(as, models.Users{}, nil)

	events, err := es.GetUsersEvents("", 10)

	assert.Nil(t, err)
	assert.Empty(t, events)
}

func TestService_GetUsersEvents_WhenAsErr_Err(t *testing.T) {
	_, as, _, es := initService(t)

	err := errors.New(test_error_msg)
	var users []models.User
	mockAuthServiceGetUsers(as, users, err)

	events, err := es.GetUsersEvents("", 10)

	assert.Error(t, err)
	assert.Exactly(t, err.Error(), test_error_msg)
	assert.Nil(t, events)
}

func TestService_GetUsersEvents_WhenUsersAndNoEvents(t *testing.T) {
	_, as, c, es := initService(t)

	users := generateUsers(2)
	mockedEvents := make(map[string]models.Events)
	mockAuthServiceGetUsers(as, users, nil)
	mockCalendarGetEventsForUser(c, mockedEvents, nil)

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
	_, as, c, es := initService(t)

	users := generateUsers(3)
	mockAuthServiceGetUsers(as, users, nil)

	mockedEvents := make(map[string]models.Events)
	user1Events := generateEvents("1", 3)
	mockedEvents[*users[0].Token] = user1Events
	user2Events := generateEvents("2", 2)
	mockedEvents[*users[1].Token] = user2Events
	user3Events := generateEvents("3", 0)
	mockedEvents[*users[2].Token] = user3Events

	mockCalendarGetEventsForUser(c, mockedEvents, nil)

	events, err := es.GetUsersEvents("", 10)

	assert.Nil(t, err)
	assert.NotEmpty(t, events)
	assert.Exactly(t, 3, len(events))

	assert.NotEmpty(t, events[users[0].Id.String()].Items)
	assert.Exactly(t, 3, len(events[users[0].Id.String()].Items))
	assert.Exactly(t, events[users[0].Id.String()].NextPageToken, "")

	assert.NotEmpty(t, events[users[1].Id.String()].Items)
	assert.Exactly(t, 2, len(events[users[1].Id.String()].Items))
	assert.Exactly(t, events[users[1].Id.String()].NextPageToken, "")

	assert.Empty(t, events[users[2].Id.String()].Items)
	assert.Exactly(t, events[users[2].Id.String()].NextPageToken, "")
}

func TestService_GetUsersEvents_UserInvalidToken(t *testing.T) {
	_, as, _, es := initService(t)

	userToken := "invalid-token-obs"
	users := generateUsers(2)
	mockedEvents := make(map[string]models.Events)
	user1Events := generateEvents("1", 0)
	users[0].Token = &userToken
	mockedEvents[userToken] = user1Events

	mockAuthServiceGetUsers(as, users, nil)

	events, err := es.GetUsersEvents("", 10)

	assert.NotNil(t, err)
	assert.Empty(t, events)
}

func TestService_GetUserEvents_UserDoesNotExist(t *testing.T) {
	_, as, _, es := initService(t)

	mockAuthServiceGetUser(as, nil, nil)

	events, err := es.GetUserEvents(uuid.New().String(), "", 10)

	assert.Nil(t, err)
	assert.Empty(t, events)
}

func TestService_GetUserEvents_NoUserEvents(t *testing.T) {
	_, as, c, es := initService(t)

	users := generateUsers(1)
	mockAuthServiceGetUser(as, &(users[0]), nil)
	mockCalendarGetEventsForUser(c, make(map[string]models.Events), nil)

	events, err := es.GetUserEvents(uuid.New().String(), "", 10)

	assert.Nil(t, err)
	assert.Empty(t, events)
}

func TestService_GetUserEvents_UserEvents(t *testing.T) {
	_, as, c, es := initService(t)

	users := generateUsers(1)
	var mockedEvents = make(map[string]models.Events)
	user1Events := generateEvents("1", 3)
	mockedEvents[*users[0].Token] = user1Events

	mockAuthServiceGetUser(as, &(users[0]), nil)
	mockCalendarGetEventsForUser(c, mockedEvents, nil)

	events, err := es.GetUserEvents(uuid.New().String(), "", 10)

	assert.Nil(t, err)
	assert.NotEmpty(t, events)
	assert.Exactly(t, 3, len(events.Items))
}

func initService(t *testing.T) (*mocks.EventsRepository, *mocks.AuthService, *mocks.Calendar, *ServiceImpl) {
	er := mocks.NewEventsRepository(t)
	as := mocks.NewAuthService(t)
	c := mocks.NewCalendar(t)
	es := NewService(er, log.Default(), as, c)
	return er, as, c, es
}

func mockAuthServiceGetUsers(as *mocks.AuthService, users []models.User, err error) {
	as.On("GetUsers").Return(users, err)
}

func mockAuthServiceGetUser(as *mocks.AuthService, user *models.User, err error) {
	as.On("GetUser", mock.Anything).Return(user, err)
}

func mockCalendarGetEventsForUser(c *mocks.Calendar, events map[string]models.Events, err error) {
	c.On("GetEventsForUser",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(
		func(_ context.Context, tok oauth2.Token, _ string, _ int) *models.Events {
			token, err := json.Marshal(tok)
			if err != nil {
				return nil
			}
			key := string(token)
			event := events[key]
			return &event
		},
		func(_ context.Context, _ oauth2.Token, _ string, _ int) string {
			return ""
		},
		func(_ context.Context, _ oauth2.Token, _ string, _ int) error {
			return err
		})
}

func generateEvents(prefix string, count int) models.Events {
	var events models.Events
	for i := 0; i < count; i++ {
		events = append(events, models.Event{
			Title:     fmt.Sprintf("Meeting %s%s", prefix, strconv.Itoa(i)),
			Organizer: fmt.Sprintf("Organizer %s%s", prefix, strconv.Itoa(i)),
		})
	}

	return events
}

func generateUsers(amount int) models.Users {
	var users models.Users
	for i := 0; i < amount; i++ {
		id, _ := uuid.NewUUID()
		userToken := generateUserToken(i, time.Now().Add(time.Hour*2))
		users = append(users, models.User{Id: &id, Token: &userToken})
	}
	return users
}

func generateUserToken(i int, expiry time.Time) string {
	expiryStr := time.Now().Add(time.Hour * 2).Format(time.RFC3339)
	token := "{\"access_token\":\"test %s\",\"token_type\":\"Bearer\",\"refresh_token\":\"test\",\"expiry\":\"%s\"}"
	userToken := fmt.Sprintf(token, strconv.Itoa(i+1), expiryStr)
	return userToken
}
