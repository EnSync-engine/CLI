package domain

import "time"

type Event struct {
	ID        string         `json:"id,omitempty"`
	Name      string         `json:"name"`
	Payload   map[string]any `json:"payload,omitempty"`
	CreatedAt time.Time      `json:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
}

type EventList struct {
	ResultsLength int      `json:"resultsLength"`
	Results       []*Event `json:"results"`
}

func (e *Event) IsZero() bool {
	return e.ID == "" && e.Name == ""
}
