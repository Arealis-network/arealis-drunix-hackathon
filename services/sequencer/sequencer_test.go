package sequencer

import (
	"testing"

	"arealis-drunix-hackathon/models"
	"arealis-drunix-hackathon/schemas"
)

func TestSequencerPreservesOrderForSameAsset(t *testing.T) {
	sequencer := NewSequencer()

	tx1 := models.TransactionIntent{
		CorrelationID: "txn-001",
		AssetID:       "SOLAR-001",
		Action:        schemas.ActionLockAsset,
	}

	tx2 := models.TransactionIntent{
		CorrelationID: "txn-002",
		AssetID:       "SOLAR-001",
		Action:        schemas.ActionLockAsset,
	}

	sequencer.Enqueue(tx1)
	sequencer.Enqueue(tx2)

	first, ok := sequencer.Next("SOLAR-001")
	if !ok {
		t.Fatal("expected first transaction")
	}

	if first.CorrelationID != "txn-001" {
		t.Fatalf(
			"expected txn-001, got %s",
			first.CorrelationID,
		)
	}

	second, ok := sequencer.Next("SOLAR-001")
	if !ok {
		t.Fatal("expected second transaction")
	}

	if second.CorrelationID != "txn-002" {
		t.Fatalf(
			"expected txn-002, got %s",
			second.CorrelationID,
		)
	}
}

func TestSequencerKeepsAssetsIndependent(t *testing.T) {
	sequencer := NewSequencer()

	tx1 := models.TransactionIntent{
		CorrelationID: "txn-001",
		AssetID:       "SOLAR-001",
		Action:        schemas.ActionLockAsset,
	}

	tx2 := models.TransactionIntent{
		CorrelationID: "txn-002",
		AssetID:       "SOLAR-002",
		Action:        schemas.ActionLockAsset,
	}

	sequencer.Enqueue(tx1)
	sequencer.Enqueue(tx2)

	first, ok := sequencer.Next("SOLAR-002")
	if !ok {
		t.Fatal("expected transaction for SOLAR-002")
	}

	if first.CorrelationID != "txn-002" {
		t.Fatalf(
			"expected txn-002, got %s",
			first.CorrelationID,
		)
	}

	second, ok := sequencer.Next("SOLAR-001")
	if !ok {
		t.Fatal("expected transaction for SOLAR-001")
	}

	if second.CorrelationID != "txn-001" {
		t.Fatalf(
			"expected txn-001, got %s",
			second.CorrelationID,
		)
	}
}
