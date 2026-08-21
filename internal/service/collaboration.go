package service

import (
	"fmt"
	"sort"
	"stickerchallenge/internal/domain"
)

func (s *Service) Notes(batchID string) ([]domain.CollaborationNote, error) {
	return s.Store.ListNotes(batchID)
}
func (s *Service) AddReviewDecision(batchID, recordID, actor, reason string, accepted bool) (domain.Batch, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return batch, err
	}
	if batch.Status != domain.StatusReviewing {
		return batch, domain.ErrTransition
	}
	if err := domain.ApplyDecision(&batch, domain.ReviewDecision{RecordID: recordID, Reason: reason, Accepted: accepted}); err != nil {
		return batch, err
	}
	if err := s.Store.UpdateBatch(batch, batch.Version); err != nil {
		return batch, err
	}
	if err := s.audit(batchID, "decision", actor, fmt.Sprintf("record=%s accepted=%t", recordID, accepted)); err != nil {
		return batch, err
	}
	return batch, nil
}
func (s *Service) ReviewProgress(batchID string) (int, int, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return 0, 0, err
	}
	_, _, confirmed := domain.CountResults(batch.Records)
	return confirmed, len(batch.Records), nil
}
func (s *Service) CanPublish(batchID string) (bool, error) {
	batch, err := s.Store.GetBatch(batchID)
	if err != nil {
		return false, err
	}
	return (batch.Status == domain.StatusConfirmed || batch.Status == domain.StatusPublished || batch.Status == domain.StatusArchived) && domain.ConfirmAll(&batch), nil
}
func (s *Service) DeleteNote(batchID, noteID string) error {
	notes, err := s.Store.ListNotes(batchID)
	if err != nil {
		return err
	}
	for _, note := range notes {
		if note.ID == noteID {
			return s.Store.Delete([]byte("notes"), noteID)
		}
	}
	return domain.ErrNotFound
}
func (s *Service) ExportHistory(batchID string) ([]domain.ExportSnapshot, error) {
	values, err := s.Store.ListExports(batchID)
	if err != nil {
		return nil, err
	}
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values, nil
}
func (s *Service) AuditHistory(batchID string) ([]domain.AuditEvent, error) {
	return s.Store.ListEvents(batchID)
}
