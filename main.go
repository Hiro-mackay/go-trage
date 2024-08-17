package main

import "time"

func getOsName() string {
	return "Linux"
}

func main() {
	switch os := getOsName(); os {
	case "Linux":
		println("Linux")
	case "Windows":
		println("Windows")
	default:
		println("Unknown")
	}

	switch {
	case time.Now().Hour() < 12:
		println("Good morning!")
	case time.Now().Hour() < 17:
		println("Good afternoon.")
	default:
		println("Good evening.")
	}
}
