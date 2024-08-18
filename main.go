package main

import (
	"context"
	"fmt"
	"time"
)

func process(ctx context.Context, ch chan string) {
	fmt.Println("Processing...")
	time.Sleep(2 * time.Second)
	fmt.Println("Done!")
	ch <- "Done"
}
func main() {
	ch := make(chan string)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	go process(ctx, ch)
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		case v := <-ch:
			fmt.Println("Received value: ", v)
			return
		}
	}
}
