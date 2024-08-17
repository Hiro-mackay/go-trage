package main

import (
	"fmt"
)

type Vertex struct {
	X, Y int
}

func (v Vertex) Plus() int {
	return v.X + v.Y
}

func (v Vertex) String() string {
	return fmt.Sprintf("X is %v! Y is %v!\n", v.X, v.Y)
}

func q1() {
	v := Vertex{3, 4}
	fmt.Println(v.Plus())
}

func q2() {
	v := Vertex{3, 4}
	fmt.Println(v)
}

func main() {
	q1()
	q2()
}
