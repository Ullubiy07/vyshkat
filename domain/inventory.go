package domain

type InventoryNumber struct {
	ID string
}

func NewInventoryNumber(id string) InventoryNumber {
	return InventoryNumber{ID: id}
}

func (n InventoryNumber) String() string {
	return n.ID
}

type Inventorizable interface {
	Name() string
	InventoryNumber() InventoryNumber
	SetInventoryNumber(num InventoryNumber)
}

type Vehicle interface {
	Inventorizable
	Simplicity() int
}

type EnergyConsumer interface {
	EnergyConsumption() int
}

type ItemInfo struct {
	ItemName string
	ItemNum  InventoryNumber
}

func (i *ItemInfo) Name() string {
	return i.ItemName
}

func (i *ItemInfo) InventoryNumber() InventoryNumber {
	return i.ItemNum
}

func (i *ItemInfo) SetInventoryNumber(num InventoryNumber) {
	i.ItemNum = num
}

type ElectricBike struct {
	*ItemInfo
	Energy int
	Simple int
}

func (b ElectricBike) EnergyConsumption() int {
	return b.Energy
}

func (b ElectricBike) Simplicity() int {
	return b.Simple
}

type ElectricScooter struct {
	*ItemInfo
	Energy int
	Simple int
}

func (s ElectricScooter) Simplicity() int {
	return s.Simple
}

func (s ElectricScooter) EnergyConsumption() int {
	return s.Energy
}

type Bike struct {
	*ItemInfo
	Simple int
}

func (b Bike) Simplicity() int {
	return b.Simple
}

type Scooter struct {
	*ItemInfo
	Simple int
}

func (s Scooter) Simplicity() int {
	return s.Simple
}

type Helmet struct {
	*ItemInfo
}

type ChargingStation struct {
	*ItemInfo
	Energy int
	Simple int
}

func (s ChargingStation) Simplicity() int {
	return s.Simple
}

func (s ChargingStation) EnergyConsumption() int {
	return s.Energy
}
