package service

import (
	"fmt"
	"strconv"
	"vyshkat/domain"
)

type InventoryRepository interface {
	Add(item domain.Inventorizable) error
	AllItems() ([]domain.Inventorizable, error)
	Size() (int, error)
}

type ServiceCenter interface {
	Inspect(item domain.Inventorizable) error
}

type InventoryService struct {
	repo      InventoryRepository
	inspector ServiceCenter
}

func NewInventoryService(repo InventoryRepository, inspector ServiceCenter) *InventoryService {
	return &InventoryService{
		repo:      repo,
		inspector: inspector,
	}
}

func (s *InventoryService) AddItemToFleet(item domain.Inventorizable) error {
	if err := s.inspector.Inspect(item); err != nil {
		return fmt.Errorf("approve failed: %w", err)
	}

	size, err := s.repo.Size()
	if err != nil {
		return fmt.Errorf("repo size: %w", err)
	}
	item.SetInventoryNumber(domain.InventoryNumber{ID: strconv.Itoa(size)})
	return s.repo.Add(item)
}

func (s *InventoryService) TotalDailyEnergy() (int, error) {
	items, err := s.repo.AllItems()
	if err != nil {
		return 0, fmt.Errorf("all items: %w", err)
	}

	sum := 0
	for _, item := range items {
		if value, ok := item.(domain.EnergyConsumer); ok {
			sum += value.EnergyConsumption()
		}
	}
	return sum, nil
}

func (s *InventoryService) BeginnerVehicles() ([]domain.Vehicle, error) {
	items, err := s.repo.AllItems()
	if err != nil {
		return nil, fmt.Errorf("all items: %w", err)
	}

	var vehicles []domain.Vehicle
	for _, item := range items {
		if vehicle, ok := item.(domain.Vehicle); ok && s.isGoodForBeginner(vehicle) {
			vehicles = append(vehicles, vehicle)
		}
	}
	return vehicles, nil
}

func (s *InventoryService) AllItems() ([]domain.Inventorizable, error) {
	return s.repo.AllItems()
}

func (s *InventoryService) isGoodForBeginner(vehicle domain.Vehicle) bool {
	return vehicle.Simplicity() >= 6
}

type ItemInspector struct {
}

func (i ItemInspector) Inspect(item domain.Inventorizable) error {
	return nil
}
