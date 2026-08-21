package service

import (
	"fmt"
	"stickerchallenge/internal/domain"
)

type MaintenanceResult struct {
	BatchID string
	Before  int
	After   int
	Changed bool
	Message string
}

func (s *Service) RefreshResults(batchID, actor string) (MaintenanceResult, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return MaintenanceResult{}, err
	}
	before := len(batch.Records)
	changed := false
	for index := range batch.Records {
		record := batch.Records[index]
		updated, updateErr := domain.Recalculate(record, s.Divisors)
		if updateErr != nil {
			return MaintenanceResult{}, updateErr
		}
		if !domain.CompareRecords(record, updated) {
			changed = true
		}
		updated.UpdatedBy = actor
		batch.Records[index] = updated
	}
	if changed {
		if err := s.Store.UpdateBatch(batch, batch.Version); err != nil {
			return MaintenanceResult{}, err
		}
		if err := s.audit(batchID, "refresh", actor, "results recalculated"); err != nil {
			return MaintenanceResult{}, err
		}
	}
	return MaintenanceResult{BatchID: batchID, Before: before, After: len(batch.Records), Changed: changed, Message: fmt.Sprintf("%d records refreshed", len(batch.Records))}, nil
}
func (s *Service) ReopenCheck(batchID string) error {
	if !s.Store.BatchExists(batchID) {
		return domain.ErrNotFound
	}
	_, err := s.Store.GetBatch(batchID)
	return err
}
func (s *Service) EnsureOwner(batchID, owner string) error {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return err
	}
	if batch.Owner != owner {
		return fmt.Errorf("owner mismatch: %w", domain.ErrInvalid)
	}
	return nil
}
func (s *Service) RenameBatch(batchID, label, actor string, expected int) (domain.Batch, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	if label == "" {
		return batch, domain.ErrInvalid
	}
	batch.Label = label
	batch.UpdatedAt = s.Clock.Now()
	if err := s.Store.UpdateBatch(batch, expected); err != nil {
		return batch, err
	}
	if err := s.audit(batchID, "rename", actor, label); err != nil {
		return batch, err
	}
	return batch, nil
}
func (s *Service) SetOwner(batchID, owner, actor string, expected int) (domain.Batch, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	if owner == "" {
		return batch, domain.ErrInvalid
	}
	batch.Owner = owner
	batch.UpdatedAt = s.Clock.Now()
	if err := s.Store.UpdateBatch(batch, expected); err != nil {
		return batch, err
	}
	if err := s.audit(batchID, "owner", actor, owner); err != nil {
		return batch, err
	}
	return batch, nil
}
func (s *Service) RemoveBatch(batchID, actor string) error {
	if !s.Store.BatchExists(batchID) {
		return domain.ErrNotFound
	}
	if err := s.Store.Delete([]byte("batches"), batchID); err != nil {
		return err
	}
	return s.audit(batchID, "remove", actor, "batch removed")
}
func (s *Service) RecordIDs(batchID string) ([]string, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(batch.Records))
	for _, record := range batch.Records {
		ids = append(ids, record.ID)
	}
	return ids, nil
}
