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

	computers := make([]*program, 50)
	inputChans := make([]chan int, 50)
	outputChans := make([]chan int, 50)

	for i := range computers {
		copystack := make(map[int]int)
		for k, v := range stack {
			copystack[k] = v
		}

		inputChans[i] = make(chan int, 100)
		outputChans[i] = make(chan int)
		computers[i] = &program{
			relativeBase: 0,
			input:        inputChans[i],
			output:       outputChans[i],
			close:        make(chan bool),
			program:      copystack,
		}
		inputChans[i] <- i // give their ids
		go computers[i].run()
	}

	natX, natY, lastY := 0, 0, 0
	count := 0
	for {
		for i := range outputChans {
			select {
			case addr := <-outputChans[i]:
				x := <-outputChans[i]
				y := <-outputChans[i]

				if addr == 255 {
					natX, natY = x, y
				} else {
					inputChans[addr] <- x
					inputChans[addr] <- y
				}
				count = 0
			default:
				count++
			}
		}

		if count > 1000000 && natX != 0 && natY != 0 {
			inputChans[0] <- natX
			inputChans[0] <- natY

			if natY == lastY {
				log.Println(natY, lastY)
				return
			}
			lastY = natY
			natX = 0
			natY = 0
			count = 0
		}
	}
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
			var inp int
			select {
			case inp = <-p.input:
			default:
				inp = -1
			}
			p.program[o] = inp
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
