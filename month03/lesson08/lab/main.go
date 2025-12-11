package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	Amount int    // Сумма заказа
	Status string // "new", "processed", "paid"
}

type Processable interface {
	Process()
}

func (o *Order) Process() {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	o.Status = "processed"
}

// generateOrders создает count заказов со статусом "new" и случайной суммой
func generateOrders(count int) <-chan Order {
	out := make(chan Order)

	go func() {
		defer close(out)
		for i := 1; i <= count; i++ {
			amount := rand.Intn(10_000)
			order := Order{
				ID:     i,
				Amount: amount,
				Status: "new",
			}
			out <- order
		}
	}()

	return out
}

// Обработчик (fan-out)
func processOrders(in <-chan Order, workers int) <-chan Order {
	out := make(chan Order)
	var wg sync.WaitGroup

	worker := func(id int) {
		defer wg.Done()
		for order := range in {
			//Используем интерфейс
			var p Processable = &order
			p.Process()

			out <- order
		}
	}

	wg.Add(workers)
	for i := 1; i <= workers; i++ {
		go worker(i)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Фильтр
func filterOrders(in <-chan Order, minAmount int) <-chan Order {
	out := make(chan Order)

	go func() {
		defer close(out)
		for order := range in {
			if order.Amount > minAmount {
				out <- order
			}
		}
	}()

	return out
}

func main() {
	ordersChan := generateOrders(50)

	processedChan := processOrders(ordersChan, 5)

	filteredChan := filterOrders(processedChan, 5000)

	for order := range filteredChan {
		fmt.Printf("Итоговый заказ: ID=%d, Amount=%d, Status=%s\n",
			order.ID, order.Amount, order.Status)
	}

}
