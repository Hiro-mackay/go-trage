package main

const (
	_      = iota
	KB int = 1 << (10 * iota)
	MB
	GB
)

func main() {
	println(KB, MB, GB)
}
