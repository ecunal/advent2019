package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

func main() {
	grid := parseInput()

	var curr [4]pair
	curridx := 0
	gates := make(map[byte]pair)
	noKeys := 0
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] >= 'A' && grid[i][j] <= 'Z' {
				gates[grid[i][j]] = pair{i, j}
			} else if isKey(grid[i][j]) {
				noKeys++
			} else if grid[i][j] == '@' {
				curr[curridx] = pair{i, j}
				curridx++
			}
		}
	}

	totalSteps := recursive(noKeys, []byte{}, curr, grid, gates, make(map[cacheKey]int))
	log.Println(totalSteps)
}

type cacheKey struct {
	curr      [4]pair
	foundKeys string
}

func hash(foundKeys []byte) string {
	sort.Slice(foundKeys, func(i, j int) bool {
		return foundKeys[i] < foundKeys[j]
	})
	return string(foundKeys)
}

func recursive(noKeys int, foundKeys []byte, currs [4]pair, grid [][]byte, gates map[byte]pair, cache map[cacheKey]int) int {
	if len(foundKeys) == noKeys {
		return 0
	}

	hashed := hash(foundKeys)
	//	log.Println("Hashed found keys:", hashed)

	if steps, ok := cache[cacheKey{curr: currs, foundKeys: hashed}]; ok {
		//		log.Println("Curr:", curr, "foundKeys:", hashed, "returning:", steps)
		return steps
	}

	minSteps := math.MaxInt64

	for i, curr := range currs {

		possibleKeys := bfs(curr, grid)

		//	log.Println("Possible keys:", possibleKeys)

		for _, pk := range possibleKeys {
			key := grid[pk.i][pk.j]
			gate := gates[key-32]
			grid[gate.i][gate.j] = '.'
			grid[pk.i][pk.j] = '@'
			grid[curr.i][curr.j] = '.'

			copyKeys := make([]byte, len(foundKeys)+1)
			copy(copyKeys, foundKeys)
			copyKeys[len(copyKeys)-1] = key

			var copycurr [4]pair
			for j := range currs {
				copycurr[j] = currs[j]
			}
			copycurr[i] = pair{pk.i, pk.j}

			//		log.Println("copykeys:", copyKeys)

			steps := recursive(noKeys, copyKeys, copycurr, grid, gates, cache)
			if steps+pk.steps < minSteps {
				minSteps = steps + pk.steps
			}

			// rollback
			grid[gate.i][gate.j] = key - 32
			grid[pk.i][pk.j] = key
			grid[curr.i][curr.j] = '@'
		}
	}

	cache[cacheKey{curr: currs, foundKeys: hashed}] = minSteps
	//	log.Println("Curr:", curr, "foundKeys:", hashed, "adding:", minSteps)
	return minSteps
}

func bfs(curr pair, grid [][]byte) []pairSteps {
	possibleKeys := []pairSteps{}
	queue := []pairSteps{{i: curr.i, j: curr.j, steps: 0}}

	visited := make(map[pair]bool)

	for len(queue) > 0 {
		top := queue[0]
		queue = queue[1:]

		toppair := pair{i: top.i, j: top.j}
		if visited[toppair] {
			continue
		}

		visited[toppair] = true

		for _, d := range directions {
			next := pairSteps{
				i:     top.i + d.i,
				j:     top.j + d.j,
				steps: top.steps + 1,
			}

			if next.i >= 0 && next.i < len(grid) && next.j >= 0 && next.j < len(grid[next.i]) {
				if isKey(grid[next.i][next.j]) {
					possibleKeys = append(possibleKeys, next)
				} else if grid[next.i][next.j] == '.' {
					queue = append(queue, next)
				}
			}
		}
	}
	return possibleKeys
}

func isKey(a byte) bool {
	return a >= 'a' && a <= 'z'
}

type pairSteps struct {
	i, j  int
	steps int
}

type pair struct {
	i, j int
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
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		inputList = append(inputList, []byte(input))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}
