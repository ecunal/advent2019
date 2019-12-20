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
	astreoids := parseInput()

	maxSeen := 0
	maxi, maxj := 0, 0

	for i := range astreoids {
		for j, cell := range astreoids[i] {
			if cell == '#' {
				seen := calculate(astreoids, pair{i, j})
				if seen > maxSeen {
					maxSeen = seen
					maxi, maxj = i, j
				}
			}
		}
	}

	log.Printf("Max: %d, i,j: %d,%d", maxSeen, maxj, maxi)

	result := vaporize(astreoids, pair{maxi, maxj})
	log.Println(result)
}

type pair struct {
	i, j int
}

func vaporize(astreoids [][]byte, center pair) int {
	angles := make(map[float64][]pair)
	for i := range astreoids {
		for j, cell := range astreoids[i] {
			if cell == '#' && !(i == center.i && j == center.j) {
				angle := math.Atan2(float64(i-center.i), float64(j-center.j))*180/math.Pi + 90
				if angle < 0 {
					angle += 360
				}
				if _, ok := angles[angle]; !ok {
					angles[angle] = make([]pair, 0)
				}
				angles[angle] = append(angles[angle], pair{i, j})
			}
		}
	}

	sortedAngles := make([]float64, 0, len(angles))
	for k := range angles {
		sortedAngles = append(sortedAngles, k)
	}
	sort.Float64s(sortedAngles)

	wantedlist := angles[sortedAngles[199]]
	log.Println(angles[sortedAngles[0]])
	log.Println(angles[sortedAngles[1]])
	log.Println(angles[sortedAngles[2]])
	log.Println(angles[sortedAngles[9]])
	log.Println(angles[sortedAngles[19]])
	log.Println(wantedlist)

	wanted := wantedlist[0]
	currDistance := math.MaxInt64
	for i, p := range wantedlist {
		if distance(p, center) < currDistance {
			wanted = wantedlist[i]
		}
	}

	// j * 100 + x
	return wanted.j*100 + wanted.i
}

func distance(p1, p2 pair) int {
	return int(math.Abs(math.Sqrt(math.Pow(float64(p1.i-p2.i), 2) + math.Pow(float64(p1.j-p2.j), 2))))
}

func calculate(astreoids [][]byte, curr pair) int {

	angles := make(map[float64][]pair)

	for i := range astreoids {
		for j, cell := range astreoids[i] {
			if cell == '#' && !(i == curr.i && j == curr.j) {
				angle := math.Atan2(float64(j-curr.j), float64(i-curr.i))
				if _, ok := angles[angle]; !ok {
					angles[angle] = make([]pair, 0)
				}
				angles[angle] = append(angles[angle], pair{i, j})
			}
		}
	}

	/*
		sortedAngles := make([]float64, 0, len(angles))
		for k := range angles {
			sortedAngles = append(sortedAngles, k)
		}
		sort.Float64s(sortedAngles)

	*/

	return len(angles)
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
		trimmed := strings.TrimSpace(input)
		if trimmed != "" {
			inputList = append(inputList, []byte(trimmed))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}
