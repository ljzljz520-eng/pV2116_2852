package service

import (
	"fmt"
	"stickerchallenge/internal/domain"
)

func (s *Service) RegisterMany(items []domain.Batch) (int, error) {
	count := 0
	for _, item := range items {
		if err := s.Store.PutBatch(item); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
func (s *Service) RecalculateBatch(batchID string) (domain.Batch, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	for index := range batch.Records {
		updated, updateErr := domain.Recalculate(batch.Records[index], s.Divisors)
		if updateErr != nil {
			return batch, updateErr
		}
		batch.Records[index] = updated
	}
	if err := s.Store.UpdateBatch(batch, batch.Version); err != nil {
		return batch, err
	}
	return batch, nil
}
func (s *Service) CopyBatch(sourceID, targetID, owner string) (domain.Batch, error) {
	source, err := s.Store.GetBatch(sourceID)
	if err != nil {
		return domain.Batch{}, err
	}
	copy := domain.CloneBatch(source)
	copy.ID = targetID
	copy.Label = source.Label + " copy"
	copy.Owner = owner
	copy.Status = domain.StatusRegistered
	copy.Version = 1
	for index := range copy.Records {
		copy.Records[index].BatchID = targetID
		copy.Records[index].ID = fmt.Sprintf("%s-%d", targetID, index+1)
		copy.Records[index].Confirmed = false
	}
	if err := s.Store.PutBatch(copy); err != nil {
		return domain.Batch{}, err
	}
	return copy, nil
}
func (s *Service) ArchiveIfReady(batchID, actor string) (domain.Batch, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	if batch.Status == domain.StatusConfirmed {
		if _, err = s.PublishBatch(batchID, actor); err != nil {
			return batch, err
		}
	}
	return s.ArchiveBatch(batchID, actor)
}
func (s *Service) ValidateBatch(batchID string) ([]domain.ValidationIssue, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return nil, err
	}
	return domain.DefaultPolicy().CheckBatch(batch), nil
}
