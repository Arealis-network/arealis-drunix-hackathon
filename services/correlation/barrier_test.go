package correlation

import (
	"testing"

	"arealis-drunix-hackathon/models"
	"arealis-drunix-hackathon/repository"
	"arealis-drunix-hackathon/schemas"
)

func TestBarrier(t *testing.T) {
	store := repository.NewCorrelationStore()
	barrier := NewBarrier(store)

	correlationID := "txn-001"

	valuation := models.AgentEvent{
		ID:            "evt-001",
		Source:        schemas.AgentValuation,
		Type:          schemas.EventValuationCompleted,
		CorrelationID: correlationID,
		AssetID:       "SOLAR-001",
	}

	compliance := models.AgentEvent{
		ID:            "evt-002",
		Source:        schemas.AgentCompliance,
		Type:          schemas.EventComplianceCompleted,
		CorrelationID: correlationID,
		AssetID:       "SOLAR-001",
	}

	treasury := models.AgentEvent{
		ID:            "evt-003",
		Source:        schemas.AgentTreasury,
		Type:          schemas.EventFundsVerified,
		CorrelationID: correlationID,
		AssetID:       "SOLAR-001",
	}

	if barrier.AddEvent(valuation) {
		t.Fatal("barrier should not be ready after valuation only")
	}

	if barrier.AddEvent(compliance) {
		t.Fatal("barrier should not be ready after valuation + compliance")
	}

	if !barrier.AddEvent(treasury) {
		t.Fatal("barrier should be ready after all required events")
	}
}
