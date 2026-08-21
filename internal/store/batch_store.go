package store

import (
	"sort"
	"stickerchallenge/internal/domain"
)

var batchBucket = []byte("batches")

func (s *Store) PutBatch(batch domain.Batch) error {
	if err := domain.ValidateBatch(batch); err != nil {
		return err
	}
	return s.Put(batchBucket, batch.ID, batch)
}

func (s *Store) GetBatch(id string) (domain.Batch, error) {
	var batch domain.Batch
	err := s.Get(batchBucket, id, &batch)
	return batch, err
}

func (s *Store) UpdateBatch(batch domain.Batch, expectedVersion int) error {
	current, err := s.GetBatch(batch.ID)
	if err != nil {
		return err
	}
	if current.Version != expectedVersion {
		return domain.ErrConflict
	}
	batch.Version = expectedVersion + 1
	if err := domain.ValidateBatch(batch); err != nil {
		return err
	}
	return s.PutBatch(batch)
}

func (s *Store) ListBatches() ([]domain.Batch, error) {
	values, err := s.List(batchBucket)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Batch, 0, len(values))
	for _, value := range values {
		var batch domain.Batch
		if err := decode(value, &batch); err != nil {
			return nil, err
		}
		result = append(result, batch)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func decode(data []byte, target any) error {
	return jsonUnmarshal(data, target)
}
