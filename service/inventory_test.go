package service

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"vyshkat/domain"
)

type mockRepository struct {
	items       []domain.Inventorizable
	addErr      error
	sizeErr     error
	allItemsErr error
}

func (r *mockRepository) Add(item domain.Inventorizable) error {
	if r.addErr != nil {
		return r.addErr
	}
	r.items = append(r.items, item)
	return nil
}

func (r *mockRepository) AllItems() ([]domain.Inventorizable, error) {
	if r.allItemsErr != nil {
		return nil, r.allItemsErr
	}
	result := make([]domain.Inventorizable, len(r.items))
	copy(result, r.items)
	return result, nil
}

func (r *mockRepository) Size() (int, error) {
	if r.sizeErr != nil {
		return 0, r.sizeErr
	}
	return len(r.items), nil
}

type mockInspector struct {
	err      error
	calls    int
	lastItem domain.Inventorizable
}

func (i *mockInspector) Inspect(item domain.Inventorizable) error {
	i.calls++
	i.lastItem = item
	return i.err
}

func makeScooter(name, id string, simplicity int) domain.Scooter {
	return domain.Scooter{
		ItemInfo: &domain.ItemInfo{
			ItemName: name,
			ItemNum:  domain.NewInventoryNumber(id),
		},
		Simple: simplicity,
	}
}

func makeElectricBike(name, id string, energy, simplicity int) domain.ElectricBike {
	return domain.ElectricBike{
		ItemInfo: &domain.ItemInfo{
			ItemName: name,
			ItemNum:  domain.NewInventoryNumber(id),
		},
		Energy: energy,
		Simple: simplicity,
	}
}

func TestAddItemToFleet(t *testing.T) {
	repo := &mockRepository{
		items: []domain.Inventorizable{
			makeScooter("Existing 1", "0", 7),
			makeScooter("Existing 2", "1", 8),
		},
	}
	inspector := &mockInspector{}
	svc := NewInventoryService(repo, inspector)

	item := makeScooter("New scooter", "old-id", 9)

	err := svc.AddItemToFleet(item)
	if err != nil {
		t.Fatalf("AddItemToFleet() returned an error: %v", err)
	}

	if len(repo.items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(repo.items))
	}

	if got := item.InventoryNumber().String(); got != "2" {
		t.Errorf("expected inventory number %q, got %q", "2", got)
	}

	if inspector.calls != 1 {
		t.Errorf("expected inspector to be called once, got %d", inspector.calls)
	}

	if inspector.lastItem.Name() != "New scooter" {
		t.Errorf("expected inspector to receive the new scooter, got %q", inspector.lastItem.Name())
	}
}

func TestAddItemToFleetInspectionError(t *testing.T) {
	repo := &mockRepository{}
	inspectionErr := errors.New("inspection failed")
	inspector := &mockInspector{err: inspectionErr}
	svc := NewInventoryService(repo, inspector)

	item := makeScooter("Scooter", "old-id", 7)

	err := svc.AddItemToFleet(item)
	if !errors.Is(err, inspectionErr) {
		t.Fatalf("expected inspection error to be wrapped, got %v", err)
	}

	if len(repo.items) != 0 {
		t.Errorf("expected no items to be added after failed inspection, got %d", len(repo.items))
	}
}

func TestAddItemToFleetSizeError(t *testing.T) {
	sizeErr := errors.New("size failed")
	repo := &mockRepository{sizeErr: sizeErr}
	inspector := &mockInspector{}
	svc := NewInventoryService(repo, inspector)

	err := svc.AddItemToFleet(makeScooter("Scooter", "old-id", 7))
	if !errors.Is(err, sizeErr) {
		t.Fatalf("expected size error to be wrapped, got %v", err)
	}

	if len(repo.items) != 0 {
		t.Errorf("expected no items to be added, got %d", len(repo.items))
	}
}

func TestAddItemToFleetRepositoryAddError(t *testing.T) {
	addErr := errors.New("add failed")
	repo := &mockRepository{addErr: addErr}
	inspector := &mockInspector{}
	svc := NewInventoryService(repo, inspector)

	err := svc.AddItemToFleet(makeScooter("Scooter", "old-id", 7))
	if !errors.Is(err, addErr) {
		t.Fatalf("expected add error to be returned, got %v", err)
	}
}

func TestTotalDailyEnergy(t *testing.T) {
	repo := &mockRepository{
		items: []domain.Inventorizable{
			makeElectricBike("Electric bike", "1", 500, 7),
			domain.ChargingStation{
				ItemInfo: &domain.ItemInfo{
					ItemName: "Charging station",
					ItemNum:  domain.NewInventoryNumber("2"),
				},
				Energy: 1000,
			},
			makeScooter("Regular scooter", "3", 8),
		},
	}
	svc := NewInventoryService(repo, &mockInspector{})

	got, err := svc.TotalDailyEnergy()
	if err != nil {
		t.Fatalf("TotalDailyEnergy() returned an error: %v", err)
	}

	const expected = 1500
	if got != expected {
		t.Errorf("expected total energy consumption %d, got %d", expected, got)
	}
}

func TestTotalDailyEnergyRepositoryError(t *testing.T) {
	repoErr := errors.New("read failed")
	repo := &mockRepository{allItemsErr: repoErr}
	svc := NewInventoryService(repo, &mockInspector{})

	got, err := svc.TotalDailyEnergy()
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error to be wrapped, got %v", err)
	}
	if got != 0 {
		t.Errorf("expected energy consumption 0 on error, got %d", got)
	}
}

func TestBeginnerVehicles(t *testing.T) {
	easyScooter := makeScooter("Easy scooter", "1", 6)
	hardScooter := makeScooter("Hard scooter", "2", 5)
	easyBike := makeElectricBike("Easy electric bike", "3", 300, 10)

	repo := &mockRepository{
		items: []domain.Inventorizable{
			easyScooter,
			hardScooter,
			easyBike,
			domain.Helmet{
				ItemInfo: &domain.ItemInfo{
					ItemName: "Helmet",
					ItemNum:  domain.NewInventoryNumber("4"),
				},
			},
		},
	}
	svc := NewInventoryService(repo, &mockInspector{})

	got, err := svc.BeginnerVehicles()
	if err != nil {
		t.Fatalf("BeginnerVehicles() returned an error: %v", err)
	}

	expected := []domain.Vehicle{easyScooter, easyBike}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("unexpected beginner vehicles:\n got: %#v\nwant: %#v", got, expected)
	}
}

func TestBeginnerVehiclesRepositoryError(t *testing.T) {
	repoErr := errors.New("read failed")
	repo := &mockRepository{allItemsErr: repoErr}
	svc := NewInventoryService(repo, &mockInspector{})

	got, err := svc.BeginnerVehicles()
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error to be wrapped, got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil vehicles on error, got %#v", got)
	}
}

func TestAllItems(t *testing.T) {
	items := []domain.Inventorizable{
		makeScooter("Scooter", "1", 7),
		makeElectricBike("Electric bike", "2", 400, 8),
	}
	repo := &mockRepository{items: items}
	svc := NewInventoryService(repo, &mockInspector{})

	got, err := svc.AllItems()
	if err != nil {
		t.Fatalf("AllItems() returned an error: %v", err)
	}

	if !reflect.DeepEqual(got, items) {
		t.Errorf("unexpected items:\n got: %#v\nwant: %#v", got, items)
	}
}

func TestAllItemsRepositoryError(t *testing.T) {
	repoErr := errors.New("read failed")
	repo := &mockRepository{allItemsErr: repoErr}
	svc := NewInventoryService(repo, &mockInspector{})

	_, err := svc.AllItems()
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestAddItemToFleetErrorMessages(t *testing.T) {
	tests := []struct {
		name      string
		repo      *mockRepository
		inspector *mockInspector
		want      string
	}{
		{
			name:      "inspection error",
			repo:      &mockRepository{},
			inspector: &mockInspector{err: errors.New("inspection failed")},
			want:      "approve failed",
		},
		{
			name:      "size error",
			repo:      &mockRepository{sizeErr: errors.New("size failed")},
			inspector: &mockInspector{},
			want:      "repo size",
		},
		{
			name:      "add error",
			repo:      &mockRepository{addErr: errors.New("add failed")},
			inspector: &mockInspector{},
			want:      "add failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewInventoryService(tt.repo, tt.inspector)
			err := svc.AddItemToFleet(makeScooter("Scooter", "1", 7))

			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("expected error containing %q, got %q", tt.want, fmt.Sprint(err))
			}
		})
	}
}

func TestStupidInspectorInspect(t *testing.T) {
	inspector := ItemInspector{}
	item := &domain.Scooter{
		ItemInfo: &domain.ItemInfo{ItemName: "Scooter"},
	}

	err := inspector.Inspect(item)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
