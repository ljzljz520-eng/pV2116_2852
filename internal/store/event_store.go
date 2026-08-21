package store

import (
	"sort"
	"stickerchallenge/internal/domain"
)

var eventBucket = []byte("events")

func (s *Store) PutEvent(event domain.AuditEvent) error {
	if event.ID == "" || event.BatchID == "" || event.Action == "" {
		return domain.ErrInvalid
	}
	return s.Put(eventBucket, event.ID, event)
}

func (s *Store) ListEvents(batchID string) ([]domain.AuditEvent, error) {
	values, err := s.List(eventBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuditEvent, 0)
	for _, value := range values {
		var event domain.AuditEvent
		if err := decode(value, &event); err != nil {
			return nil, err
		}
		if event.BatchID == batchID {
			result = append(result, event)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}
