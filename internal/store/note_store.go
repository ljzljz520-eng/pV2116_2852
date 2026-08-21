package store

import (
	"sort"
	"stickerchallenge/internal/domain"
)

var noteBucket = []byte("notes")

func (s *Store) PutNote(note domain.CollaborationNote) error {
	if note.ID == "" || note.BatchID == "" || note.Author == "" || note.Body == "" {
		return domain.ErrInvalid
	}
	return s.Put(noteBucket, note.ID, note)
}

func (s *Store) ListNotes(batchID string) ([]domain.CollaborationNote, error) {
	values, err := s.List(noteBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.CollaborationNote, 0)
	for _, value := range values {
		var note domain.CollaborationNote
		if err := decode(value, &note); err != nil {
			return nil, err
		}
		if note.BatchID == batchID {
			result = append(result, note)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}
