package main

import (
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	inputStr, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatal(err)
	}
	input := make([]int, (len(inputStr)-1)*10000)
	for j := 0; j < 10000; j++ {
		for i := 0; i < len(inputStr)-1; i++ {
			input[(j*(len(inputStr)-1))+i] = int(inputStr[i] - '0')
		}
	}

	offset, err := strconv.Atoi(string(inputStr[:7]))
	if err != nil {
		log.Fatal(err)
	}

	for p := 0; p < 100; p++ {

		newInput := make([]int, len(input))

		sum := 0

		for i := len(input) - 1; i >= offset; i-- {
			sum += input[i]
			newInput[i] = abs(sum) % 10
		}
		input = newInput
	}

	log.Println(input[offset : offset+8])
}

func part1() {
	inputStr, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatal(err)
	}
	input := make([]int, len(inputStr)-1)
	for i := range input {
		input[i] = int(inputStr[i] - '0')
	}

	basePattern := []int{0, 1, 0, -1}

	for p := 0; p < 100; p++ {
		// phase
		//log.Printf("Phase %d, input: %v", p, input)

		newInput := make([]int, len(input))

		for n := range newInput {
			currentPattern := make([]int, len(basePattern)*(n+1))
			for i := range basePattern {
				for j := 0; j < n+1; j++ {
					currentPattern[i*(n+1)+j] = basePattern[i]
				}
			}

			//log.Println("current output element", n+1)
			//log.Println(currentPattern)

			// repeat base pattern by n+1, skip 0th element

			x := 0
			for i := range input {
				x += input[i] * currentPattern[(i+1)%len(currentPattern)] // TODO, repeating pattern
				//fmt.Printf("%d * %d + ", input[i], currentPattern[(i+1)%len(currentPattern)])
			}
			//fmt.Println()
			newInput[n] = abs(x) % 10
			//log.Println(newInput[n])
		}

		input = newInput
	}

	log.Println(input[:8])
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
