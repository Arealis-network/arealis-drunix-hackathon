package repository

import (
	"sync"

	"arealis-drunix-hackathon/models"
)

type CorrelationStore struct {
	mu     sync.RWMutex
	events map[string][]models.AgentEvent
}

func NewCorrelationStore() *CorrelationStore {
	return &CorrelationStore{
		events: make(map[string][]models.AgentEvent),
	}
}

func (s *CorrelationStore) Add(event models.AgentEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.CorrelationID] = append(
		s.events[event.CorrelationID],
		event,
	)
}

func (s *CorrelationStore) Get(correlationID string) []models.AgentEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := s.events[correlationID]

	result := make([]models.AgentEvent, len(events))
	copy(result, events)

	return result
}