package calendar

import (
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"time"
)

type ICalendar interface {
	GetEventsForUser(ctx context.Context, tok oauth2.Token) (*calendar.Events, error)
}

type Calendar struct {
	config *oauth2.Config
}

func NewCalendar(c *oauth2.Config) *Calendar {
	return &Calendar{config: c}
}

func (c Calendar) GetEventsForUser(ctx context.Context, tok oauth2.Token) (*calendar.Events, error) {
	client := c.config.Client(context.Background(), &tok)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.
		List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(10).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, err
	}

	return events, nil
}
