package main

import (
	"io/ioutil"
	"log"
	"math"
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

	grid := make([][]bool, 2000)
	for i := range grid {
		grid[i] = make([]bool, 2000)
	}

	firstj := 0
	for i := range grid {
		beamSeen := false

		for j := firstj; j < len(grid[i]); j++ {

			copystack := make(map[int]int)
			for k, v := range stack {
				copystack[k] = v
			}

			pInput := make(chan int)
			pOutput := make(chan int)
			close := make(chan bool)
			prog := &program{
				relativeBase: 0,
				input:        pInput,
				output:       pOutput,
				close:        close,
				program:      copystack,
			}
			go prog.run()

			pInput <- i
			pInput <- j

			pulled := <-pOutput
			if pulled == 0 {
				grid[i][j] = false
				if beamSeen {
					break
				}
			} else {
				grid[i][j] = true
				if !beamSeen {
					beamSeen = true
					firstj = j
				}
			}
		}
	}

	for i, row := range grid {
		for j, tile := range row {
			if tile {
				if isSquare(i, j, grid) {
					log.Println("i:", i, "j:", j, "result:", i*10000+j)
					return
				}
			}
		}
	}
}

func isSquare(i, j int, grid [][]bool) bool {
	for x := i; x < i+100; x++ {
		if x >= len(grid) {
			return false
		}
		for y := j; y < j+100; y++ {
			if y >= len(grid[x]) {
				return false
			}
			if !grid[x][y] {
				return false
			}
		}
	}
	return true
}

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
	//log.Printf("EOF")
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
