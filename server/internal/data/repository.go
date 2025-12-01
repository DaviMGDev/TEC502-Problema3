package data

import (
	"errors"
	"sync"
)

type Repository[T any] interface {
	Create(id string, entity T) error 
	Read(id string) (T, error)
	Update(id string, entity T) error 
	Delete(id string) error 
	List() ([]T, error)
	ListBy(func (T) bool) ([]T, error)
}

type inMemoryRepository[T any] struct {
	entities map[string]T
	mutex sync.RWMutex
}

func NewInMemoryRepository[T any]() Repository[T] {
	return &inMemoryRepository[T]{
		entities: make(map[string]T),
	}
}

func (r *inMemoryRepository[T]) Create(id string, entity T) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.entities[id] = entity
	return nil
}

func (r *inMemoryRepository[T]) Read(id string) (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	entity, exists := r.entities[id]
	if !exists {
		var zero T
		return zero, errors.New("entity not found")
	}
	return entity, nil
}

func (r *inMemoryRepository[T]) Update(id string, entity T) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, exists := r.entities[id]
	if !exists {
		return errors.New("entity not found")
	}
	r.entities[id] = entity
	return nil
}

func (r *inMemoryRepository[T]) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, exists := r.entities[id]
	if !exists {
		return errors.New("entity not found")
	}
	delete(r.entities, id)
	return nil
}

func (r *inMemoryRepository[T]) List() ([]T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	entities := make([]T, 0, len(r.entities))
	for _, entity := range r.entities {
		entities = append(entities, entity)
	}
	return entities, nil
}

func (r *inMemoryRepository[T]) ListBy(predicate func(T) bool) ([]T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	entities := make([]T, 0)
	for _, entity := range r.entities {
		if predicate(entity) {
			entities = append(entities, entity)
		}
	}
	return entities, nil
}
