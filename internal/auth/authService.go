package auth

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"log"
	"manny-reminder/internal/models"
	"net/http"
	"os"
)

type AuthService interface {
	SaveUser(authCode string) error
	GetUsers() ([]models.User, error)
	GetTokenFromWeb() string
	GetClient(user string) *http.Client
	GetUser(id string) (*models.User, error)
}

type ServiceImpl struct {
	l      *log.Logger
	r      AuthRepository
	config *oauth2.Config
}

func NewService(l *log.Logger, r AuthRepository, config *oauth2.Config) *ServiceImpl {
	return &ServiceImpl{l, r, config}
}

func (s ServiceImpl) GetUsers() ([]models.User, error) {
	return s.r.GetUsers()
}

func (s ServiceImpl) GetUser(userId string) (*models.User, error) {
	return s.r.GetUser(userId)
}

// GetClient Retrieve a token, saves the token, then returns the generated client.
func (s *ServiceImpl) GetClient(user string) *http.Client {
	// The file credentials.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := s.tokenFromFile("users/" + user)
	if err != nil {
		s.l.Fatalf("Unable to read token file: %v", err)
	}
	return s.config.Client(context.Background(), tok)
}

// GetTokenFromWeb Request a token from the web, then returns the retrieved token.
func (s *ServiceImpl) GetTokenFromWeb() string {
	authURL := s.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

// Retrieves a token from a local file.
func (s *ServiceImpl) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			s.l.Fatalf(err.Error())
		}
	}(f)
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func (s ServiceImpl) SaveUser(authCode string) error {
	tok, err := s.config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	userId := uuid.NewString()

	ts, err := json.Marshal(tok)
	if err != nil {
		return err
	}

	err = s.r.AddUser(userId, string(ts))
	if err != nil {
		return err
	}

	return nil
}
