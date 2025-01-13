package xsync

import "sync"

type Map[K comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}

func (m *Map[K, V]) CompareAndSwap(key K, oldVal, newVal V) bool {
	return m.m.CompareAndSwap(key, oldVal, newVal)
}

func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *Map[K, V]) Load(key K) (value V, loaded bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return value, false
	}
	value, ok = v.(V)
	return value, ok
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, ok := m.m.LoadAndDelete(key)
	if !ok {
		return value, false
	}
	value, ok = v.(V)
	return value, ok
}

func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	p, ok := m.m.Swap(key, value)
	if !ok {
		return previous, false
	}
	previous, ok = p.(V)
	return previous, ok
}
