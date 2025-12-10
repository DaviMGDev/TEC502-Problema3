package data

import (
	"cod-server/internal/utils"
	"errors"
)

type Repository[T any] interface {
	Create(id string, entity T) error
	Read(id string) (T, error)
	Update(id string, entity T) error
	Delete(id string) error
	List() ([]T, error)
}

type InMemoryRepository[T any] struct {
	entities utils.Map[string, T]
}

func NewInMemoryRepository[T any]() *InMemoryRepository[T] {
	return &InMemoryRepository[T]{
		entities: utils.NewSafeMap[string, T](),
	}
}

func (repo *InMemoryRepository[T]) Create(id string, entity T) error {
	if _, exists := repo.entities.Get(id); exists {
		return errors.New("entity with this ID already exists")
	}
	repo.entities.Set(id, entity)
	return nil
}

func (repo *InMemoryRepository[T]) Read(id string) (T, error) {
	entity, exists := repo.entities.Get(id)
	if !exists {
		var zero T
		return zero, errors.New("entity not found")
	}
	return entity, nil
}

func (repo *InMemoryRepository[T]) Update(id string, entity T) error {
	_, exists := repo.entities.Get(id)
	if !exists {
		return errors.New("entity not found")
	}
	repo.entities.Set(id, entity)
	return nil
}

func (repo *InMemoryRepository[T]) Delete(id string) error {
	_, exists := repo.entities.Get(id)
	if !exists {
		return errors.New("entity not found")
	}
	repo.entities.Delete(id)
	return nil
}

func (repo *InMemoryRepository[T]) List() ([]T, error) {
	return repo.entities.Values(), nil
}
