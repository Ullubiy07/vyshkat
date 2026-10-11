package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"vyshkat/domain"
	"vyshkat/repository"
	"vyshkat/service"
)

type failingRepository struct {
	service.InventoryRepository
}

func (r failingRepository) AllItems() ([]domain.Inventorizable, error) {
	return nil, fmt.Errorf("repository error")
}

func (r failingRepository) Add(domain.Inventorizable) error {
	return fmt.Errorf("repository error")
}

func (r failingRepository) Size() (int, error) {
	return 0, fmt.Errorf("repository error")
}

type failingInspector struct{}

func (failingInspector) Inspect(domain.Inventorizable) error {
	return fmt.Errorf("inspection error")
}

func TestHandlersServiceErrors(t *testing.T) {
	inventory := service.NewInventoryService(
		failingRepository{},
		service.ItemInspector{},
	)

	var output bytes.Buffer
	handler := NewCLIHandler(
		inventory,
		bufio.NewScanner(strings.NewReader("")),
		&output,
	)

	handler.handleAllItems()
	handler.handleBeginnerVehicles()
	handler.handleTotalDailyEnergy()

	if strings.Count(output.String(), "Ошибка:") != 3 {
		t.Fatalf("expected 3 errors, got:\n%s", output.String())
	}
}

func TestAddItemServiceError(t *testing.T) {
	inventory := service.NewInventoryService(
		repository.NewMemoryInventoryRepository(),
		failingInspector{},
	)

	output := runCLIWithInventory("1\n1\nСамокат\n8\n6\n", inventory)

	if !strings.Contains(output, "Ошибка:") {
		t.Fatalf("expected inspection error, got:\n%s", output)
	}
}

func TestElectricBikeEnergyInput(t *testing.T) {
	output := runCLI("1\n4\nЭлектровелосипед\n8\n15\n5\n6\n")

	if !strings.Contains(output, "Электровелосипед") {
		t.Fatalf("electric bike was not added, output:\n%s", output)
	}
}

func TestInputErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"menu EOF", ""},
		{"vehicle type EOF", "1\n"},
		{"vehicle name EOF", "1\n1\n"},
		{"simplicity EOF", "1\n1\nСамокат\n"},
		{"electric vehicle energy EOF", "1\n3\nСамокат\n8\n"},
		{"item type EOF", "2\n"},
		{"item name EOF", "2\n1\n"},
		{"charging station energy EOF", "2\n2\nСтанция\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := runCLI(tt.input)
			if !strings.Contains(output, "Ошибка чтения:") {
				t.Fatalf("expected input error, got:\n%s", output)
			}
		})
	}
}

func TestInvalidTypeChoices(t *testing.T) {
	tests := []string{
		"1\nabc\n",
		"1\n5\n",
		"2\nabc\n",
		"2\n3\n",
	}

	for _, input := range tests {
		output := runCLI(input)
		if !strings.Contains(output, "введите целое число") {
			t.Errorf("expected invalid choice error for %q, got:\n%s", input, output)
		}
	}
}

func runCLIWithInventory(input string, inventory *service.InventoryService) string {
	var output bytes.Buffer
	handler := NewCLIHandler(
		inventory,
		bufio.NewScanner(strings.NewReader(input)),
		&output,
	)
	handler.HandleRequests()
	return output.String()
}

func runCLI(input string) string {
	repo := repository.NewMemoryInventoryRepository()
	inventory := service.NewInventoryService(repo, service.ItemInspector{})

	var output bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(input))

	handler := NewCLIHandler(inventory, scanner, &output)
	handler.HandleRequests()

	return output.String()
}

func TestHandleRequestsExit(t *testing.T) {
	output := runCLI("6\n")

	if !strings.Contains(output, "Меню управления") {
		t.Fatal("menu was not displayed")
	}
}

func TestAddVehicleAndListItems(t *testing.T) {
	output := runCLI("1\n1\nОбычный самокат\n8\n5\n6\n")

	if !strings.Contains(output, "Обычный самокат") {
		t.Fatal("vehicle was not added or displayed")
	}
	if !strings.Contains(output, "0") {
		t.Fatal("inventory number was not displayed")
	}
}

func TestAddElectricVehicleAndCalculateEnergy(t *testing.T) {
	output := runCLI(
		"1\n3\nЭлектросамокат\n7\n12\n" +
			"2\n2\nЗарядная станция\n30\n" +
			"3\n6\n",
	)

	if !strings.Contains(output, "Потребление энергии: 42") {
		t.Fatalf("expected total energy 42, output:\n%s", output)
	}
}

func TestAddItemAndListItems(t *testing.T) {
	output := runCLI("2\n1\nШлем\n5\n6\n")

	if !strings.Contains(output, "Шлем") {
		t.Fatal("item was not added or displayed")
	}
}

func TestBeginnerVehicles(t *testing.T) {
	output := runCLI(
		"1\n1\nПростой самокат\n8\n" +
			"1\n2\nСложный велосипед\n4\n" +
			"4\n6\n",
	)

	if !strings.Contains(output, "Простой самокат") {
		t.Fatal("beginner-friendly vehicle was not displayed")
	}
	if strings.Contains(output, "Сложный велосипед") {
		t.Fatal("vehicle with simplicity below 6 was displayed")
	}
}

func TestInvalidMenuChoice(t *testing.T) {
	output := runCLI("abc\n")

	if !strings.Contains(output, "введите целое число от 1 до 6") {
		t.Fatal("expected invalid choice error")
	}
}

func TestEmptyItemName(t *testing.T) {
	output := runCLI("1\n1\n   \n")

	if !strings.Contains(output, "название не может быть пустым") {
		t.Fatal("expected empty name error")
	}
}

func TestInvalidSimplicity(t *testing.T) {
	output := runCLI("1\n1\nСамокат\n11\n")

	if !strings.Contains(output, "введите целое число от 1 до 10") {
		t.Fatal("expected simplicity validation error")
	}
}

func TestInvalidEnergyConsumption(t *testing.T) {
	output := runCLI("1\n3\nЭлектросамокат\n7\n-1\n")

	if !strings.Contains(output, "неотрицательное целое число") {
		t.Fatal("expected energy validation error")
	}
}

func TestUnexpectedEndOfInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"menu choice", ""},
		{"item name", "1\n1\n"},
		{"simplicity", "1\n1\nСамокат\n"},
		{"energy", "1\n3\nСамокат\n7\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := runCLI(tt.input)

			if !strings.Contains(output, "Ошибка чтения:") {
				t.Fatal("expected input reading error")
			}
		})
	}
}

func TestAddVehicleInvalidEnergyCLI(t *testing.T) {
	runCLI("1\n4\nЭлектровелосипед\n8\nnot-a-number\n6\n")
}

func TestAddItemToFleetServiceError(t *testing.T) {
	inventory := service.NewInventoryService(
		repository.NewMemoryInventoryRepository(),
		failingInspector{},
	)

	output := runCLIWithInventory("2\n1\nШлем\n6\n", inventory)

	if !strings.Contains(output, "Ошибка:") {
		t.Fatalf("expected fleet addition error for item, got:\n%s", output)
	}
}

var _ domain.Inventorizable = (*domain.ItemInfo)(nil)
