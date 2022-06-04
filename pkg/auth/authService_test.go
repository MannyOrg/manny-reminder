package auth

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"log"
	"manny-reminder/pkg/models"
	"testing"
)

type RepositoryMock struct {
	Users []models.User
}

func (r RepositoryMock) GetUsers() ([]models.User, error) {
	return r.Users, nil
}

func (r RepositoryMock) GetUser(userId string) (models.User, error) {
	for _, user := range r.Users {
		if user.Id.String() == userId {
			return user, nil
		}
	}

	return models.User{}, nil
}

func (r RepositoryMock) AddUser(string, string) error {
	return nil
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}

func TestGetUsers_EmptyResponse(t *testing.T) {
	as, _ := getService()

	users, err := as.GetUsers()

	assert.Nil(t, err)
	assert.Empty(t, users)
}

func TestGetUsers_OneUser(t *testing.T) {
	as, r := getService()

	r.Users = []models.User{
		{},
	}

	users, err := as.GetUsers()

	assert.Nil(t, err)
	assert.NotEmpty(t, users)
	assert.Exactly(t, 1, len(users))
}

func TestGetUsers_MultipleUsers(t *testing.T) {
	as, r := getService()

	r.Users = []models.User{
		{},
		{},
		{},
		{},
	}

	users, err := as.GetUsers()

	assert.Nil(t, err)
	assert.NotEmpty(t, users)
	assert.Exactly(t, 4, len(users))
}

func getService() (*Service, *RepositoryMock) {
	l := log.Default()
	r := NewRepositoryMock()
	c := &oauth2.Config{}
	as := NewService(l, r, c)
	return as, r
}
