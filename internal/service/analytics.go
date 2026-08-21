package service

import (
	"sort"
	"stickerchallenge/internal/domain"
)

type OperatorStats struct {
	Actor   string
	Events  int
	Batches int
	Notes   int
	Exports int
}
type BatchTimeline struct {
	BatchID  string
	Actions  []string
	Actors   []string
	Complete bool
}

func (s *Service) OperatorStats(batchID string) (OperatorStats, error) {
	events, err := s.Store.ListEvents(batchID)
	if err != nil {
		return OperatorStats{}, err
	}
	notes, err := s.Store.ListNotes(batchID)
	if err != nil {
		return OperatorStats{}, err
	}
	exports, err := s.Store.ListExports(batchID)
	if err != nil {
		return OperatorStats{}, err
	}
	stats := OperatorStats{Events: len(events), Batches: 1, Notes: len(notes), Exports: len(exports)}
	actors := map[string]int{}
	for _, event := range events {
		actors[event.Actor]++
	}
	for actor, count := range actors {
		if stats.Actor == "" || count > actors[stats.Actor] || (count == actors[stats.Actor] && actor < stats.Actor) {
			stats.Actor = actor
		}
	}
	return stats, nil
}

func (s *Service) Timeline(batchID string) (BatchTimeline, error) {
	events, err := s.Store.ListEvents(batchID)
	if err != nil {
		return BatchTimeline{}, err
	}
	timeline := BatchTimeline{BatchID: batchID, Actions: make([]string, 0), Actors: make([]string, 0)}
	seenActors := map[string]bool{}
	for _, event := range events {
		timeline.Actions = append(timeline.Actions, event.Action)
		if !seenActors[event.Actor] {
			timeline.Actors = append(timeline.Actors, event.Actor)
			seenActors[event.Actor] = true
		}
	}
	sort.Strings(timeline.Actors)
	timeline.Complete = hasAction(timeline.Actions, "register") && hasAction(timeline.Actions, "review") && hasAction(timeline.Actions, "confirm")
	return timeline, nil
}

func hasAction(actions []string, target string) bool {
	for _, action := range actions {
		if action == target {
			return true
		}
	}
	return false
}
func (s *Service) ConfirmedCount(batchID string) (int, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return 0, err
	}
	_, _, confirmed := domain.CountResults(batch.Records)
	return confirmed, nil
}
func (s *Service) PassingCount(batchID string) (int, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return 0, err
	}
	passing, _, _ := domain.CountResults(batch.Records)
	return passing, nil
}
func (s *Service) FailingCount(batchID string) (int, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return 0, err
	}
	_, failing, _ := domain.CountResults(batch.Records)
	return failing, nil
}
func (s *Service) BatchReadyForArchive(batchID string) (bool, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return false, err
	}
	return batch.Status == domain.StatusPublished && len(batch.Records) > 0, nil
}
func (s *Service) Labels(status domain.Status) ([]string, error) {
	batches, err := s.Search(domain.SearchQuery{Status: status})
	if err != nil {
		return nil, err
	}
	labels := make([]string, 0, len(batches))
	for _, batch := range batches {
		labels = append(labels, batch.Label)
	}
	sort.Strings(labels)
	return labels, nil
}
