package main

import (
	"log"
	"sync"
	"time"
)

func worker(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Start working: ", n)
	time.Sleep(3 * time.Second)
	log.Println("Done: ", n)
}

func main() {
	var wg sync.WaitGroup
	log.Println("Hello world")
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
		time.Sleep(2 * time.Second)
	}
	wg.Wait()
	log.Println("Test over")
}
