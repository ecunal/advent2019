package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
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

	set := []int{5, 6, 7, 8, 9} // part1: set := []int{0, 1, 2, 3, 4}
	current := make([]int, 0)
	perms := make([][]int, 0, 5*4*3*2)
	perms = generatePermutation(set, current, perms)

	maxOutput := 0

	for _, n := range perms {
		log.Println(n)

		amfis := make([]*amfi, len(n))
		for i := range amfis {
			amfis[i] = &amfi{
				name:  'A' + byte(i),
				close: make(chan bool, 1),
			}
		}

		for i := range amfis {
			bridge := make(chan int, 1)
			nexti := (i + 1) % len(amfis)
			amfis[i].output = bridge
			amfis[nexti].input = bridge
			bridge <- n[nexti]
		}

		wg := &sync.WaitGroup{}

		for i := range amfis {
			wg.Add(1)
			go amfis[i].run(program, wg)
		}

		// give first amfi its first input
		amfis[len(amfis)-1].output <- 0

		wg.Wait()

		if output := amfis[len(amfis)-1].lastoutput; output > maxOutput {
			maxOutput = output
		}
	}

	log.Println("Max signal:", maxOutput)
}

func generatePermutation(set, current []int, result [][]int) [][]int {
	if len(set) == 0 {
		return append(result, current)
	}

	for i := range set {
		current = append(current, set[i])

		copyset := make([]int, len(set))
		copy(copyset, set)
		copyset = append(copyset[:i], copyset[i+1:]...)

		result = generatePermutation(copyset, current, result)
		current = current[:len(current)-1]
	}

	return result
}

type amfi struct {
	name       byte
	input      <-chan int
	output     chan<- int
	close      chan bool
	lastoutput int
}

func (a *amfi) run(program []int, wg *sync.WaitGroup) {
	for pc := 0; program[pc] != 99; {
		opcode := program[pc] % 100
		switch opcode {
		case 1: // add
			a, b := getTwoInputs(program, pc)
			program[program[pc+3]] = a + b
			pc += 4
		case 2: // mult
			a, b := getTwoInputs(program, pc)
			program[program[pc+3]] = a * b
			pc += 4
		case 3: // input

			select {
			case in := <-a.input:
				program[program[pc+1]] = in
				pc += 2
			case <-a.close:
				log.Printf("Amfi %s closing.", string(a.name))
				return
			}

		case 4: // output
			a.lastoutput = program[program[pc+1]]
			log.Printf("Amfi %s Output: %d", string(a.name), a.lastoutput)
			a.output <- a.lastoutput
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

	// halt
	log.Printf("Amfi %s halted.", string(a.name))
	wg.Done()
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
