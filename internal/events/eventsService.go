package events

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	calendar2 "manny-reminder/internal/calendar"

	"log"
	"manny-reminder/internal/auth"
	"manny-reminder/internal/models"
)

type IService interface {
	GetUsersEvents(pageToken string, size int) (map[string]models.EventsResponse, error)
	GetUserEvents(userId string, pageToken string, size int) (models.EventsResponse, error)
}

type Service struct {
	l  *log.Logger
	r  IRepository
	as auth.IService
	c  calendar2.Calendar
}

func NewService(r IRepository, l *log.Logger, as auth.IService, c calendar2.Calendar) *Service {
	return &Service{l: l, r: r, as: as, c: c}
}

func (s Service) GetUsersEvents(pageToken string, size int) (map[string]models.EventsResponse, error) {
	response := make(map[string]models.EventsResponse)
	users, err := s.as.GetUsers()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return response, nil
	}
	ctx := context.Background()
	for _, user := range users {
		events, err := s.getUserEvents(ctx, &user, pageToken, size)
		if err != nil {
			return nil, err
		}

		response[user.Id.String()] = events
	}
	return response, nil
}

func (s Service) GetUserEvents(userId string, pageToken string, size int) (models.EventsResponse, error) {
	user, err := s.as.GetUser(userId)
	if err != nil {
		return models.EventsResponse{}, err
	}
	if user == nil {
		return models.EventsResponse{}, nil
	}
	ctx := context.Background()
	events, err := s.getUserEvents(ctx, user, pageToken, size)
	if err != nil {
		return models.EventsResponse{}, err
	}

	return events, nil
}

func (s Service) getUserEvents(ctx context.Context, user *models.User, pageToken string, size int) (models.EventsResponse, error) {
	var result []models.Event
	var tok oauth2.Token
	err := json.Unmarshal([]byte(*user.Token), &tok)
	if err != nil {
		return models.EventsResponse{}, err
	}

	events, npt, err := s.c.GetEventsForUser(ctx, tok, pageToken, size)

	if events == nil || len(*events) == 0 {
		return models.EventsResponse{Items: result, NextPageToken: ""}, nil
	}

	return models.EventsResponse{Items: result, NextPageToken: npt}, nil
}
