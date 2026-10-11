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
		fmt.Fprintln(h.writer, "\n============Меню управления============\n",
			"1. Добавить транспорт в парк\n",
			"2. Добавить вещь в парк\n",
			"3. Вывести общую дневную энергию парка\n",
			"4. Вывести предметы для новичков\n",
			"5. Вывести все предметы\n",
			"6. Выход\n",
		)

		choice, ok := h.readChoice(6)
		if !ok {
			return
		}

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
		default:
			fmt.Fprintln(h.writer, "Ошибка: неправильный тип")
		}
	}
}

func (h *CLIHandler) handleAddVehicleToFleet() {
	fmt.Fprintln(h.writer, "Выберите тип предмета:\n",
		"1. Самокат\n",
		"2. Велосипед\n",
		"3. Электросамокат\n",
		"4. Электровелосипед\n")

	choice, ok := h.readChoice(4)
	if !ok {
		return
	}
	name, ok := h.readItemName()
	if !ok {
		return
	}
	simplicity, ok := h.readSimplicity()
	if !ok {
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
		energy, ok := h.readEnergyConsumption()
		if !ok {
			return
		}
		scooter := domain.ElectricScooter{ItemInfo: &info, Simple: simplicity}
		scooter.Energy = energy
		item = scooter

	case 4:
		energy, ok := h.readEnergyConsumption()
		if !ok {
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
	fmt.Fprintln(h.writer, "Выберите тип вещи:\n",
		"1. Шлем\n",
		"2. Зарядная станция")

	choice, ok := h.readChoice(2)
	if !ok {
		return
	}
	name, ok := h.readItemName()
	if !ok {
		return
	}

	info := domain.ItemInfo{ItemName: name}
	var item domain.Inventorizable

	switch choice {
	case 1:
		item = domain.Helmet{ItemInfo: &info}
	case 2:
		energy, ok := h.readEnergyConsumption()
		if !ok {
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

func (h *CLIHandler) readChoice(choiceCount int) (int, bool) {
	if !h.scanner.Scan() {
		fmt.Fprintln(h.writer, "Ошибка чтения:", h.scanner.Err())
		return 0, false
	}

	choice, err := strconv.Atoi(strings.TrimSpace(h.scanner.Text()))
	if err != nil || choice < 1 || choice > choiceCount {
		fmt.Fprintf(h.writer, "Ошибка: введите целое число от 1 до %d\n", choiceCount)
		return 0, false
	}
	return choice, true
}

func (h *CLIHandler) readItemName() (string, bool) {
	fmt.Fprintln(h.writer, "Введите наименование предмета:")
	if !h.scanner.Scan() {
		fmt.Fprintln(h.writer, "Ошибка чтения:", h.scanner.Err())
		return "", false
	}

	name := strings.TrimSpace(h.scanner.Text())
	if name == "" {
		fmt.Fprintln(h.writer, "Ошибка: название не может быть пустым")
		return "", false
	}
	return name, true
}

func (h *CLIHandler) readSimplicity() (int, bool) {
	fmt.Fprintln(h.writer, "Введите простоту для новичка [1, 10]:")
	if !h.scanner.Scan() {
		fmt.Fprintln(h.writer, "Ошибка чтения:", h.scanner.Err())
		return 0, false
	}

	simplicity, err := strconv.Atoi(strings.TrimSpace(h.scanner.Text()))
	if err != nil || simplicity < 1 || simplicity > 10 {
		fmt.Fprintln(h.writer, "Ошибка: введите целое число от 1 до 10")
		return 0, false
	}
	return simplicity, true
}

func (h *CLIHandler) readEnergyConsumption() (int, bool) {
	fmt.Fprintln(h.writer, "Введите целое число - суточный расход энергии (кВт·ч):")
	if !h.scanner.Scan() {
		fmt.Fprintln(h.writer, "Ошибка чтения:", h.scanner.Err())
		return 0, false
	}

	energy, err := strconv.Atoi(strings.TrimSpace(h.scanner.Text()))
	if err != nil || energy < 0 {
		fmt.Fprintln(h.writer, "Ошибка: введите неотрицательное целое число")
		return 0, false
	}
	return energy, true
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

	fmt.Fprintf(h.writer, "%-17s | %-20s | %s\n", "Инвентарный номер", "Название", "Простота для новичка")
	for _, vehicle := range vehicles {
		fmt.Fprintf(h.writer, "%-17s | %-20s | %d\n", vehicle.InventoryNumber(), vehicle.Name(), vehicle.Simplicity())
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
