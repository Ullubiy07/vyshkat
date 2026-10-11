package repository

import "vyshkat/domain"

type MemoryInventoryRepository struct {
	items []domain.Inventorizable
}

func NewMemoryInventoryRepository() *MemoryInventoryRepository {
	return &MemoryInventoryRepository{}
}

func (r *MemoryInventoryRepository) Add(item domain.Inventorizable) error {
	r.items = append(r.items, item)
	return nil
}

func (r *MemoryInventoryRepository) AllItems() ([]domain.Inventorizable, error) {
	res := make([]domain.Inventorizable, len(r.items))
	copy(res, r.items)
	return res, nil
}

func (r *MemoryInventoryRepository) Size() (int, error) {
	return len(r.items), nil
}
