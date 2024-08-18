package lib

func Average(s []int) int {
	total := 0
	for _, v := range s {
		total += v
	}
	return total / len(s)
}
