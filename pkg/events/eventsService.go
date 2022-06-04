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
	GetUsersEvents(pageToken string, size int) (map[string]models.EventsResponse, error)
	GetUserEvents(userId string, pageToken string, size int) (models.EventsResponse, error)
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

func (s Service) GetUsersEvents(pageToken string, size int) (map[string]models.EventsResponse, error) {
	response := make(map[string]models.EventsResponse)
	users, err := s.as.GetUsers()
	if err != nil {
		return nil, err
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
	ctx := context.Background()
	events, err := s.getUserEvents(ctx, &user, pageToken, size)
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

	events, err := s.c.GetEventsForUser(ctx, tok, pageToken, size)

	if len(events.Items) != 0 {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			var attendees []string
			for _, attendee := range item.Attendees {
				attendees = append(attendees, attendee.Email)
			}
			event := &models.Event{
				Title:     item.Summary,
				Start:     item.Start.DateTime,
				End:       item.End.DateTime,
				Organizer: item.Organizer.Email,
				Attendees: attendees,
			}
			result = append(result, *event)
		}
	}

	return models.EventsResponse{Items: result, NextPageToken: events.NextPageToken}, nil
}
