package storage

import "sync"

type Repository[T any] interface {
	Create(id string, entity T) error 
	Read(id string) (T, error)
	Update(id string, entity T) error 
	Delete(id string) error 
	List() ([]T, error) 
	ListBy(func(T) bool) ([]T, error)
}

type InMemoryRepository[T any] struct {
	data map[string]T
	mutex sync.RWMutex 
}

func NewInMemoryRepository[T any]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{
		data: make(map[string]T),
		mutex: sync.RWMutex{},
	}
}

func (r *InMemoryRepository[T]) Create(id string, entity T) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[id] = entity
	return nil
}

func (r *InMemoryRepository[T]) Read(id string) (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	entity, exists := r.data[id]
	if !exists {
		var zero T
		return zero, nil 
	}
	return entity, nil
}

func (r *InMemoryRepository[T]) Update(id string, entity T) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.data[id] = entity
	return nil
}

func (r *InMemoryRepository[T]) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.data, id)
	return nil 
}

func (r *InMemoryRepository[T]) List() ([]T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	entities := make([]T, 0, len(r.data)) 
	for _, entity := range r.data {
		entities = append(entities, entity)
	}
	return entities, nil
}

func (r *InMemoryRepository[T]) ListBy(predicate func(T) bool) ([]T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	entities := make([]T, 0) 
	for _, entity := range r.data {
		if predicate(entity) {
			entities = append(entities, entity)
		}
	}
	return entities, nil
}

