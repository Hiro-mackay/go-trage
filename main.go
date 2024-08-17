package main

import "fmt"

func main() {
	m := map[string][]string{
		"something": []string{"a", "b", "c"},
		"nothing":   []string{},
	}

	v, ok := m["nothing"]
	fmt.Println(v, ok)
}
