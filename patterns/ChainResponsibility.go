package main

import "fmt"

// Request представляет заявку на обработку
type Request struct {
	Level   int
	Message string
}

// Handler определяет интерфейс обработчика
type Handler interface {
	SetNext(handler Handler)
	Handle(request *Request)
}

// BaseHandler представляет базовый обработчик с реализацией метода SetNext
type BaseHandler struct {
	next Handler
}

func (b *BaseHandler) SetNext(handler Handler) {
	b.next = handler
}

func (b *BaseHandler) Handle(request *Request) {
	if b.next != nil {
		b.next.Handle(request)
	}
}

// LowLevelSupport представляет обработчик заявок низкого уровня
type LowLevelSupport struct {
	BaseHandler
}

func (h *LowLevelSupport) Handle(request *Request) {
	if request.Level == 1 {
		fmt.Println("LowLevelSupport: Handling request -", request.Message)
	} else {
		fmt.Println("LowLevelSupport: Passing request to the next handler")
		h.BaseHandler.Handle(request)
	}
}

// MidLevelSupport представляет обработчик заявок среднего уровня
type MidLevelSupport struct {
	BaseHandler
}

func (h *MidLevelSupport) Handle(request *Request) {
	if request.Level == 2 {
		fmt.Println("MidLevelSupport: Handling request -", request.Message)
	} else {
		fmt.Println("MidLevelSupport: Passing request to the next handler")
		h.BaseHandler.Handle(request)
	}
}

// HighLevelSupport представляет обработчик заявок высокого уровня
type HighLevelSupport struct {
	BaseHandler
}

func (h *HighLevelSupport) Handle(request *Request) {
	if request.Level == 3 {
		fmt.Println("HighLevelSupport: Handling request -", request.Message)
	} else {
		fmt.Println("HighLevelSupport: Passing request to the next handler")
		h.BaseHandler.Handle(request)
	}
}

// Основная функция
func main() {
	// Создаем обработчики
	lowLevel := &LowLevelSupport{}
	midLevel := &MidLevelSupport{}
	highLevel := &HighLevelSupport{}

	// Формируем цепочку обработчиков
	lowLevel.SetNext(midLevel)
	midLevel.SetNext(highLevel)

	// Создаем заявки
	requests := []*Request{
		{Level: 1, Message: "Password reset"},
		{Level: 2, Message: "Software installation"},
		{Level: 3, Message: "Server down"},
		{Level: 4, Message: "Unknown issue"},
	}

	// Обрабатываем заявки
	for _, request := range requests {
		lowLevel.Handle(request)
	}
}
