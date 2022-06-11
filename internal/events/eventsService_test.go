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

	userToken := "{\"access_token\":\"test\",\"token_type\":\"Bearer\",\"refresh_token\":\"test\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	uuid1, _ := uuid.NewUUID()
	uuid2, _ := uuid.NewUUID()

	users := []models.User{
		{Id: &uuid1, Token: &userToken},
		{Id: &uuid2, Token: &userToken},
	}
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

	userToken1 := "{\"access_token\":\"test1\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	userToken2 := "{\"access_token\":\"test2\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	userToken3 := "{\"access_token\":\"test3\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	uuid3 := uuid.New()
	users := []models.User{
		{Id: &uuid1, Token: &userToken1},
		{Id: &uuid2, Token: &userToken2},
		{Id: &uuid3, Token: &userToken3},
	}
	mockAuthServiceGetUsers(as, users, nil)

	mockedEvents := make(map[string]models.Events)
	user1Events := generateEvents("1", 3)
	mockedEvents[userToken1] = user1Events
	user2Events := generateEvents("2", 2)
	mockedEvents[userToken2] = user2Events
	user3Events := generateEvents("3", 0)
	mockedEvents[userToken3] = user3Events

	mockCalendarGetEventsForUser(c, mockedEvents, nil)

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
	_, as, _, es := initService(t)

	userToken := "invalid-token-obs"
	uuid1, _ := uuid.NewUUID()
	uuid2, _ := uuid.NewUUID()
	users := []models.User{
		{Id: &uuid1, Token: &userToken},
		{Id: &uuid2, Token: &userToken},
	}
	mockedEvents := make(map[string]models.Events)
	user1Events := generateEvents("1", 0)
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

	userToken := "{\"access_token\":\"test\",\"token_type\":\"Bearer\",\"refresh_token\":\"test\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	uuid1, _ := uuid.NewUUID()

	users := []models.User{{Id: &uuid1, Token: &userToken}}
	mockAuthServiceGetUser(as, &(users[0]), nil)
	mockCalendarGetEventsForUser(c, make(map[string]models.Events), nil)

	events, err := es.GetUserEvents(uuid.New().String(), "", 10)

	assert.Nil(t, err)
	assert.Empty(t, events)
}

func TestService_GetUserEvents_UserEvents(t *testing.T) {
	_, as, c, es := initService(t)

	userToken := "{\"access_token\":\"test\",\"token_type\":\"Bearer\",\"refresh_token\":\"test\",\"expiry\":\"2022-06-04T17:36:36.65039945+03:00\"}"
	uuid1, _ := uuid.NewUUID()

	users := []models.User{{Id: &uuid1, Token: &userToken}}

	var mockedEvents = make(map[string]models.Events)
	user1Events := generateEvents("1", 3)
	mockedEvents[userToken] = user1Events

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
