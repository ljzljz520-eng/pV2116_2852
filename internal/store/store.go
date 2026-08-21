package store

import (
	"encoding/json"
	"os"

	"go.etcd.io/bbolt"
	"stickerchallenge/internal/domain"
)

var buckets = [][]byte{[]byte("batches"), []byte("events"), []byte("notes"), []byte("exports")}

type Store struct {
	db *bbolt.DB
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0o600, &bbolt.Options{NoSync: true})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range buckets {
			if _, createErr := tx.CreateBucketIfNotExists(bucket); createErr != nil {
				return createErr
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func OpenMemory(path string) (*Store, error) {
	if path == "" {
		path = "stickerchallenge.db"
	}
	return Open(path)
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Put(bucket []byte, key string, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).Put([]byte(key), encoded)
	})
}

func (s *Store) Get(bucket []byte, key string, target any) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		value := tx.Bucket(bucket).Get([]byte(key))
		if value == nil {
			return domain.ErrNotFound
		}
		return json.Unmarshal(append([]byte(nil), value...), target)
	})
}

func (s *Store) Delete(bucket []byte, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Delete([]byte(key)) })
}

func (s *Store) List(bucket []byte) ([][]byte, error) {
	items := make([][]byte, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(k, v []byte) error {
			if v != nil {
				items = append(items, append([]byte(nil), v...))
			}
			return nil
		})
	})
	return items, err
}

func Remove(path string) error {
	if path == "" {
		return nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}
