package repository

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"vyshkat/domain"
)

func testItems() []domain.Inventorizable {
	return []domain.Inventorizable{
		domain.ElectricBike{
			ItemInfo: &domain.ItemInfo{
				ItemName: "Electric Bike",
				ItemNum:  domain.NewInventoryNumber("28"),
			},
			Energy: 500,
			Simple: 4,
		},
		domain.ElectricScooter{
			ItemInfo: &domain.ItemInfo{
				ItemName: "Electric Scooter",
				ItemNum:  domain.NewInventoryNumber("29"),
			},
			Energy: 50,
			Simple: 5,
		},
		domain.Bike{
			ItemInfo: &domain.ItemInfo{
				ItemName: "Bike",
				ItemNum:  domain.NewInventoryNumber("30"),
			},
			Simple: 8,
		},
		domain.Scooter{
			ItemInfo: &domain.ItemInfo{
				ItemName: "Scooter",
				ItemNum:  domain.NewInventoryNumber("31"),
			},
			Simple: 10,
		},
		domain.Helmet{
			ItemInfo: &domain.ItemInfo{
				ItemName: "Helmet",
				ItemNum:  domain.NewInventoryNumber("32"),
			},
		},
		domain.ChargingStation{
			ItemInfo: &domain.ItemInfo{
				ItemName: "Charging Station",
				ItemNum:  domain.NewInventoryNumber("33"),
			},
		},
	}
}

func TestEmptyRepository(t *testing.T) {
	repo := NewMemoryInventoryRepository()

	size, err := repo.Size()
	if err != nil {
		t.Fatalf("Size() returned an error: %v", err)
	}
	if size != 0 {
		t.Errorf("expected empty repository size to be 0, got %d", size)
	}

	items, err := repo.AllItems()
	if err != nil {
		t.Fatalf("AllItems() returned an error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestAddItems(t *testing.T) {
	repo := NewMemoryInventoryRepository()
	expected := testItems()

	for _, item := range expected {
		if err := repo.Add(item); err != nil {
			t.Fatalf("Add() returned an error: %v", err)
		}
	}

	size, err := repo.Size()
	if err != nil {
		t.Fatalf("Size() returned an error: %v", err)
	}
	if size != len(expected) {
		t.Errorf("expected size %d, got %d", len(expected), size)
	}

	actual, err := repo.AllItems()
	if err != nil {
		t.Fatalf("AllItems() returned an error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("AllItems() returned unexpected items:\n got: %#v\nwant: %#v", actual, expected)
	}
}

func TestAllItemsReturnsIndependentSlice(t *testing.T) {
	repo := NewMemoryInventoryRepository()
	expected := testItems()

	if err := repo.Add(expected[0]); err != nil {
		t.Fatalf("Add() returned an error: %v", err)
	}

	items, err := repo.AllItems()
	if err != nil {
		t.Fatalf("AllItems() returned an error: %v", err)
	}

	items[0] = expected[1]

	actual, err := repo.AllItems()
	if err != nil {
		t.Fatalf("AllItems() returned an error: %v", err)
	}

	if !reflect.DeepEqual(actual[0], expected[0]) {
		t.Error("modifying the returned slice changed the repository")
	}
}

func TestParallelInteraction(t *testing.T) {
	repo := NewMemoryInventoryRepository()

	const itemCount = 100
	const readerCount = 10

	var wg sync.WaitGroup

	for i := 0; i < itemCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			item := domain.Scooter{
				ItemInfo: &domain.ItemInfo{
					ItemName: fmt.Sprintf("Scooter %d", i),
					ItemNum:  domain.NewInventoryNumber(fmt.Sprintf("%d", i)),
				},
				Simple: 7,
			}

			if err := repo.Add(item); err != nil {
				t.Errorf("Add() returned an error: %v", err)
			}
		}(i)
	}

	for i := 0; i < readerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if _, err := repo.Size(); err != nil {
				t.Errorf("Size() returned an error: %v", err)
			}
			if _, err := repo.AllItems(); err != nil {
				t.Errorf("AllItems() returned an error: %v", err)
			}
		}()
	}

	wg.Wait()

	size, err := repo.Size()
	if err != nil {
		t.Fatalf("Size() returned an error: %v", err)
	}
	if size != itemCount {
		t.Errorf("expected size %d after concurrent additions, got %d", itemCount, size)
	}

	items, err := repo.AllItems()
	if err != nil {
		t.Fatalf("AllItems() returned an error: %v", err)
	}
	if len(items) != itemCount {
		t.Errorf("expected %d items, got %d", itemCount, len(items))
	}
}
