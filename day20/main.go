package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)

	input := parseInput()
	portals := make(map[string]portal)
	maze := make([][]string, len(input))
	for i := range maze {
		maze[i] = make([]string, len(input[i]))
	}

	for i := 0; i < len(input)-1; i++ {
		for j := 0; j < len(input[i])-1; j++ {
			if isLetter(input[i][j]) {

				maze[i][j] = " "

				// vertical
				if isLetter(input[i+1][j]) {
					gate := string(input[i][j]) + string(input[i+1][j])
					var pos pair

					if i+2 < len(input) && input[i+2][j] == '.' {
						maze[i+2][j] = gate
						pos = pair{i + 2, j, 0}
					} else if i-1 >= 0 && input[i-1][j] == '.' {
						maze[i-1][j] = gate
						pos = pair{i - 1, j, 0}
					}

					prt := portals[gate]
					if i == 0 || i == len(input)-2 {
						prt.outer = pos
					} else {
						prt.inner = pos
					}
					portals[gate] = prt
				}

				// horizontal
				if isLetter(input[i][j+1]) {
					gate := string(input[i][j]) + string(input[i][j+1])
					var pos pair

					if j+2 < len(input[i]) && input[i][j+2] == '.' {
						maze[i][j+2] = gate
						pos = pair{i, j + 2, 0}
					} else if j-1 >= 0 && input[i][j-1] == '.' {
						maze[i][j-1] = gate
						pos = pair{i, j - 1, 0}
					}

					prt := portals[gate]
					if j == 0 || j == len(input[i])-2 {
						prt.outer = pos
					} else {
						prt.inner = pos
					}
					portals[gate] = prt
				}

			} else if maze[i][j] == "" {
				maze[i][j] = string(input[i][j])
			}
		}
	}
	/*
		for _, row := range maze {
			for _, tile := range row {
				fmt.Print(tile)
			}
			fmt.Println()
		}

		log.Println(portals)
	*/
	begin := portals["AA"].outer

	// log.Println(begin)

	shortest := bfs(begin, maze, portals)
	log.Println("Shortest path:", shortest)
}

type portal struct {
	inner pair
	outer pair
}

func bfs(curr pair, grid [][]string, portals map[string]portal) int {
	queue := []pairSteps{{p: curr, steps: 0}}
	visited := make(map[pair]bool)

	for len(queue) > 0 {
		top := queue[0]
		queue = queue[1:]

		if visited[top.p] {
			continue
		}
		visited[top.p] = true

		//	log.Println("Current:", top)

		for _, d := range directions {
			next := pairSteps{
				p: pair{
					i:     top.p.i + d.i,
					j:     top.p.j + d.j,
					level: top.p.level,
				},
				steps: top.steps + 1,
			}

			if next.p.i >= 0 && next.p.i < len(grid) && next.p.j >= 0 && next.p.j < len(grid[next.p.i]) {
				nextTile := grid[next.p.i][next.p.j]

				if nextTile == "ZZ" && next.p.level == 0 {
					log.Println(next.steps)
					return next.steps
				}

				if len(nextTile) == 2 && nextTile != "AA" { // portal
					nextPortal := pairSteps{steps: next.steps + 1}
					locations := portals[nextTile]

					if next.p.i == locations.inner.i && next.p.j == locations.inner.j {
						// one level down
						nextPortal.p = locations.outer
						nextPortal.p.level = next.p.level + 1

						log.Println("taking downwards (+1) portal: ", nextTile, top, nextPortal)
						queue = append(queue, nextPortal)
					} else if next.p.i == locations.outer.i && next.p.j == locations.outer.j && next.p.level != 0 {
						// one level up
						nextPortal.p = locations.inner
						nextPortal.p.level = next.p.level - 1

						log.Println("taking upwards (-1) portal: ", nextTile, top, nextPortal)
						queue = append(queue, nextPortal)
					}
				} else if nextTile == "." {
					queue = append(queue, next)
				}
			}
		}
	}
	return 0
}

type pair struct {
	i, j  int
	level int
}

type pairSteps struct {
	p     pair
	steps int
}

func isLetter(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

var directions = []pair{
	{i: 0, j: -1},
	{i: 0, j: 1},
	{i: 1, j: 0},
	{i: -1, j: 0},
}

func parseInput() [][]byte {
	f, err := os.Open("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}
	defer f.Close()
	inputList := make([][]byte, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()
		inputList = append(inputList, []byte(input))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}
