package in_memory

import (
	"distributed-algorithms/src/key-value"
	"fmt"
	"sync"
)

type Storage struct {
	storage map[string]any
	mutex   sync.RWMutex
}

func NewStorage() key_value.Storage {
	return &Storage{
		storage: make(map[string]any),
		mutex:   sync.RWMutex{},
	}
}

func (s *Storage) Get(key string) key_value.Value {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	v, ok := s.storage[key]

	value := key_value.Value{
		Value:  v,
		Exists: ok,
	}

	return value
}

func (s *Storage) Set(key string, value any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.storage[key] = value
}

func (s *Storage) CompareAndSet(key string, oldValue any, newValue any) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.storage[key]

	if !ok {
		return false, fmt.Errorf("key \"%s\" doesn't exist", key)
	}

	var success = false
	if v == oldValue {
		success = true
		s.storage[key] = newValue
	}

	return success, nil
}

func (s *Storage) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.storage, key)
}

func (s *Storage) AddElement(key string, subKey string, value any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.storage[key]

	if !ok {
		s.storage[key] = map[string]any{subKey: value}
	} else {
		s.storage[key].(map[string]any)[subKey] = value
	}
}
