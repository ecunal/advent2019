package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)

	input := parseInput()
	portals := make(map[string][]pair)
	maze := make([][]string, len(input))
	for i := range maze {
		maze[i] = make([]string, len(input[i]))
	}

	for i := 0; i < len(input)-1; i++ {
		for j := 0; j < len(input[i])-1; j++ {
			if isLetter(input[i][j]) {

				// vertical
				if isLetter(input[i+1][j]) {
					gate := string(input[i][j]) + string(input[i+1][j])
					log.Println(gate, i, j)

					if _, ok := portals[gate]; !ok {
						portals[gate] = make([]pair, 0)
					}

					if i+2 < len(input) && input[i+2][j] == '.' {
						maze[i][j] = " "
						maze[i+1][j] = gate

						portals[gate] = append(portals[gate], pair{i + 1, j})
					} else if i-1 >= 0 && input[i-1][j] == '.' {
						maze[i+1][j] = " "
						maze[i][j] = gate

						portals[gate] = append(portals[gate], pair{i, j})
					}
				}

				// horizontal
				if isLetter(input[i][j+1]) {
					gate := string(input[i][j]) + string(input[i][j+1])

					log.Println(gate, i, j)

					if _, ok := portals[maze[i][j+1]]; !ok {
						portals[maze[i][j+1]] = make([]pair, 0)
					}

					if j+2 < len(input[i]) && input[i][j+2] == '.' {
						maze[i][j] = " "
						maze[i][j+1] = gate

						portals[gate] = append(portals[gate], pair{i, j + 1})
					} else if j-1 >= 0 && input[i][j-1] == '.' {
						maze[i][j+1] = " "
						maze[i][j] = gate

						portals[gate] = append(portals[gate], pair{i, j})
					}
				}

			} else {
				maze[i][j] = string(input[i][j])
			}
		}
	}

	log.Println(portals)

	var begin pair
	for _, d := range directions {
		n := pair{
			i: portals["AA"][0].i + d.i,
			j: portals["AA"][0].j + d.j,
		}
		if maze[n.i][n.j] == "." {
			begin = n
			break
		}
	}

	shortest := bfs(begin, maze, portals)
	log.Println("Shortest path:", shortest)
}

func bfs(curr pair, grid [][]string, portals map[string][]pair) int {
	queue := []pairSteps{{p: curr, steps: 0}}
	visited := make(map[pair]bool)

	for len(queue) > 0 {
		top := queue[0]
		queue = queue[1:]

		if visited[top.p] {
			continue
		}
		visited[top.p] = true

		for _, d := range directions {
			next := pairSteps{
				p: pair{
					i: top.p.i + d.i,
					j: top.p.j + d.j,
				},
				steps: top.steps + 1,
			}

			if next.p.i >= 0 && next.p.i < len(grid) && next.p.j >= 0 && next.p.j < len(grid[next.p.i]) {
				nextTile := grid[next.p.i][next.p.j]

				if nextTile == "ZZ" {
					return top.steps
				}

				if len(nextTile) == 2 && nextTile != "AA" { // portal
					other := getOtherLocation(portals[nextTile], next.p)
					log.Println("taking portal from", next, "to", other)
					queue = append(queue, pairSteps{p: other, steps: top.steps})
				} else if nextTile == "." {
					queue = append(queue, next)
				}
			}
		}
	}
	return 0
}

func getOtherLocation(locations []pair, curr pair) pair {
	for _, p := range locations {
		if p != curr {
			return p
		}
	}
	return curr
}

type pair struct {
	i, j int
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
