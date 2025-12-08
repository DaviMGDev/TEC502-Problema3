package utils 

import (
	"sync"
	"slices"
)

type List[T any] interface {
	Append(item T)
	Remove(item T) bool
	Contains(item T) bool
	Size() int
	Get(index int) (T, bool)
	Search(item T) int
	Pop() (T, bool)
}

type SafeList[T comparable] struct {
	items []T
	mutex sync.RWMutex
}

func NewSafeList[T comparable]() *SafeList[T] {
	return &SafeList[T]{
		items: make([]T, 0),
	}
}

func NewSafeListFromSlice[T comparable](slice []T) *SafeList[T] {
	return &SafeList[T]{
		items: slice,
	}
}

func (l *SafeList[T]) Append(item T) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items = append(l.items, item)
}
func (l *SafeList[T]) Remove(item T) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for i, v := range l.items {
		if v == item {
			l.items = append(l.items[:i], l.items[i+1:]...)
			return true
		}
	}
	return false
}
func (l *SafeList[T]) Contains(item T) bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return slices.Contains(l.items, item)
}
func (l *SafeList[T]) Size() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return len(l.items)
}
func (l *SafeList[T]) Get(index int) (T, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if index < 0 || index >= len(l.items) {
		var zero T
		return zero, false
	}
	return l.items[index], true
}
func (l *SafeList[T]) Search(item T) int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for i, v := range l.items {
		if v == item {
			return i
		}
	}
	return -1
}
func (l *SafeList[T]) Pop() (T, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if len(l.items) == 0 {
		var zero T
		return zero, false
	}
	item := l.items[len(l.items)-1]
	l.items = l.items[:len(l.items)-1]
	return item, true
}	
