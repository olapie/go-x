package xsync

import "sync"

type Set[K comparable] struct {
	mu sync.RWMutex
	m  map[K]struct{}
}

func (s *Set[K]) Add(key K) {
	s.mu.Lock()
	s.m[key] = struct{}{}
	s.mu.Unlock()
}

func (s *Set[K]) Contains(key K) bool {
	s.mu.RLock()
	_, ok := s.m[key]
	s.mu.RUnlock()
	return ok
}

func (s *Set[K]) Delete(key K) {
	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()
}
