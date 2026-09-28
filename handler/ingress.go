package handler

import (
	"encoding/json"
	"fmt"

	"arealis-drunix-hackathon/models"
)

func ParseEvent(raw []byte) (models.AgentEvent, error) {
	var event models.CloudEvent

	if err := json.Unmarshal(raw, &event); err != nil {
		return models.AgentEvent{}, fmt.Errorf("invalid CloudEvent: %w", err)
	}

	if event.SpecVersion == "" {
		return models.AgentEvent{}, fmt.Errorf("missing specversion")
	}

	if event.ID == "" {
		return models.AgentEvent{}, fmt.Errorf("missing event id")
	}

	if event.Source == "" {
		return models.AgentEvent{}, fmt.Errorf("missing event source")
	}

	if event.Type == "" {
		return models.AgentEvent{}, fmt.Errorf("missing event type")
	}

	if event.CorrelationID == "" {
		return models.AgentEvent{}, fmt.Errorf("missing correlationid")
	}

	var data map[string]any

	if err := json.Unmarshal(event.Data, &data); err != nil {
		return models.AgentEvent{}, fmt.Errorf("invalid event data: %w", err)
	}

	assetID, ok := data["asset_id"].(string)
	if !ok || assetID == "" {
		return models.AgentEvent{}, fmt.Errorf("missing asset_id")
	}

	return models.AgentEvent{
		ID:            event.ID,
		Source:        event.Source,
		Type:          event.Type,
		CorrelationID: event.CorrelationID,
		AssetID:       assetID,
		Data:          data,
	}, nil
}