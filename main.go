package main

import (
	"fmt"
	"sync"
	"time"
)

func Producer(ch chan int, i int) {
	ch <- i * 2
}

func Consumer(ch chan int, wg *sync.WaitGroup) {
	for i := range ch {
		func() {
			defer wg.Done()
			fmt.Println("process", i*1000)
		}()
	}

	fmt.Println("chanel closed")

}

func main() {

	var wg sync.WaitGroup
	ch := make(chan int)

	// Producer
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go Producer(ch, i)
	}

	// Consumer
	go Consumer(ch, &wg)
	wg.Wait()
	close(ch)

	time.Sleep(1 * time.Second)
	fmt.Println("DONE")

}
