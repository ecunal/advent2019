package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	inputs := parseInput()

	total := 0

	for _, i := range inputs {

		tmp := fuel(i)

		for tmp > 0 {
			total += tmp
			tmp = fuel(tmp)
		}

	}

	log.Println("Total:", total)
}

func fuel(i int) int {
	return (i / 3) - 2
}

func parseInput() []int {
	f, err := os.Open("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}
	defer f.Close()

	inputList := make([]int, 0)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		input := scanner.Text()
		if f, err := strconv.Atoi(input); err != nil {
			log.Printf("could not parse input string (%s)\n", input)
		} else {
			inputList = append(inputList, f)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}
