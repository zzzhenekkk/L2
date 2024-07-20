package facade

import (
	"fmt"
	"net/http"
)

//Паттерн фасад используется для упрощения интерфейса сложных систем. Он предоставляет единый интерфейс для набора интерфейсов в подсистеме,
//тем самым облегчая использование системы. Основные преимущества:
//
//1.	Упрощение использования системы: Фасад скрывает сложность взаимодействия с подсистемами.
//2.	Снижение зависимости от деталей реализации: Клиенты взаимодействуют с фасадом, а не с отдельными компонентами, что снижает количество зависимостей.
//3.	Повышение читабельности кода: Фасад делает код более читаемым и поддерживаемым, особенно если система имеет множество компонентов.

// OrderService предоставляет методы для работы с заказами
type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(userID, productID string) (string, error) {
	fmt.Printf("Creating order for user %s and product %s\n", userID, productID)
	// Логика создания заказа
	return "new_order_id", nil
}

// PaymentService предоставляет методы для обработки оплаты
type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) ProcessPayment(orderID, paymentMethod string) error {
	fmt.Printf("Processing payment for order %s with method %s\n", orderID, paymentMethod)
	// Логика обработки оплаты
	return nil
}

// NotificationService предоставляет методы для отправки уведомлений
type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendOrderConfirmation(userID, orderID string) error {
	fmt.Printf("Sending order confirmation to user %s for order %s\n", userID, orderID)
	// Логика отправки уведомления
	return nil
}

// OrderFacade объединяет взаимодействие с различными сервисами для обработки заказа
type OrderFacade struct {
	orderService        *OrderService
	paymentService      *PaymentService
	notificationService *NotificationService
}

func NewOrderFacade() *OrderFacade {
	return &OrderFacade{
		orderService:        NewOrderService(),
		paymentService:      NewPaymentService(),
		notificationService: NewNotificationService(),
	}
}

func (f *OrderFacade) PlaceOrder(userID, productID, paymentMethod string) error {
	// 1. Создание заказа
	orderID, err := f.orderService.CreateOrder(userID, productID)
	if err != nil {
		return err
	}

	// 2. Обработка оплаты
	err = f.paymentService.ProcessPayment(orderID, paymentMethod)
	if err != nil {
		return err
	}

	// 3. Отправка уведомления
	err = f.notificationService.SendOrderConfirmation(userID, orderID)
	if err != nil {
		return err
	}

	return nil
}

// Основная функция, где создается HTTP сервер и обрабатываются запросы на размещение заказа
func main() {
	http.HandleFunc("/place_order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		userID := r.FormValue("user_id")
		productID := r.FormValue("product_id")
		paymentMethod := r.FormValue("payment_method")

		orderFacade := NewOrderFacade()
		err := orderFacade.PlaceOrder(userID, productID, paymentMethod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Order placed successfully")
	})

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
