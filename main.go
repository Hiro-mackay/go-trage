package main

import (
	"encoding/json"
	"fmt"
)

type T struct{}

type Person struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Nicknames string `json:"nicknames"`
	T         *T     `json:"T,omitempty"`
}

func (p *Person) UnmarshalJSON(b []byte) error {
	type PersonAlias Person
	pa := &PersonAlias{
		T: &T{},
	}
	if err := json.Unmarshal(b, pa); err != nil {
		return err
	}
	*p = Person(*pa)
	p.Name = p.Name + "!"
	return nil
}

func (p *Person) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name string
	}{
		Name: "Mr." + p.Name,
	})
}

func main() {
	b := []byte(`{"name":"John","age":30,"nicknames":"Johny"}`)
	var p Person
	if err := json.Unmarshal(b, &p); err != nil {
		fmt.Println(err)
	}

	fmt.Println(p.Name, p.Age, p.Nicknames)
	v, _ := json.Marshal(p)
	fmt.Println(string(v))

}
