package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
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

	arcade := &game{
		grid: make(map[pair]int),
		mu:   &sync.Mutex{},
	}

	stack[0] = 2 // part2
	arcade.play(stack)
}

type game struct {
	grid         map[pair]int
	ball, paddle pair
	score        int
	mu           *sync.Mutex
	rounds       int
}

func (g *game) getInput() int {
	g.mu.Lock()
	defer g.mu.Unlock()

	/*time.Sleep(10 * time.Millisecond)
	fmt.Println("\x1B[2J\x1B[H")
	printBoard(g.grid)
	time.Sleep(10 * time.Millisecond)*/
	g.rounds++

	if g.paddle.x < g.ball.x {
		return 1
	}
	if g.paddle.x > g.ball.x {
		return -1
	}
	return 0
}

func (g *game) play(stack map[int]int) {
	pOutput := make(chan int)
	close := make(chan bool)
	prog := &program{
		relativeBase: 0,
		input:        g.getInput,
		output:       pOutput,
		close:        close,
		program:      stack,
	}
	go prog.run()

	for {
		select {
		case x := <-pOutput:
			g.mu.Lock()
			y := <-pOutput
			tileID := <-pOutput

			if x == -1 && y == 0 {
				g.score = tileID
			} else {
				curr := pair{x, y}
				g.grid[curr] = tileID

				if tileID == 4 {
					g.ball = curr
				} else if tileID == 3 {
					g.paddle = curr
				}
			}
			g.mu.Unlock()

			//log.Printf("Ball: %v, paddle: %v, score: %d", g.ball, g.paddle, g.score)
			// fmt.Print("\x1B[2J\x1B[H")
			// printBoard(g.grid)

		case <-close:
			printBoard(g.grid)
			log.Println("Score:", g.score)
			log.Println("Rounds:", g.rounds)
			return
		}
	}
}

func printBoard(grid map[pair]int) {
	sorted := make([]pair, 0, len(grid))
	for k := range grid {
		sorted = append(sorted, k)
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].y == sorted[j].y {
			return sorted[i].x < sorted[j].x
		}
		return sorted[i].y < sorted[j].y
	})

	for i, k := range sorted {
		if i != 0 && sorted[i-1].y != sorted[i].y {
			fmt.Println()
		}
		switch grid[k] {
		case 0:
			fmt.Print(" ")
		case 1:
			fmt.Print("|")
		case 2:
			fmt.Print("#")
		case 3:
			fmt.Print("-")
		case 4:
			fmt.Print("*")
		}
	}

	fmt.Println()
}

type pair struct {
	x, y int
}

type program struct {
	relativeBase int
	input        func() int
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
			p.program[o] = p.input()
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
