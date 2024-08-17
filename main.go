package main

func Producer(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}

func double(ch1 chan int, ch2 chan int) {
	defer close(ch2)
	for i := range ch1 {
		ch2 <- i * 2
	}
}

func quadruple(ch2 chan int, ch3 chan int) {
	defer close(ch3)
	for i := range ch2 {
		ch3 <- i * 4
	}
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go Producer(ch1)
	go double(ch1, ch2)
	go quadruple(ch2, ch3)
	for r := range ch3 {
		println(r)
	}

}
