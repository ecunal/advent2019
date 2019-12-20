package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type coord struct {
	x, y int
}

type value struct {
	visited bool
	steps   int
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return 0 - x
}

func main() {
	wires := parseInput()

	grid := make(map[coord]value)

	shortestDistance := math.MaxInt64

	for i, wire := range wires {
		x, y, steps := 0, 0, 0

		for _, path := range wire {

			newx, newy := move(path, x, y)

			// log.Printf("Wire %d: %s. Current x,y: %d, %d. New x,y: %d, %d", i, path, x, y, newx, newy)

			if x != newx {
				from, to, inc := x, newx, 1
				if x > newx {
					inc = -1
				}

				for from = from + inc; from != to+inc; from += inc {
					c := coord{from, newy}
					v, ok := grid[c]

					steps++

					// log.Printf("Visiting %d, %d, current steps %d", c.x, c.y, steps)

					if i == 1 && ok && v.visited {
						// intersection!

						log.Printf("Intersection found on %d, %d. steps: %d, %d", c.x, c.y, v.steps, steps)

						distance := v.steps + steps
						if distance < shortestDistance && c.x != 0 && c.y != 0 {
							shortestDistance = distance
						}
					}

					if i == 0 {
						grid[c] = value{
							visited: true,
							steps:   steps,
						}
					}
				}

			} else {
				from, to, inc := y, newy, 1
				if y > newy {
					inc = -1
				}

				for from = from + inc; from != to+inc; from += inc {
					c := coord{newx, from}
					v, ok := grid[c]

					steps++

					//	log.Printf("Visiting %d, %d, current steps %d", c.x, c.y, steps)

					if i == 1 && ok && v.visited {
						// intersection!

						log.Printf("Intersection found on %d, %d. steps: %d, %d", c.x, c.y, v.steps, steps)

						distance := v.steps + steps
						if distance < shortestDistance && c.x != 0 && c.y != 0 {
							shortestDistance = distance
						}
					}
					if i == 0 {
						grid[c] = value{
							visited: true,
							steps:   steps,
						}
					}
				}
			}

			x, y = newx, newy
		}
	}

	log.Println(shortestDistance)
}

func move(path string, i, j int) (int, int) {
	count, err := strconv.Atoi(path[1:])
	if err != nil {
		log.Fatal(err)
	}

	direction := path[0]

	switch direction {
	case 'R':
		return i, j + count
	case 'D':
		return i - count, j
	case 'U':
		return i + count, j
	case 'L':
		return i, j - count
	}

	log.Fatal("unknown direction:", direction)
	return 0, 0
}

func parseInput() [][]string {
	f, err := os.Open("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}
	defer f.Close()

	inputList := make([][]string, 0)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			break
		}

		inputList = append(inputList, strings.Split(input, ","))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}
