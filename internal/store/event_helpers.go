package store

import (
	"sort"
	"stickerchallenge/internal/domain"
)

func (s *Store) EventsForAction(batchID, action string) ([]domain.AuditEvent, error) {
	events, err := s.ListEvents(batchID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuditEvent, 0)
	for _, event := range events {
		if event.Action == action {
			result = append(result, event)
		}
	}
	return result, nil
}
func (s *Store) Actors(batchID string) ([]string, error) {
	events, err := s.ListEvents(batchID)
	if err != nil {
		return nil, err
	}
	values := map[string]bool{}
	for _, event := range events {
		values[event.Actor] = true
	}
	result := make([]string, 0, len(values))
	for actor := range values {
		result = append(result, actor)
	}
	sort.Strings(result)
	return result, nil
}
func (s *Store) EventCount(batchID string) (int, error) {
	events, err := s.ListEvents(batchID)
	return len(events), err
}
func (s *Store) LatestEvent(batchID string) (domain.AuditEvent, error) {
	events, err := s.ListEvents(batchID)
	if err != nil {
		return domain.AuditEvent{}, err
	}
	if len(events) == 0 {
		return domain.AuditEvent{}, domain.ErrNotFound
	}
	return events[len(events)-1], nil
}
