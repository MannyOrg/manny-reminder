package events

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	calendar2 "manny-reminder/pkg/calendar"

	"log"
	"manny-reminder/pkg/auth"
	"manny-reminder/pkg/models"
)

type IService interface {
	GetUsersEvents() (map[string][]models.Event, error)
	GetUserEvents(userId string) ([]models.Event, error)
}

type Service struct {
	l  *log.Logger
	r  IRepository
	as auth.IService
	c  calendar2.ICalendar
}

func NewEvents(r *Repository, l *log.Logger, as auth.IService, c calendar2.ICalendar) *Service {
	return &Service{l: l, r: r, as: as, c: c}
}

func (s Service) GetUsersEvents() (map[string][]models.Event, error) {
	response := make(map[string][]models.Event)
	users, err := s.as.GetUsers()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	for _, user := range users {
		events, err := s.getUserEvents(ctx, &user)
		if err != nil {
			return nil, err
		}

		response[user.Id.String()] = events
	}
	return response, nil
}

func (s Service) GetUserEvents(userId string) ([]models.Event, error) {
	user, err := s.as.GetUser(userId)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	events, err := s.getUserEvents(ctx, &user)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s Service) getUserEvents(ctx context.Context, user *models.User) ([]models.Event, error) {
	var result []models.Event
	var tok oauth2.Token
	err := json.Unmarshal([]byte(*user.Token), &tok)
	if err != nil {
		return nil, err
	}

	events, err := s.c.GetEventsForUser(ctx, tok)

	if len(events.Items) != 0 {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			event := &models.Event{Title: item.Summary}
			result = append(result, *event)
		}
	}

	return result, nil
}
