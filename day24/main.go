package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	inputGrid := parseInput()

	midi, midj := len(inputGrid)/2, len(inputGrid[0])/2

	levels := make(map[int][][]byte)
	levels[0] = inputGrid
	levels[-1] = initializeEmptyLevel(len(inputGrid))
	levels[1] = initializeEmptyLevel(len(inputGrid))

	for m := 0; m < 200; m++ {
		newLevels := make(map[int][][]byte)

		for level, grid := range levels {
			newGrid := make([][]byte, len(grid))
			for i := range newGrid {
				newGrid[i] = make([]byte, len(grid[i]))
			}

			for i, row := range grid {
				for j, tile := range row {

					if i == midi && j == midj {
						// skip
						continue
					}

					bugCount := 0

					for neighbourDirection, d := range directions {
						ni := i + d.i
						nj := j + d.j

						if ni >= 0 && nj >= 0 && ni < len(grid) && nj < len(row) {

							if ni == midi && nj == midj { // go one level deeper
								if levelDown, ok := levels[level-1]; !ok {
									newLevels[level-1] = initializeEmptyLevel(len(grid))
								} else {
									bugCount += countLevelDownBugs(levelDown, neighbourDirection)
								}
							} else if grid[ni][nj] == '#' { // inside current grid
								bugCount++
							}

						} else { // go one level upper

							if levelUp, ok := levels[level+1]; !ok {
								newLevels[level+1] = initializeEmptyLevel(len(grid))
							} else {
								bugCount += countLevelUpBugs(levelUp, neighbourDirection)
							}
						}
					}
					newGrid[i][j] = game(tile, bugCount)
				}
			}

			newLevels[level] = newGrid
		}

		levels = newLevels
	}

	log.Println("Total bug count:", countBugs(levels))
}

func countLevelDownBugs(level [][]byte, d direction) int {
	var iidx, jidx []int

	switch d {
	case north:
		iidx = []int{len(level) - 1}
		jidx = make([]int, len(level[0]))
		for i := range jidx {
			jidx[i] = i
		}
	case south:
		iidx = []int{0}
		jidx = make([]int, len(level[0]))
		for i := range jidx {
			jidx[i] = i
		}
	case west:
		jidx = []int{len(level[0]) - 1}
		iidx = make([]int, len(level))
		for i := range iidx {
			iidx[i] = i
		}
	case east:
		jidx = []int{0}
		iidx = make([]int, len(level))
		for i := range iidx {
			iidx[i] = i
		}
	}

	count := 0
	for _, i := range iidx {
		for _, j := range jidx {
			if level[i][j] == '#' {
				count++
			}
		}
	}
	return count
}

func countLevelUpBugs(level [][]byte, d direction) int {
	midi, midj := len(level)/2, len(level[0])/2

	var i, j int
	switch d {
	case north:
		j = midj
		i = midi - 1
	case south:
		j = midj
		i = midi + 1
	case west:
		i = midi
		j = midj - 1
	case east:
		i = midi
		j = midj + 1
	}
	if level[i][j] == '#' {
		return 1
	}
	return 0
}

func countBugs(levels map[int][][]byte) int {
	count := 0
	for _, grid := range levels {
		for _, row := range grid {
			for _, tile := range row {
				if tile == '#' {
					count++
				}
			}
		}
	}
	return count
}

func game(tile byte, bugCount int) byte {
	if tile == '.' && (bugCount == 1 || bugCount == 2) {
		return '#'
	}
	if tile == '#' && bugCount == 1 {
		return '#'
	}
	return '.'
}

func initializeEmptyLevel(l int) [][]byte {
	emptyLevel := make([][]byte, l)
	for i := range emptyLevel {
		emptyLevel[i] = make([]byte, l)
		for j := range emptyLevel[i] {
			emptyLevel[i][j] = '.'
		}
	}
	return emptyLevel
}

func part1() {
	grid := parseInput()

	prevGrids := make(map[string]bool)

	for {
		currRep := serialize(grid)
		if prevGrids[currRep] {
			score := 0
			for i, c := range currRep {
				if c == '#' {
					score += 1 << i
				}
			}
			log.Println("Biodiversity rating:", score)
			return
		}
		prevGrids[currRep] = true

		newGrid := make([][]byte, len(grid))
		for i := range newGrid {
			newGrid[i] = make([]byte, len(grid[i]))
		}

		for i, row := range grid {
			for j, tile := range row {
				bugCount := 0

				for _, d := range directions {
					ni := i + d.i
					nj := j + d.j

					if ni >= 0 && nj >= 0 && ni < len(grid) && nj < len(row) && grid[ni][nj] == '#' {
						bugCount++
					}
				}

				newGrid[i][j] = game(tile, bugCount)
			}
		}
		grid = newGrid
	}
}

type pair struct {
	i, j int
}

type direction int

const (
	north direction = iota
	east
	south
	west
)

var directions = map[direction]pair{
	west:  {i: 0, j: -1},
	east:  {i: 0, j: 1},
	south: {i: 1, j: 0},
	north: {i: -1, j: 0},
}

func serialize(m [][]byte) string {
	r := ""
	for i := range m {
		r += string(m[i])
	}
	return r
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
		if input != "" {
			inputList = append(inputList, []byte(input))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}
