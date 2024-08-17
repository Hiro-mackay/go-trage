package main

import (
	"fmt"
)

func connectDB() {
	panic("Unable to connect database!")
}

func save() {
	defer func() {
		s := recover()
		fmt.Println(s)
	}()
	connectDB()
}
func main() {
	save()
	fmt.Println("OK?")

}
