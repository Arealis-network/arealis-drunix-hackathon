package models

type AgentEvent struct {
	ID            string
	Source        string
	Type          string
	CorrelationID string
	AssetID       string
	Data          map[string]any
}
