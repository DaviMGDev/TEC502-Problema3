package utils

import "sync"

type Map[K comparable, V Comparable] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
	Delete(key K)
	Has(key K) bool
	Keys() []K
	Values() []V
	Size() int
	Clear()
	Equals(other Comparable) bool
}

type SafeMap[K comparable, V Comparable] struct {
	items map[K]V
	mutex sync.RWMutex
}

func NewSafeMap[K comparable, V Comparable]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		items: make(map[K]V),
	}
}

func (m *SafeMap[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, exists := m.items[key]
	return value, exists
}

func (m *SafeMap[K, V]) Set(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.items[key] = value
}

func (m *SafeMap[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.items, key)
}

func (m *SafeMap[K, V]) Has(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, exists := m.items[key]
	return exists
}

func (m *SafeMap[K, V]) Keys() []K {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	keys := make([]K, 0, len(m.items))
	for key := range m.items {
		keys = append(keys, key)
	}
	return keys
}

func (m *SafeMap[K, V]) Values() []V {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	values := make([]V, 0, len(m.items))
	for _, value := range m.items {
		values = append(values, value)
	}
	return values
}

func (m *SafeMap[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.items)
}

func (m *SafeMap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.items = make(map[K]V)
}

func (m *SafeMap[K, V]) Equals(other Comparable) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	otherMap, ok := other.(*SafeMap[K, V])
	if !ok {
		return false
	}
	if m.Size() != otherMap.Size() {
		return false
	}
	for key, value := range m.items {
		otherValue, exists := otherMap.Get(key)
		if !exists || !value.Equals(otherValue) {
			return false
		}
	}
	return true
}


