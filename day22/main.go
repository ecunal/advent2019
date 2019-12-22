package main

import (
	"bufio"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func main() {
	input := parseInput()

	deckSize := big.NewInt(119315717514047)
	times := big.NewInt(101741582076661)

	inc, offset := big.NewInt(1), big.NewInt(0)

	for _, pi := range input {
		switch pi.t {
		case tNew:
			inc = inc.Mul(inc, big.NewInt(-1))
			inc = inc.Mod(inc, deckSize)
			offset = offset.Add(offset, inc)
			offset = offset.Mod(offset, deckSize)
		case tCut:
			arg := big.NewInt(pi.arg)
			arg = arg.Mul(arg, inc)
			offset = offset.Add(offset, arg)
			offset = offset.Mod(offset, deckSize)
		case tInc:
			arg := big.NewInt(pi.arg)
			arg = arg.ModInverse(arg, deckSize)
			inc = inc.Mul(inc, arg)
			inc = inc.Mod(inc, deckSize)
		}
	}

	incIt := &big.Int{}
	incIt = incIt.Set(inc).Exp(inc, times, deckSize)

	minusinc := big.NewInt(1)
	minusinc = minusinc.Sub(minusinc, incIt)

	modinc := big.NewInt(1)
	modinc = modinc.Sub(modinc, inc)
	modinc = modinc.Mod(modinc, deckSize)
	modinc = modinc.ModInverse(modinc, deckSize)

	offsetIt := &big.Int{}
	offsetIt = offsetIt.Mul(offset, minusinc)
	offsetIt = offsetIt.Mul(offsetIt, modinc)
	offsetIt = offsetIt.Mod(offsetIt, deckSize)

	pos := big.NewInt(2020)
	pos = pos.Mul(pos, incIt)
	pos = pos.Add(pos, offsetIt)
	pos = pos.Mod(pos, deckSize)

	log.Println(pos.String())
}

func part1() {
	const deckSize = 10007
	input := parseInput()

	deck := make([]int, deckSize)
	for i := range deck {
		deck[i] = i
	}

	for _, pi := range input {
		switch pi.t {
		case tNew:

			for i := len(deck)/2 - 1; i >= 0; i-- {
				opp := len(deck) - 1 - i
				deck[i], deck[opp] = deck[opp], deck[i]
			}

		case tCut:

			N := pi.arg
			if N < 0 {
				N = int64(len(deck)) + N
			}
			deck = append(deck[N:], deck[:N]...)

		case tInc:

			n := pi.arg
			res := make([]int, len(deck))
			ridx := int64(0)
			for i := range deck {
				res[ridx] = deck[i]
				ridx = (ridx + n) % int64(len(deck))
			}
			deck = res

		}
	}

	for i, n := range deck {
		if n == 2009 {
			log.Println(i)
			return
		}
	}
}

type technique int

const (
	tNew technique = iota
	tCut
	tInc
)

type puzzleInput struct {
	t   technique
	arg int64
}

func parseInput() []puzzleInput {
	f, err := os.Open("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}
	defer f.Close()
	inputList := make([]puzzleInput, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if input == "deal into new stack" {
			inputList = append(inputList, puzzleInput{t: tNew})
		} else if strings.HasPrefix(input, "cut") {
			n, err := strconv.ParseInt(strings.TrimPrefix(input, "cut "), 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			inputList = append(inputList, puzzleInput{t: tCut, arg: n})
		} else if strings.HasPrefix(input, "deal with increment") {
			n, err := strconv.ParseInt(strings.TrimPrefix(input, "deal with increment "), 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			inputList = append(inputList, puzzleInput{t: tInc, arg: n})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}
