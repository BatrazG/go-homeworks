package main

import (
	"fmt"
	"sync"
)

// BankAccount содержит мьютекс и баланс
type BankAccount struct {
	mu      sync.Mutex
	Balance int
}

// Deposit безопасно изменяет баланс.
// Получатель (a *BankAccount) — указатель.
// Это важно: мьютексы нельзя копировать по значению.
func (a *BankAccount) Deposit(amount int) {
	a.mu.Lock()
	a.Balance += amount
	a.mu.Unlock()
}

func main() {
	account := BankAccount{}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account.Deposit(1)
		}()
	}

	wg.Wait()
	fmt.Println("Итоговый счетчик:", account.Balance)
}
