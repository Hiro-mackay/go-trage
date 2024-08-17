package main

import (
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("./main.go")
	defer file.Close()
	data := make([]byte, 100)
	file.Read(data)
	fmt.Println(string(data))
}
