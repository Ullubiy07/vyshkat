package repository

import (
	"sync"
	"vyshkat/domain"
)

type MemoryInventoryRepository struct {
	items []domain.Inventorizable
	mutex sync.Mutex
}

func NewMemoryInventoryRepository() *MemoryInventoryRepository {
	return &MemoryInventoryRepository{}
}

func (r *MemoryInventoryRepository) Add(item domain.Inventorizable) error {
	r.mutex.Lock()
	r.items = append(r.items, item)
	r.mutex.Unlock()
	return nil
}

func (r *MemoryInventoryRepository) AllItems() ([]domain.Inventorizable, error) {
	r.mutex.Lock()
	res := make([]domain.Inventorizable, len(r.items))
	copy(res, r.items)
	r.mutex.Unlock()
	return res, nil
}

func (r *MemoryInventoryRepository) Size() (int, error) {
	r.mutex.Lock()
	size := len(r.items)
	r.mutex.Unlock()
	return size, nil
}
