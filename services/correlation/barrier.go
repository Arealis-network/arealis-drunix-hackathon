package correlation

import (
	"arealis-drunix-hackathon/models"
	"arealis-drunix-hackathon/repository"
	"arealis-drunix-hackathon/schemas"
)

type Barrier struct {
	store *repository.CorrelationStore
}

func NewBarrier(store *repository.CorrelationStore) *Barrier {
	return &Barrier{
		store: store,
	}
}

func (b *Barrier) AddEvent(event models.AgentEvent) bool {
	b.store.Add(event)

	return b.IsReady(event.CorrelationID)
}

func (b *Barrier) IsReady(correlationID string) bool {
	events := b.store.Get(correlationID)

	required := map[string]bool{
		schemas.EventValuationCompleted:  false,
		schemas.EventComplianceCompleted: false,
		schemas.EventFundsVerified:       false,
	}

	for _, event := range events {
		if _, exists := required[event.Type]; exists {
			required[event.Type] = true
		}
	}

	for _, received := range required {
		if !received {
			return false
		}
	}

	return true
}
