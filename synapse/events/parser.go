package events

import (
	"encoding/json"
	"fmt"
)

func Parse(raw []byte) (CloudEvent, error) {
	var event CloudEvent

	if err := json.Unmarshal(raw, &event); err != nil {
		return CloudEvent{}, fmt.Errorf("invalid CloudEvent JSON: %w", err)
	}

	if event.SpecVersion == "" {
		return CloudEvent{}, fmt.Errorf("missing specversion")
	}

	if event.ID == "" {
		return CloudEvent{}, fmt.Errorf("missing event id")
	}

	if event.Source == "" {
		return CloudEvent{}, fmt.Errorf("missing event source")
	}

	if event.Type == "" {
		return CloudEvent{}, fmt.Errorf("missing event type")
	}

	if event.CorrelationID == "" {
		return CloudEvent{}, fmt.Errorf("missing correlationid")
	}

	return event, nil
}
