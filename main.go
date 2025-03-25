package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func mailing(name string, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range ch {
		fmt.Printf("[%s] Сообщение отправлено: %s\n", name, msg)
		time.Sleep(time.Second)
	}
}

func department(name string, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("[%s] Уведомление #%d", name, i+1)
		out <- msg
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000))) // Разное время отправки
	}
	close(out)
}

func fanIn(inputs ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	for _, ch := range inputs {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for msg := range c {
				out <- msg
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	contacts := []string{"Арни", "Посвеж", "Жанна", "Кореец", "Жонас(Брат)", "Игорь Ким"}
	message := []string{"Осторожно, возможно воспламенение, приготовьте огнетушитель!", "Произошел жирный доезд!!!", "Адилька сгорел"}

	var wg sync.WaitGroup
	messages := make([]chan string, len(contacts))

	for i, name := range contacts {
		messages[i] = make(chan string, len(messages))
		wg.Add(1)
		go mailing(name, messages[i], &wg)
	}

	for _, msg := range message {
		fmt.Println("🔹 Рассылаем:", msg)
		for _, ch := range messages {
			ch <- msg
		}
	}

	for _, ch := range messages {
		close(ch)
	}

	wg.Wait()
	fmt.Println("Прежупреждение отправлено")

	rand.Seed(time.Now().UnixNano())
	
	dept1 := make(chan string)
	dept2 := make(chan string)
	dept3 := make(chan string)

	wg.Add(3)
	go department(" Отдел продаж", dept1, &wg)
	go department(" Отдел HR", dept2, &wg)
	go department(" Техподдержка", dept3, &wg)

	notifications := fanIn(dept1, dept2, dept3)

	for msg := range notifications {
		fmt.Println(" Получено уведомление:", msg)
	}

	wg.Wait()
	fmt.Println(" Все отделы завершили рассылку")
}
