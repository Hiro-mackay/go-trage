package main

func cal(price, quantity int) (amount int) {
	amount = price * quantity
	return
}

func main() {
	r := cal(100, 10)
	println(r)
}
