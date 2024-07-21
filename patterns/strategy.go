package strategy

import "fmt"

// PaymentMethod определяет интерфейс стратегии оплаты
type PaymentMethod interface {
	Pay(amount float64) string
}

// CreditCard реализует стратегию оплаты с использованием кредитной карты
type CreditCard struct{}

func (c *CreditCard) Pay(amount float64) string {
	return fmt.Sprintf("%.2f paid using Credit Card", amount)
}

// PayPal реализует стратегию оплаты через PayPal
type PayPal struct{}

func (p *PayPal) Pay(amount float64) string {
	return fmt.Sprintf("%.2f paid using PayPal", amount)
}

// PaymentContext содержит контекст оплаты
type PaymentContext struct {
	Method PaymentMethod
}

func (p *PaymentContext) Pay(amount float64) {
	fmt.Println(p.Method.Pay(amount))
}

// Основная функция
func main() {
	// Создание контекста оплаты с использованием кредитной карты
	creditCardPayment := &PaymentContext{
		Method: &CreditCard{},
	}
	creditCardPayment.Pay(100.0)

	// Смена стратегии оплаты на PayPal
	paypalPayment := &PaymentContext{
		Method: &PayPal{},
	}
	paypalPayment.Pay(200.0)
}
