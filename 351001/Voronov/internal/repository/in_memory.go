package repository

import (
	"Voronov/internal/errors"
	"sync"
	"sync/atomic"
)

type InMemoryRepository[T any] struct {
	items    map[int64]*T
	nextID   int64
	mu       sync.RWMutex
	idGetter func(*T) int64
	idSetter func(*T, int64)
}

func NewInMemoryRepository[T any](idGetter func(*T) int64, idSetter func(*T, int64)) *InMemoryRepository[T] {
	return &InMemoryRepository[T]{
		items:    make(map[int64]*T),
		nextID:   1,
		idGetter: idGetter,
		idSetter: idSetter,
	}
}

func (r *InMemoryRepository[T]) FindByID(id int64) (*T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, errors.ErrNotFound
	}
	return item, nil
}

func (r *InMemoryRepository[T]) FindAll() ([]*T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*T, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}
	return result, nil
}

func (r *InMemoryRepository[T]) Create(entity *T) (*T, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := atomic.AddInt64(&r.nextID, 1) - 1
	r.idSetter(entity, id)
	r.items[id] = entity
	return entity, nil
}

func (r *InMemoryRepository[T]) Update(id int64, entity *T) (*T, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return nil, errors.ErrNotFound
	}

	r.idSetter(entity, id)
	r.items[id] = entity
	return entity, nil
}

func (r *InMemoryRepository[T]) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return errors.ErrNotFound
	}

	delete(r.items, id)
	return nil
}
