package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range jobs {
		fmt.Printf("worker %d скачивает %s\n", id, url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("worker %d: ошибОЧКА %s: %v\n", id, url, err)
			continue
		}

		fmt.Printf("Worker %d: %s -> статус %d\n", id, url, resp.StatusCode)
		resp.Body.Close()
		time.Sleep(time.Second * 3)
	}
}
func main() {
	urls := []string{
		"https://examplelhl.com",
		"https://www.google.com",
		"https://www.github.com",
		"https://www.facebook.com",
		"https://www.stackoverflow.com",
	}

	const numWorkers = 3

	jobs := make(chan string, len(urls))
	var wg sync.WaitGroup

	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobs, &wg)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	wg.Wait()
	fmt.Println("я покакаль")
}
