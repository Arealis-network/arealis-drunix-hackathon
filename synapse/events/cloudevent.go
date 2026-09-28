package events

import "time"

type CloudEvent struct {
	SpecVersion   string    `json:"specversion"`
	ID            string    `json:"id"`
	Source        string    `json:"source"`
	Type          string    `json:"type"`
	Time          time.Time `json:"time"`
	CorrelationID string    `json:"correlationid"`
	Data          any       `json:"data"`
}
