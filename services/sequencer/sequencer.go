package sequencer

import (
	"sync"

	"arealis-drunix-hackathon/models"
)

type Sequencer struct {
	mu     sync.Mutex
	queues map[string][]models.TransactionIntent
}

func NewSequencer() *Sequencer {
	return &Sequencer{
		queues: make(map[string][]models.TransactionIntent),
	}
}

func (s *Sequencer) Enqueue(intent models.TransactionIntent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.queues[intent.AssetID] = append(
		s.queues[intent.AssetID],
		intent,
	)
}

func (s *Sequencer) Next(assetID string) (models.TransactionIntent, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue := s.queues[assetID]

	if len(queue) == 0 {
		return models.TransactionIntent{}, false
	}

	next := queue[0]

	s.queues[assetID] = queue[1:]

	return next, true
}
