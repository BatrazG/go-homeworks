package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

type PhoneBook struct {
	record map[string]string
	mu     sync.RWMutex
}

func NewPhoneBook() *PhoneBook {
	return &PhoneBook{
		record: make(map[string]string),
		// RWMutex не требует инициализации, его нулевое значение готово к работе
	}
}

func (p *PhoneBook) Set(name, phone string) {
	p.mu.Lock()
	p.record[name] = phone
	p.mu.Unlock()
}

func (p *PhoneBook) Get(name string) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	phone, ok := p.record[name] // Идиома "comma ok" для проверки наличия в map
	return phone, ok
}

func main() {
	var wg sync.WaitGroup
	book := NewPhoneBook()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			lastDigits := rand.IntN(100)
			phone := fmt.Sprintf("+79999999%02d", lastDigits)
			name := fmt.Sprintf("record %d", i)
			book.Set(name, phone)
		}
	}()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				name := fmt.Sprintf("record %d", rand.IntN(100))
				phone, ok := book.Get(name)
				if ok {
					fmt.Printf("Reader %d: %s -> %s\n", id+1, name, phone)
				} else {
					fmt.Printf("Reader %d: record %s not found\n", id+1, name)
				}
			}
		}(i)
	}

	wg.Wait()
}
