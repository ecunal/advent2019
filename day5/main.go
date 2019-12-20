package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	inputstr, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}

	inputstrings := strings.Split(strings.TrimSpace(string(inputstr)), ",")
	program := make([]int, len(inputstrings))

	for i := range inputstrings {
		program[i], err = strconv.Atoi(inputstrings[i])
		if err != nil {
			log.Fatal("Error converting str to int:", err)
		}
	}

	input := 5 // part 2

	for pc := 0; program[pc] != 99; {

		opcode := program[pc] % 100

		log.Println("Opcode:", opcode, "PC:", pc)

		switch opcode {
		case 1: // add

			a, b := getTwoInputs(program, pc)

			//	log.Println("Add:", a, b)

			program[program[pc+3]] = a + b

			pc += 4
		case 2: // mult

			a, b := getTwoInputs(program, pc)

			// log.Println("Mult:", a, b)

			program[program[pc+3]] = a * b

			pc += 4
		case 3: // input
			program[program[pc+1]] = input
			pc += 2
		case 4: // output
			log.Println("Output:", program[program[pc+1]])
			pc += 2
		case 5: // jump-if-true

			a, b := getTwoInputs(program, pc)
			if a != 0 {
				pc = b
			} else {
				pc += 3
			}

		case 6: // jump-if-false

			a, b := getTwoInputs(program, pc)
			if a == 0 {
				pc = b
			} else {
				pc += 3
			}

		case 7: // less than

			a, b := getTwoInputs(program, pc)
			if a < b {
				program[program[pc+3]] = 1
			} else {
				program[program[pc+3]] = 0
			}

			pc += 4

		case 8: // equals

			a, b := getTwoInputs(program, pc)
			if a == b {
				program[program[pc+3]] = 1
			} else {
				program[program[pc+3]] = 0
			}

			pc += 4

		default:
			log.Print("Error:", opcode)
		}
	}
}

func getTwoInputs(program []int, pc int) (int, int) {
	a, b := program[pc+1], program[pc+2]

	if (program[pc]/100)%10 == 0 {
		a = program[a]
	}

	if (program[pc]/1000)%10 == 0 {
		b = program[b]
	}
	return a, b
}
