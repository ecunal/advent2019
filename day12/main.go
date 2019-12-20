package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type triplet struct {
	x, y, z int
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func main() {
	moons := parseInput()
	returned := make(map[int]int)
	initials := [3][4]int{}
	for i := range initials {
		initials[i] = [4]int{}
		for j := range moons {
			initials[i][j] = moons[j].pos[i]
		}
	}

	currentTurn := 1

	for {
		if len(returned) == 3 {
			break
		}

		if currentTurn != 1 {
			for i := range initials {
				allEqual := true
				for j := range moons {
					if moons[j].pos[i] != initials[i][j] {
						allEqual = false
					}
				}
				if allEqual {
					log.Println("Found match for", i, moons, currentTurn)
					if _, ok := returned[i]; !ok {
						returned[i] = currentTurn
					}
				}
			}
		}

		for m, mymoon := range moons {
			for j := range moons {
				if m == j {
					continue
				}

				for p := range mymoon.pos {
					if mymoon.pos[p] > moons[j].pos[p] {
						moons[m].vel[p]--
					} else if mymoon.pos[p] < moons[j].pos[p] {
						moons[m].vel[p]++
					}
				}
			}
		}

		for m := range moons {
			for p := range moons[m].pos {
				moons[m].pos[p] += moons[m].vel[p]
			}
		}
		currentTurn++
	}

	log.Println(returned)
	log.Println(LCM(returned[0], returned[1], returned[2]))

	/* part 1
	totalEnergy := 0

	for m := range moons {
		kin, pot := 0, 0
		for p := range moons[m].pos {
			pot += abs(moons[m].pos[p])
			kin += abs(moons[m].vel[p])
		}
		totalEnergy += pot * kin
	}

	log.Println(totalEnergy)
	*/
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type moon struct {
	pos [3]int
	vel [3]int
}

func parseInput() []moon {
	f, err := os.Open("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}
	defer f.Close()
	inputList := make([]moon, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()
		trimmed := strings.TrimSpace(input)
		if trimmed == "" {
			continue
		}
		positions := strings.Split(trimmed[1:len(trimmed)-1], ",")
		pos := [3]int{0, 0, 0}
		for i, str := range positions {
			pos[i], err = strconv.Atoi(strings.Split(str, "=")[1])
			if err != nil {
				log.Fatal("Error with input:", positions[i], err)
			}
		}
		inputList = append(inputList, moon{
			pos: pos,
			vel: [3]int{0, 0, 0},
		})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}
