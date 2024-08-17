package main

import "fmt"

type UserNotFound struct {
	Username string
}

func (e *UserNotFound) Error() string {
	return fmt.Sprintf("User not found %v", e.Username)
}

func MyFunc() error {

	ok := false

	if ok {
		return nil
	}

	return &UserNotFound{Username: "gopher"}
}
func main() {
	if err := MyFunc(); err != nil {
		fmt.Println(err)
	}

}
