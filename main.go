package main

import "fmt"

func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Println("int", v)
	case string:
		fmt.Println("string", v)
	default:
		fmt.Printf("unknown type '%T'\n", v)
	}
}

func main() {
	do(21)
	do("hello")
	do(true)
}
