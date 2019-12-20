package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile)
	inputstr, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}

	inputstrings := strings.Split(strings.TrimSpace(string(inputstr)), ",")
	stack := make(map[int]int)

	for i := range inputstrings {
		stack[i], err = strconv.Atoi(inputstrings[i])
		if err != nil {
			log.Fatal("Error converting str to int:", err)
		}
	}

	pInput := make(chan int, 1)
	pOutput := make(chan int)
	close := make(chan bool)
	prog := &program{
		relativeBase: 0,
		input:        pInput,
		output:       pOutput,
		close:        close,
		program:      stack,
	}
	go prog.run()

	curr := pair{0, 0}
	currDirection := 0
	grid := make(map[pair]int)

	pInput <- 1 // part2

	for {
		// log.Println(curr)
		select {
		case o := <-pOutput:
			if o != 0 && o != 1 {
				log.Fatal("Error in output:", o)
			}
			/* part1
			_, ok := grid[curr]
			if !ok {
				//log.Println("Painting", curr, "Total:", painted)
				painted++
			}
			*/
			grid[curr] = o

			turn := <-pOutput
			if turn == 0 {
				// log.Println("Turning left")
				newDirection := currDirection - 1
				if newDirection < 0 {
					newDirection += len(directions)
				}
				currDirection = newDirection
			} else if turn == 1 {
				// log.Println("Turning right")
				currDirection = (currDirection + 1) % len(directions)
			} else {
				log.Fatal("Error in direction:", turn)
			}

			curr.i += directions[currDirection].i
			curr.j += directions[currDirection].j

			pInput <- grid[curr]
		case <-close:

			cols, rows := 0, 0

			for k := range grid {
				if k.i > cols {
					cols = k.i
				}
				if k.j > rows {
					rows = k.j
				}
			}

			for i := 0; i <= cols; i++ {
				for j := 0; j <= rows; j++ {
					if grid[pair{i, j}] == 1 {
						fmt.Print("#")
					} else {
						fmt.Print(" ")
					}
				}
				fmt.Println()
			}
			fmt.Println()
			os.Exit(0)
		}
	}
}

type pair struct {
	i, j int
}

var directions = []pair{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

type program struct {
	relativeBase int
	input        <-chan int
	output       chan<- int
	close        chan bool
	program      map[int]int
}

func (p *program) run() {
	for pc := 0; p.program[pc] != 99; {
		opcode := p.program[pc] % 100
		// log.Println("Opcode:", opcode)
		// log.Println("PC:", pc)
		switch opcode {
		case 1: // add
			a, b, o := p.get3(pc)
			p.program[o] = a + b
			pc += 4
		case 2: // mult
			a, b, o := p.get3(pc)
			p.program[o] = a * b
			pc += 4
		case 3: // input
			o := p.getOutputAddr(pc, 1)
			p.program[o] = <-p.input
			pc += 2
		case 4: // output
			o := p.getInputValue(pc, 1)
			//log.Println("Output:", o)
			p.output <- o
			pc += 2
		case 5: // jump-if-true
			a, b := p.getTwoInputs(pc)
			if a != 0 {
				pc = b
			} else {
				pc += 3
			}
		case 6: // jump-if-false
			a, b := p.getTwoInputs(pc)
			if a == 0 {
				pc = b
			} else {
				pc += 3
			}
		case 7: // less than
			a, b, o := p.get3(pc)
			if a < b {
				p.program[o] = 1
			} else {
				p.program[o] = 0
			}
			pc += 4
		case 8: // equals
			a, b, o := p.get3(pc)
			if a == b {
				p.program[o] = 1
			} else {
				p.program[o] = 0
			}
			pc += 4
		case 9:
			a := p.getInputValue(pc, 1)
			p.relativeBase += a
			pc += 2
		default:
			log.Print("Error:", opcode)
		}
	}

	// halt
	log.Printf("EOF")
	close(p.close)
}

func (p *program) getTwoInputs(pc int) (int, int) {
	return p.getInputValue(pc, 1), p.getInputValue(pc, 2)
}

func (p *program) getInputValue(pc, offset int) int {
	a := p.program[pc+offset]
	if mode := (p.program[pc] / (100 * int(math.Pow10(offset-1)))) % 10; mode == 0 { // position
		a = p.program[a]
	} else if mode == 2 { // relative
		a = p.program[a+p.relativeBase]
	}
	return a
}

func (p *program) getOutputAddr(pc, offset int) int {
	o := p.program[pc+offset]
	if mode := (p.program[pc] / (100 * int(math.Pow10(offset-1)))) % 10; mode == 2 {
		o += p.relativeBase
	}
	return o
}

func (p *program) get3(pc int) (int, int, int) {
	a, b := p.getTwoInputs(pc)
	c := p.getOutputAddr(pc, 3)
	return a, b, c
}
