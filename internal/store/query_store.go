package store

import (
	"sort"
	"stickerchallenge/internal/domain"
	"strings"
)

func (s *Store) FindBatches(label string, status domain.Status) ([]domain.Batch, error) {
	batches, err := s.ListBatches()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Batch, 0)
	for _, batch := range batches {
		if label != "" && !strings.Contains(strings.ToLower(batch.Label), strings.ToLower(label)) {
			continue
		}
		if status != "" && batch.Status != status {
			continue
		}
		result = append(result, batch)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}
func (s *Store) CountBatches(status domain.Status) (int, error) {
	batches, err := s.ListBatches()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, batch := range batches {
		if status == "" || batch.Status == status {
			count++
		}
	}
	return count, nil
}
func (s *Store) BatchExists(id string) bool { _, err := s.GetBatch(id); return err == nil }
func (s *Store) ReplaceRecord(batchID string, replacement domain.StickerRecord, expected int) error {
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return err
	}
	for index := range batch.Records {
		if batch.Records[index].ID == replacement.ID {
			batch.Records[index] = replacement
			return s.UpdateBatch(batch, expected)
		}
	}
	return domain.ErrNotFound
}
func (s *Store) AddRecord(batchID string, record domain.StickerRecord, expected int) error {
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return err
	}
	for _, existing := range batch.Records {
		if existing.ID == record.ID {
			return domain.ErrConflict
		}
	}
	batch.Records = append(batch.Records, record)
	return s.UpdateBatch(batch, expected)
}
func (s *Store) RemoveRecord(batchID, recordID string, expected int) error {
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return err
	}
	kept := make([]domain.StickerRecord, 0, len(batch.Records))
	found := false
	for _, record := range batch.Records {
		if record.ID == recordID {
			found = true
			continue
		}
		kept = append(kept, record)
	}
	if !found {
		return domain.ErrNotFound
	}
	batch.Records = kept
	return s.UpdateBatch(batch, expected)
}
