package calendar

import (
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"manny-reminder/internal/models"
	"time"
)

type Calendar interface {
	GetEventsForUser(ctx context.Context, tok oauth2.Token, nextPageToken string, size int) (*models.Events, string, error)
}

type GoogleCalendar struct {
	config *oauth2.Config
}

func NewCalendar(c *oauth2.Config) *GoogleCalendar {
	return &GoogleCalendar{config: c}
}

func (c GoogleCalendar) GetEventsForUser(ctx context.Context, tok oauth2.Token, pageToken string, size int) (*models.Events, string, error) {
	client := c.config.Client(context.Background(), &tok)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, "", err
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.
		List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(int64(size)).
		PageToken(pageToken).
		OrderBy("startTime").
		Do()

	if err != nil {
		return nil, "", err
	}

	var result models.Events
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

	return &result, events.NextPageToken, nil
}
