package storage

import (
	"fmt"
	"sync"

	"github.com/abdelmounim-dev/redis/pkg/parser"
)

type Store interface {
	Set(key string, value *parser.Token) error
	Get(key string) (*parser.Token, error)
	Delete(key string) error
}

type KeyValueStore struct {
	data map[string]*parser.Token
	mu   sync.RWMutex
}

func (s *KeyValueStore) Set(key string, value *parser.Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *KeyValueStore) Get(key string) (*parser.Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("STORE GET: value with key %v not found", key)
	}
	return t, nil
}

func (s *KeyValueStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)

	return nil
}
