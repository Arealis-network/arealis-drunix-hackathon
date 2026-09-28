package orchestration

import (
	"fmt"

	"arealis-drunix-hackathon/models"
	"arealis-drunix-hackathon/schemas"
	"arealis-drunix-hackathon/services/correlation"
)

type Service struct {
	barrier *correlation.Barrier
}

func NewService(barrier *correlation.Barrier) *Service {
	return &Service{
		barrier: barrier,
	}
}

func (s *Service) ProcessEvent(event models.AgentEvent) (*models.TransactionIntent, bool, error) {
	ready := s.barrier.AddEvent(event)

	if !ready {
		return nil, false, nil
	}

	events := s.barrier.GetEvents(event.CorrelationID)

	assetID, err := validateEvents(events)
	if err != nil {
		return nil, false, err
	}

	intent := &models.TransactionIntent{
		CorrelationID: event.CorrelationID,
		AssetID:       assetID,
		Action:        schemas.ActionLockAsset,
	}

	return intent, true, nil
}

func validateEvents(events []models.AgentEvent) (string, error) {
	if len(events) == 0 {
		return "", fmt.Errorf("no events found")
	}

	assetID := events[0].AssetID

	for _, event := range events {
		if event.AssetID != assetID {
			return "", fmt.Errorf(
				"asset mismatch: expected %s, got %s",
				assetID,
				event.AssetID,
			)
		}
	}

	return assetID, nil
}
