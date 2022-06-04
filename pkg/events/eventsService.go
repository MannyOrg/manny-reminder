package events

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"manny-reminder/pkg/auth"
	"manny-reminder/pkg/models"
	"time"
)

type IService interface {
	GetUsersEvents() ([]models.Event, error)
}

type Service struct {
	l      *log.Logger
	r      IRepository
	config *oauth2.Config
	as     auth.IService
}

func NewEvents(r *Repository, l *log.Logger, config *oauth2.Config, as *auth.Service) *Service {
	return &Service{l, r, config, as}
}

func (e Service) AddUser(authCode string) error {
	tok, err := e.config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	userId := uuid.NewString()
	err = e.as.AddUser(userId)
	if err != nil {
		return err
	}

	e.as.SaveToken("users/"+userId+".json", tok)
	return nil
}

func (e Service) GetUsersEvents() ([]models.Event, error) {
	var response []models.Event
	files, err := ioutil.ReadDir("users/")
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	for _, file := range files {
		client := e.as.GetClient(file.Name())

		srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return nil, err
		}

		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			return nil, err
		}

		if len(events.Items) != 0 {
			for _, item := range events.Items {
				date := item.Start.DateTime
				if date == "" {
					date = item.Start.Date
				}
				event := &models.Event{Title: item.Summary}
				response = append(response, *event)
			}
		}
	}
	return response, nil
}
