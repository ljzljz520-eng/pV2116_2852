package store

import (
	"fmt"
	"stickerchallenge/internal/domain"
)

type IntegrityIssue struct {
	Bucket  string
	Key     string
	Message string
}

func (s *Store) CheckIntegrity() ([]IntegrityIssue, error) {
	issues := make([]IntegrityIssue, 0)
	batches, err := s.ListBatches()
	if err != nil {
		return issues, err
	}
	for _, batch := range batches {
		if err := domain.ValidateBatch(batch); err != nil {
			issues = append(issues, IntegrityIssue{Bucket: "batches", Key: batch.ID, Message: err.Error()})
		}
		for _, record := range batch.Records {
			if record.BatchID != batch.ID {
				issues = append(issues, IntegrityIssue{Bucket: "batches", Key: batch.ID, Message: fmt.Sprintf("record %s has wrong batch", record.ID)})
			}
		}
	}
	return issues, nil
}
func (s *Store) RepairBatch(batch domain.Batch) error {
	if err := domain.ValidateBatch(batch); err != nil {
		return err
	}
	return s.PutBatch(batch)
}
func (s *Store) Flush() error { return s.db.Sync() }
func (s *Store) PathInfo() string {
	if s == nil || s.db == nil {
		return "closed"
	}
	return "open"
}
func (s *Store) IsClosed() bool { return s == nil || s.db == nil }
