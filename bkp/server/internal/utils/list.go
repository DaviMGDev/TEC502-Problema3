package utils

import (
	"errors"
	"sync"
)

type List[T Comparable] interface {
	Append(item T)
	Extract(intdex int) (T, error)
	Get(index int) (T, error)
	Size() int
	Clear()
	ToSlice() []T
	Map(func(T) T) List[T]
	Filter(func(T) bool) List[T]
	ForEach(func(T))
	Search(T) int
	Equals(other Comparable) bool
}

type SafeList[T Comparable] struct {
	items []T
	mutex sync.RWMutex
}

func NewSafeList[T Comparable]() *SafeList[T] {
	return &SafeList[T]{
		items: make([]T, 0),
	}
}	

func (list *SafeList[T]) Append(item T) { 	
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.items = append(list.items, item)
}

func (list *SafeList[T]) Extract(index int) (T, error) {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	if index < 0 || index >= len(list.items) {
		var zero T
		return zero, errors.New("index out of range")
	}
	item := list.items[index]
	list.items = append(list.items[:index], list.items[index+1:]...)
	return item, nil
}

func (list *SafeList[T]) Get(index int) (T, error) {
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	if index < 0 || index >= len(list.items) {
		var zero T
		return zero, errors.New("index out of range")
	}
	return list.items[index], nil
}

func (list *SafeList[T]) Size() int { 
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	return len(list.items)
}

func (list *SafeList[T]) Clear() {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	list.items = make([]T, 0)
}

func (list *SafeList[T]) ToSlice() []T {
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	copied := make([]T, len(list.items))
	copy(copied, list.items)
	return copied
}

func (list *SafeList[T]) Map(f func(T) T) List[T] {
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	newList := NewSafeList[T]()
	for _, item := range list.items {
		newList.Append(f(item))
	}
	return newList
}

func (list *SafeList[T]) Filter(f func(T) bool) List[T] {
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	newList := NewSafeList[T]()
	for _, item := range list.items {
		if f(item) {
			newList.Append(item)
		}
	}
	return newList
}

func (list *SafeList[T]) ForEach(f func(T)) {
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	for _, item := range list.items {
		f(item)
	}
}

func (list *SafeList[T]) Search(target T) int {
	list.mutex.RLock()
	defer list.mutex.RUnlock()
	for index, item := range list.items {
		if item.Equals(target) {
			return index
		}
	}
	return -1
}	 

func (list *SafeList[T]) Equals(other Comparable) bool {
	otherList, ok := other.(*SafeList[T])
	if !ok {
		return false
	}
	list.mutex.RLock()	
	otherList.mutex.RLock()
	defer list.mutex.RUnlock()
	defer otherList.mutex.RUnlock()

	if len(list.items) != len(otherList.items) {
		return false
	}
	for i, item := range list.items {
		if !item.Equals(otherList.items[i]) {
			return false
		}
	}
	return true
}
