package ram

import (
	"distributed-algorithms/src/key-value"
	"sync"
)

type Storage struct {
	storage map[string]any
	mutex   sync.RWMutex
}

func NewStorage() key_value.KeyValueStorage {
	return &Storage{storage: make(map[string]any)}
}

func (s *Storage) Get(key string) (key_value.Value, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	v, ok := s.storage[key]

	value := key_value.Value{
		Value:  v,
		Exists: ok,
	}

	return value, nil
}

func (s *Storage) Set(key string, value any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.storage[key] = value
	return nil
}

func (s *Storage) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.storage, key)
	return nil
}
