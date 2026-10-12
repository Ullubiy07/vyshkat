package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"vyshkat/domain"
	"vyshkat/service"
)

type CLIHandler struct {
	inventory *service.InventoryService
	scanner   *bufio.Scanner
	writer    io.Writer
}

func NewCLIHandler(inventory *service.InventoryService, scanner *bufio.Scanner, writer io.Writer) CLIHandler {
	return CLIHandler{
		inventory: inventory,
		scanner:   scanner,
		writer:    writer,
	}
}

func (h *CLIHandler) HandleRequests() {
	for {
		title := "============Меню управления============"
		choices := []string{"Добавить транспорт в парк", "Добавить вещь в парк", "Вывести общую дневную энергию парка",
			"Вывести транспорт для новичков", "Вывести все предметы", "Выход"}

		choice, err := h.readChoice(title, choices)
		if err != nil {
			return
		}
		
		fmt.Fprintln(h.writer)
		switch choice {
		case 1:
			h.handleAddVehicleToFleet()
		case 2:
			h.handleAddItemToFleet()
		case 3:
			h.handleTotalDailyEnergy()
		case 4:
			h.handleBeginnerVehicles()
		case 5:
			h.handleAllItems()
		case 6:
			return
		}
		fmt.Fprintln(h.writer)
	}
}

func (h *CLIHandler) handleAddVehicleToFleet() {
	title := "Выберите тип предмета:"
	choices := []string{"Самокат", "Велосипед", "Электросамокат", "Электровелосипед"}

	choice, err := h.readChoice(title, choices)
	if err != nil {
		return
	}
	name, err := h.readItemName()
	if err != nil {
		return
	}
	simplicity, err := h.readSimplicity()
	if err != nil {
		return
	}

	info := domain.ItemInfo{ItemName: name}
	var item domain.Inventorizable

	switch choice {
	case 1:
		item = domain.Scooter{ItemInfo: &info, Simple: simplicity}
	case 2:
		item = domain.Bike{ItemInfo: &info, Simple: simplicity}
	case 3:
		energy, err := h.readEnergyConsumption()
		if err != nil {
			return
		}
		scooter := domain.ElectricScooter{ItemInfo: &info, Simple: simplicity}
		scooter.Energy = energy
		item = scooter

	case 4:
		energy, err := h.readEnergyConsumption()
		if err != nil {
			return
		}
		bike := domain.ElectricBike{ItemInfo: &info, Simple: simplicity}
		bike.Energy = energy
		item = bike
	}

	if err := h.inventory.AddItemToFleet(item); err != nil {
		fmt.Fprintln(h.writer, "Ошибка:", err)
	}
}

func (h *CLIHandler) handleAddItemToFleet() {
	title := "Выберите тип вещи:"
	choices := []string{"Шлем", "Зарядная станция"}

	choice, err := h.readChoice(title, choices)
	if err != nil {
		return
	}
	name, err := h.readItemName()
	if err != nil {
		return
	}

	info := domain.ItemInfo{ItemName: name}
	var item domain.Inventorizable

	switch choice {
	case 1:
		item = domain.Helmet{ItemInfo: &info}
	case 2:
		energy, err := h.readEnergyConsumption()
		if err != nil {
			return
		}
		station := domain.ChargingStation{ItemInfo: &info}
		station.Energy = energy
		item = station
	}

	if err := h.inventory.AddItemToFleet(item); err != nil {
		fmt.Fprintln(h.writer, "Ошибка:", err)
	}
}

func (h *CLIHandler) readChoice(title string, choices []string) (int, error) {
	fmt.Fprintln(h.writer, title)
	for i, choice := range choices {
		fmt.Fprintf(h.writer, "%d. %s\n", i+1, choice)
	}

	for {
		if err := h.readLine(); err != nil {
			return 0, fmt.Errorf("read line: %w", err)
		}

		choice, err := strconv.Atoi(strings.TrimSpace(h.scanner.Text()))
		if err != nil || choice < 1 || choice > len(choices) {
			fmt.Fprintf(h.writer, "Ошибка: введите целое число от 1 до %d\n", len(choices))
		} else {
			return choice, nil
		}
	}
}

func (h *CLIHandler) readItemName() (string, error) {
	fmt.Fprintln(h.writer, "Введите наименование предмета:")

	for {
		if err := h.readLine(); err != nil {
			return "", fmt.Errorf("read line: %w", err)
		}

		name := strings.TrimSpace(h.scanner.Text())
		if name == "" {
			fmt.Fprintln(h.writer, "Ошибка: название не может быть пустым")
		} else {
			return name, nil
		}
	}
}

func (h *CLIHandler) readSimplicity() (int, error) {
	fmt.Fprintln(h.writer, "Введите простоту для новичка [1, 10]:")

	for {
		if err := h.readLine(); err != nil {
			return 0, fmt.Errorf("read line: %w", err)
		}

		simplicity, err := strconv.Atoi(strings.TrimSpace(h.scanner.Text()))
		if err != nil || simplicity < 1 || simplicity > 10 {
			fmt.Fprintln(h.writer, "Ошибка: введите целое число от 1 до 10")
		} else {
			return simplicity, nil
		}
	}
}

func (h *CLIHandler) readEnergyConsumption() (int, error) {
	fmt.Fprintln(h.writer, "Введите целое число - суточный расход энергии (кВт·ч):")

	for {
		if err := h.readLine(); err != nil {
			return 0, fmt.Errorf("read line: %w", err)
		}

		energy, err := strconv.Atoi(strings.TrimSpace(h.scanner.Text()))
		if err != nil || energy < 0 {
			fmt.Fprintln(h.writer, "Ошибка: введите неотрицательное целое число")
		} else {
			return energy, nil
		}
	}
}

func (h *CLIHandler) readLine() error {
	if !h.scanner.Scan() {
		fmt.Fprintf(h.writer, "Ошибка чтения")
		return fmt.Errorf("scan: %w", h.scanner.Err())
	}
	return nil
}

func (h *CLIHandler) handleTotalDailyEnergy() {
	energy, err := h.inventory.TotalDailyEnergy()
	if err != nil {
		fmt.Fprintln(h.writer, "Ошибка:", err)
		return
	}
	fmt.Fprintf(h.writer, "Потребление энергии: %d\n", energy)
}

func (h *CLIHandler) handleBeginnerVehicles() {
	vehicles, err := h.inventory.BeginnerVehicles()
	if err != nil {
		fmt.Fprintln(h.writer, "Ошибка:", err)
		return
	}

	fmt.Fprintf(h.writer, "%-17s | %s | %s\n", "Инвентарный номер", "Простота для новичка", "Название")
	for _, vehicle := range vehicles {
		fmt.Fprintf(h.writer, "%-17s | %-20d | %s\n", vehicle.InventoryNumber(), vehicle.Simplicity(), vehicle.Name())
	}
}

func (h *CLIHandler) handleAllItems() {
	items, err := h.inventory.AllItems()
	if err != nil {
		fmt.Fprintln(h.writer, "Ошибка:", err)
		return
	}

	fmt.Fprintf(h.writer, "%-17s | %s\n", "Инвентарный номер", "Название")
	for _, item := range items {
		fmt.Fprintf(h.writer, "%-17s | %s\n", item.InventoryNumber(), item.Name())
	}
}
