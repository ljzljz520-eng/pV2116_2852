package store

import (
	"sort"
	"stickerchallenge/internal/domain"
)

var exportBucket = []byte("exports")

func (s *Store) PutExport(snapshot domain.ExportSnapshot) error {
	if snapshot.ID == "" || snapshot.BatchID == "" || snapshot.Payload == "" {
		return domain.ErrInvalid
	}
	return s.Put(exportBucket, snapshot.ID, snapshot)
}

func (s *Store) GetExport(id string) (domain.ExportSnapshot, error) {
	var snapshot domain.ExportSnapshot
	err := s.Get(exportBucket, id, &snapshot)
	return snapshot, err
}

func (s *Store) ListExports(batchID string) ([]domain.ExportSnapshot, error) {
	values, err := s.List(exportBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ExportSnapshot, 0)
	for _, value := range values {
		var snapshot domain.ExportSnapshot
		if err := decode(value, &snapshot); err != nil {
			return nil, err
		}
		if snapshot.BatchID == batchID {
			result = append(result, snapshot)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}
