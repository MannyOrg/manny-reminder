package models

type Event struct {
	Title     string   `json:"title"`
	Start     string   `json:"start"`
	End       string   `json:"end"`
	Organizer string   `json:"organizer"`
	Attendees []string `json:"attendees"`
}

type EventsResponse struct {
	Items         []Event `json:"items"`
	NextPageToken string  `json:"nextPageToken"`
}
