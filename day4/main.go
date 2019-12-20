package main

import "log"

func main() {
	from, to := 152085, 670283
	//from, to := 223450, 223460
	no := 0

	for i := from; i <= to; i++ {
		d := digits(i)
		if isValid(d) {
			no++
		}
	}

	log.Println(no)
}

func isValid(d []int) bool {
	doubles := false

	for j := len(d) - 1; j > 0; j-- {
		if d[j] > d[j-1] {
			return false
		} else if d[j] == d[j-1] {
			if !((j+1 < len(d) && d[j] == d[j+1]) || (j-2 >= 0 && d[j] == d[j-2])) {
				doubles = true
			}
		}
	}

	return doubles
}

func digits(n int) []int {
	result := make([]int, 0)
	for x := n; x > 0; x /= 10 {
		result = append(result, x%10)
	}
	return result
}
