package main

import (
	"fmt"
	"go-trade/lib"
)

func main() {
	s := []int{1, 2, 3, 4, 5}
	fmt.Println(lib.Average(s))

	person := lib.Person{Name: "John", Age: 25}
	fmt.Println(person)
}
