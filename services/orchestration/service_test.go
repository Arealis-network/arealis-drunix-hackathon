package orchestration

import (
	"testing"

	"arealis-drunix-hackathon/models"
	"arealis-drunix-hackathon/repository"
	"arealis-drunix-hackathon/schemas"
	"arealis-drunix-hackathon/services/correlation"
)

func TestProcessEvent(t *testing.T) {
	store := repository.NewCorrelationStore()
	barrier := correlation.NewBarrier(store)
	service := NewService(barrier)

	events := []models.AgentEvent{
		{
			ID:            "evt-001",
			Source:        schemas.AgentValuation,
			Type:          schemas.EventValuationCompleted,
			CorrelationID: "txn-001",
			AssetID:       "SOLAR-001",
		},
		{
			ID:            "evt-002",
			Source:        schemas.AgentCompliance,
			Type:          schemas.EventComplianceCompleted,
			CorrelationID: "txn-001",
			AssetID:       "SOLAR-001",
		},
		{
			ID:            "evt-003",
			Source:        schemas.AgentTreasury,
			Type:          schemas.EventFundsVerified,
			CorrelationID: "txn-001",
			AssetID:       "SOLAR-001",
		},
	}

	for i := 0; i < 2; i++ {
		intent, ready, err := service.ProcessEvent(events[i])

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if ready {
			t.Fatalf("transaction should not be ready after %d events", i+1)
		}

		if intent != nil {
			t.Fatal("intent should be nil before quorum")
		}
	}

	intent, ready, err := service.ProcessEvent(events[2])

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ready {
		t.Fatal("transaction should be ready after all required events")
	}

	if intent == nil {
		t.Fatal("expected transaction intent")
	}

	if intent.CorrelationID != "txn-001" {
		t.Errorf("unexpected correlation ID: %s", intent.CorrelationID)
	}

	if intent.AssetID != "SOLAR-001" {
		t.Errorf("unexpected asset ID: %s", intent.AssetID)
	}

	if intent.Action != schemas.ActionLockAsset {
		t.Errorf("unexpected action: %s", intent.Action)
	}
}
