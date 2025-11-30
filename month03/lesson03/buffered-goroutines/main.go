package main

import (
	"fmt"
	"time"
)

type Task struct {
	ID int
	TimeCreated string
}

func worker(id int, tasks <-chan Task) {
	for task := range tasks {
		fmt.Printf("Воркер %d начал задачу %d (создана в %s)\n", id, task.ID, task.TimeCreated)
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Воркер %d закончил задачу %d\n", id, task.ID)
	}
}

func main() {
	taskCh := make(chan Task, 5)
	
	go worker(1, taskCh)
	go worker(2, taskCh)
	go worker(3, taskCh)

	for i := 1; i <= 15; i++ {
		//timeCreated := time.Now()
		task := Task{
			ID: i,
			TimeCreated: time.Now().Format("15:04:05.000"),
		}
		taskCh <- task
		fmt.Printf("Отправлена задача %d\n", i)
		time.Sleep(100 * time.Millisecond)
	}
	close(taskCh)
	
	time.Sleep(8 * time.Second)
	fmt.Println("Работа окончена")
}