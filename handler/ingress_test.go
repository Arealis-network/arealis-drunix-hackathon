package handler

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseEvent(t *testing.T) {
	raw, err := os.ReadFile("../examples/solar-asset/events.json")
	if err != nil {
		t.Fatalf("failed to read events.json: %v", err)
	}

	var rawEvents []json.RawMessage

	if err := json.Unmarshal(raw, &rawEvents); err != nil {
		t.Fatalf("failed to parse events.json: %v", err)
	}

	if len(rawEvents) != 3 {
		t.Fatalf("expected 3 events, got %d", len(rawEvents))
	}

	for _, rawEvent := range rawEvents {
		event, err := ParseEvent(rawEvent)
		if err != nil {
			t.Fatalf("failed to parse event: %v", err)
		}

		if event.CorrelationID != "txn-001" {
			t.Errorf("unexpected correlation ID: %s", event.CorrelationID)
		}

		if event.AssetID != "SOLAR-001" {
			t.Errorf("unexpected asset ID: %s", event.AssetID)
		}
	}
}
