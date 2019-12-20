package main

import (
	"fmt"
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
	grid := make(map[pair]int)
	dfs(curr, grid, pInput, pOutput)
	printBoard(grid)

	oxygen := pair{16, 14}

	queue := []pairMin{{oxygen, 0}}
	mins := 0
	visited := make(map[pair]bool)

	for len(queue) != 0 {
		top := queue[0]
		queue = queue[1:]

		if visited[top.p] {
			continue
		}
		visited[top.p] = true

		for _, v := range directions {
			next := pair{top.p.x + v.x, top.p.y + v.y}

			if grid[next] != 0 {
				queue = append(queue, pairMin{next, top.minutes + 1})
			}
		}
		mins = top.minutes
	}

	fmt.Println(mins)
}

type pairMin struct {
	p       pair
	minutes int
}

func dfs(curr pair, grid map[pair]int, input, output chan int) {

	for i, v := range directions {
		next := pair{curr.x + v.x, curr.y + v.y}

		if _, ok := grid[next]; ok {
			continue
		}

		input <- i
		o := <-output
		grid[next] = o
		if o == 0 {
			continue
		}
		if o == 2 {
			log.Println("found!", next)
		}

		dfs(next, grid, input, output)

		if i%2 == 1 {
			input <- i + 1
		} else {
			input <- i - 1
		}
		<-output
	}
}

func printBoard(grid map[pair]int) {

	min, max := pair{math.MaxInt64, math.MaxInt64}, pair{math.MinInt64, math.MinInt64}
	for k := range grid {
		if k.x < min.x {
			min.x = k.x
		}
		if k.x > max.x {
			max.x = k.x
		}
		if k.y < min.y {
			min.y = k.y
		}
		if k.y > max.y {
			max.y = k.y
		}
	}

	fmt.Println(min, max)

	for i := min.x; i <= max.x; i++ {
		for j := min.y; j <= max.y; j++ {
			if v, ok := grid[pair{i, j}]; !ok {
				fmt.Print(" ")
			} else if v == 0 {
				fmt.Print("#")
			} else if v == 1 {
				fmt.Print(".")
			} else if v == 2 {
				fmt.Print("O")
			}
		}
		fmt.Println()
	}

	fmt.Println()
}

type pair struct {
	x, y int
}

var directions = map[int]pair{
	1: {0, -1}, // north
	2: {0, 1},  // south
	3: {1, 0},  // west
	4: {-1, 0}, // east
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
