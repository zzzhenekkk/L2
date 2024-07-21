package visitor

import "fmt"

//Паттерн Посетитель (Visitor) — это поведенческий паттерн проектирования,
//который позволяет добавлять новые операции к объектам, не изменяя их классы.
//Он разделяет алгоритмы и структуры объектов, над которыми эти алгоритмы выполняются.
//Это достигается за счет введения нового класса (посетителя), который содержит все операции,
//выполняемые над объектами.

// Shape определяет интерфейс фигуры
type Shape interface {
	getType() string
	accept(Visitor)
}

// Square представляет конкретную фигуру: квадрат
type Square struct {
	side int
}

func (s *Square) accept(v Visitor) {
	v.visitForSquare(s)
}

func (s *Square) getType() string {
	return "Square"
}

// Circle представляет конкретную фигуру: круг
type Circle struct {
	radius int
}

func (c *Circle) accept(v Visitor) {
	v.visitForCircle(c)
}

func (c *Circle) getType() string {
	return "Circle"
}

// Rectangle представляет конкретную фигуру: прямоугольник
type Rectangle struct {
	l int
	b int
}

func (t *Rectangle) accept(v Visitor) {
	v.visitForRectangle(t)
}

func (t *Rectangle) getType() string {
	return "Rectangle"
}

// Visitor определяет интерфейс посетителя
type Visitor interface {
	visitForSquare(*Square)
	visitForCircle(*Circle)
	visitForRectangle(*Rectangle)
}

// AreaCalculator представляет конкретного посетителя для вычисления площади
type AreaCalculator struct {
	area int
}

func (a *AreaCalculator) visitForSquare(s *Square) {
	// Логика вычисления площади квадрата
	fmt.Println("Calculating area for square")
}

func (a *AreaCalculator) visitForCircle(c *Circle) {
	// Логика вычисления площади круга
	fmt.Println("Calculating area for circle")
}

func (a *AreaCalculator) visitForRectangle(r *Rectangle) {
	// Логика вычисления площади прямоугольника
	fmt.Println("Calculating area for rectangle")
}

// MiddleCoordinates представляет конкретного посетителя для вычисления средних координат
type MiddleCoordinates struct {
	x int
	y int
}

func (m *MiddleCoordinates) visitForSquare(s *Square) {
	// Логика вычисления средних координат для квадрата
	fmt.Println("Calculating middle point coordinates for square")
}

func (m *MiddleCoordinates) visitForCircle(c *Circle) {
	// Логика вычисления средних координат для круга
	fmt.Println("Calculating middle point coordinates for circle")
}

func (m *MiddleCoordinates) visitForRectangle(r *Rectangle) {
	// Логика вычисления средних координат для прямоугольника
	fmt.Println("Calculating middle point coordinates for rectangle")
}

// Основная функция
func main() {
	square := &Square{side: 2}
	circle := &Circle{radius: 3}
	rectangle := &Rectangle{l: 2, b: 3}

	areaCalculator := &AreaCalculator{}

	square.accept(areaCalculator)
	circle.accept(areaCalculator)
	rectangle.accept(areaCalculator)

	fmt.Println()
	middleCoordinates := &MiddleCoordinates{}
	square.accept(middleCoordinates)
	circle.accept(middleCoordinates)
	rectangle.accept(middleCoordinates)
}
