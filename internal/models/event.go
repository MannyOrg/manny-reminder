package models

type Event struct {
	Title     string   `json:"title"`
	Start     string   `json:"start"`
	End       string   `json:"end"`
	Organizer string   `json:"organizer"`
	Attendees []string `json:"attendees"`
}

type Events []Event

type EventsResponse struct {
	Items         Events `json:"items"`
	NextPageToken string `json:"nextPageToken"`
}
