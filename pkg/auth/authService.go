package auth

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"log"
	"manny-reminder/pkg/models"
	"net/http"
	"os"
)

type IService interface {
	AddUser(authCode string) error
	GetUsers() ([]models.User, error)
	GetTokenFromWeb() string
	SaveToken(path string, token *oauth2.Token)
	GetClient(user string) *http.Client
}

type Service struct {
	l      *log.Logger
	r      *Repository
	config *oauth2.Config
}

func NewAuth(l *log.Logger, r *Repository, config *oauth2.Config) *Service {
	return &Service{l, r, config}
}

func (s Service) GetUsers() ([]models.User, error) {
	return s.r.GetUsers()
}

// GetClient Retrieve a token, saves the token, then returns the generated client.
func (s *Service) GetClient(user string) *http.Client {
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
func (s *Service) GetTokenFromWeb() string {
	authURL := s.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

// Retrieves a token from a local file.
func (s *Service) tokenFromFile(file string) (*oauth2.Token, error) {
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

// SaveToken Saves a token to a file path.
func (s *Service) SaveToken(path string, token *oauth2.Token) {
	s.l.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		s.l.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			s.l.Fatalf(err.Error())
		}
	}(f)
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		s.l.Fatalf(err.Error())
		return
	}
}

func (s *Service) AddUser(authCode string) error {
	return s.r.AddUser(authCode)
}
