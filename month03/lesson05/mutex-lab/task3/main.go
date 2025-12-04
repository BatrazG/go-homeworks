package main

import (
	"fmt"
	"sync"
	"unsafe"
)

// BankAccount содержит мьютекс и баланс
type BankAccount struct {
	mu      sync.Mutex
	Balance int
}

func (a *BankAccount) GetBalance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Balance
}

// Deposit безопасно изменяет баланс.
// Получатель (a *BankAccount) — указатель.
// Это важно: мьютексы нельзя копировать по значению.
func (a *BankAccount) Deposit(amount int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance += amount
}

func (from *BankAccount) Transfer(to *BankAccount, amount int) error {
	if from == to {
		return fmt.Errorf("нельзя перевести на тот же самый счет")
	}

	//Определяем порядок блокировки, чтобы избежать дедлока
	first, second := from, to
	/*Безопасно, но долго
		if fmt.Sprintf("%p", first) > fmt.Sprintf("%p", second) {
		first, second = second, first
	}*/
	if uintptr(unsafe.Pointer(first)) > uintptr(unsafe.Pointer(second)) {
		first, second = second, first
	}

	first.mu.Lock()
	second.mu.Lock()
	defer first.mu.Unlock()
	defer second.mu.Unlock()

	if amount > from.Balance {
		return fmt.Errorf("Невозможно перевести %d, остаток: %d", amount, from.Balance)
	}

	from.Balance -= amount
	to.Balance += amount

	return nil
}

func main() {
	account := BankAccount{}
	account1 := BankAccount{Balance: 1000}
	account2 := BankAccount{Balance: 1000}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account.Deposit(1)
		}()
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := account1.Transfer(&account2, 200); err != nil {
			fmt.Println("Ошибка перевода:", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := account2.Transfer(&account1, 300); err != nil {
			fmt.Println("Ошибка перевода:", err)
		}
	}()

	wg.Wait()
	fmt.Println("Итоговый счетчик:", account.Balance)
	fmt.Println(account1.GetBalance())
	fmt.Println(account2.GetBalance())

	/*
		Чтобы избежать дедлока при одновременных операциях, мьютексы обоих
		счетов нужно блокировать всегда в одном и том же глобальном порядке.
		Самый простой способ — упорядочить их по какому-то уникальному неизменному признаку,
		например, по адресу в памяти:
		сначала блокируется мьютекс счета с меньшим адресом, затем — с большим.
	*/

}
