package main

import "fmt"

type Human interface {
	Say() string
}

type Person struct {
	Name string
}

func (p *Person) Say() string {
	p.Name = "Mr." + p.Name
	fmt.Println("Hello", p.Name)
	return p.Name
}

type Dog struct {
	Name string
}

func (d Dog) Say() string {
	fmt.Println("Bark")
	return "Dog"
}

func Drive(human Human) {
	if human.Say() == "Mr.Mike" {
		fmt.Println("Drive")
	} else {
		fmt.Println("Get out")
	}

}

func main() {
	var mike Human = &Person{Name: "Mike"}
	Drive(mike)
	var x Human = &Person{Name: "X"}
	Drive(x)

	var dog Dog = Dog{Name: "Dog"}
	Drive(dog) // Ops! This is not a human. But Dog struct has Say() method. So, it is a Human interface.
}
