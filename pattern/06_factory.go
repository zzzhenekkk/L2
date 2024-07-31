package main

import "fmt"

// Transport определяет интерфейс продукта
type Transport interface {
	Deliver() string
}

// Car представляет продукт "Автомобиль"
type Car struct{}

func (c *Car) Deliver() string {
	return "Delivering by driving a car"
}

// Bike представляет продукт "Велосипед"
type Bike struct{}

func (b *Bike) Deliver() string {
	return "Delivering by riding a bike"
}

// Scooter представляет продукт "Самокат"
type Scooter struct{}

func (s *Scooter) Deliver() string {
	return "Delivering by riding a scooter"
}

// TransportFactory определяет интерфейс фабрики
type TransportFactory interface {
	CreateTransport() Transport
}

// CarFactory реализует фабрику для создания автомобилей
type CarFactory struct{}

func (cf *CarFactory) CreateTransport() Transport {
	return &Car{}
}

// BikeFactory реализует фабрику для создания велосипедов
type BikeFactory struct{}

func (bf *BikeFactory) CreateTransport() Transport {
	return &Bike{}
}

// ScooterFactory реализует фабрику для создания самокатов
type ScooterFactory struct{}

func (sf *ScooterFactory) CreateTransport() Transport {
	return &Scooter{}
}

// Основная функция
func main() {
	carFactory := &CarFactory{}
	bikeFactory := &BikeFactory{}
	scooterFactory := &ScooterFactory{}

	// Создаем транспортные средства
	car := carFactory.CreateTransport()
	bike := bikeFactory.CreateTransport()
	scooter := scooterFactory.CreateTransport()

	// Демонстрируем доставку
	fmt.Println(car.Deliver())
	fmt.Println(bike.Deliver())
	fmt.Println(scooter.Deliver())
}
