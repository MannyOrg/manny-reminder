package auth

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"log"
	"manny-reminder/internal/models"
	"manny-reminder/mocks"
	"testing"
)

func TestGetUsers_EmptyResponse(t *testing.T) {
	as, _ := getService(t)

	users, err := as.GetUsers()

	assert.Nil(t, err)
	assert.Empty(t, users)
}

func TestGetUsers_OneUser(t *testing.T) {
	as, _ := getService(t)

	users := []models.User{
		{},
	}

	users, err := as.GetUsers()

	assert.Nil(t, err)
	assert.NotEmpty(t, users)
	assert.Exactly(t, 1, len(users))
}

func TestGetUsers_MultipleUsers(t *testing.T) {
	as, _ := getService(t)

	users := []models.User{
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

func getService(t *testing.T) (*ServiceImpl, *mocks.AuthRepository) {
	l := log.Default()
	r := mocks.NewAuthRepository(t)
	c := &oauth2.Config{}
	as := NewService(l, r, c)
	return as, r
}
