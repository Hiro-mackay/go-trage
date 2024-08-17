package main

import (
	"fmt"
	"slices"
)

func q1() {
	l := []int{100, 300, 23, 11, 23, 2, 4, 6, 4}
	fmt.Println(slices.Min(l))
}

func q2() {
	m := map[string]int{
		"apple":  200,
		"banana": 300,
		"grapes": 150,
		"orange": 80,
		"papaya": 500,
		"kiwi":   90,
	}

	sum := 0

	for _, v := range m {
		sum += v
	}

	fmt.Println(sum)

}

func main() {
	q1()
	q2()
}
