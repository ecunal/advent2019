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
	/*	grid := getGrid(stack)

		var robotLocation pair
		currentDirection := 0

		for i := range grid {
			for j := range grid[i] {
				if grid[i][j] == '^' {
					robotLocation = pair{i, j}
					break
				}
			}
		}

		var path []byte
		currentCount := 0

		for {
			next := pair{
				x: robotLocation.x + directions[currentDirection].x,
				y: robotLocation.y + directions[currentDirection].y,
			}

			if isScaffold(grid, next) {
				currentCount++
				robotLocation = next
				continue
			}

			if currentCount != 0 {
				path = append(path, []byte(strconv.Itoa(currentCount))...)
				path = append(path, ',')
				currentCount = 0
			}

			// check right
			right := directions[(currentDirection+1)%len(directions)]
			next = pair{
				x: robotLocation.x + right.x,
				y: robotLocation.y + right.y,
			}

			if isScaffold(grid, next) {
				path = append(path, 'R', ',')
				currentDirection = (currentDirection + 1) % len(directions)
				continue
			}

			// check left
			leftDirection := currentDirection - 1
			if leftDirection < 0 {
				leftDirection += len(directions)
			}
			left := directions[leftDirection]
			next = pair{
				x: robotLocation.x + left.x,
				y: robotLocation.y + left.y,
			}

			if isScaffold(grid, next) {
				path = append(path, 'L', ',')
				currentDirection = leftDirection
				continue
			}

			break
		}

		fmt.Println(string(path))
	*/
	movements := "A,B,A,C,A,A,C,B,C,B\nL,12,L,8,R,12\nL,10,L,8,L,12,R,12\nR,12,L,8,L,10\nn\n"

	stack[0] = 2 // part2

	pInput := make(chan int, len(movements))
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

	for _, bt := range []byte(movements) {
		pInput <- int(bt)
	}

	for {
		select {
		case o := <-pOutput:
			fmt.Print(o)
		case <-close:
			os.Exit(0)
		}
	}

}

func isScaffold(grid [][]byte, next pair) bool {
	return next.x >= 0 && next.x < len(grid) && next.y >= 0 && next.y < len(grid[next.x]) && grid[next.x][next.y] == '#'
}

type pair struct {
	x, y int
}

var directions = []pair{
	{-1, 0}, // north
	{0, 1},  // east
	{1, 0},  // south
	{0, -1}, // west
}

func getGrid(stack map[int]int) [][]byte {

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

	grid := [][]byte{[]byte{}}

	for {
		select {
		case o := <-pOutput:
			if o == '\n' {
				grid = append(grid, []byte{})
			} else {
				grid[len(grid)-1] = append(grid[len(grid)-1], byte(o))
			}
		case <-close:
			/* part1
			total := 0

			for i, row := range grid {
				if i == 0 || i >= len(grid)-3 {
					continue
				}

				for j, tile := range grid[i] {
					if tile != '#' || j == 0 || j == len(row)-1 {
						continue
					}
					if isIntersection(i, j, grid) {
						total += i * j
					}
				}
			}
			log.Println(total)
			*/
			return grid
		}
	}
}

/* part1
func isIntersection(i, j int, grid [][]byte) bool {
	for _, d := range directions {
		if grid[i+d.x][j+d.y] != '#' {
			return false
		}
	}
	return true
}
*/

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
