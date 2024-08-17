package main

import "fmt"

func main() {
	var i int = 100
	var j int = 200
	var p1 *int
	var p2 *int
	p1 = &i // 100
	fmt.Println(*p1)
	p2 = &j // 200
	fmt.Println(*p2)
	i = *p1 + *p2 // 100 + 200 & *&i -> 300
	fmt.Println(i)
	p2 = p1 //100
	fmt.Println(*p2)
	j = *p2 + i // 200 + 300
	fmt.Println(j)

}
