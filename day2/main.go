package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}

	inputstrings := strings.Split(strings.TrimSpace(string(input)), ",")
	original := make([]int, len(inputstrings))

	for i := range inputstrings {
		original[i], err = strconv.Atoi(inputstrings[i])
		if err != nil {
			log.Fatal("Error converting str to int:", err)
		}
	}

	for noun := 0; noun < 100; noun++ {

		for verb := 0; verb < 100; verb++ {

			inputs := make([]int, len(original))
			copy(inputs, original)

			inputs[1] = noun
			inputs[2] = verb

			for p := 0; inputs[p] != 99; p += 4 {

				if inputs[p] == 1 {
					// add

					inputs[inputs[p+3]] = inputs[inputs[p+1]] + inputs[inputs[p+2]]
				} else if inputs[p] == 2 {
					// multiply

					inputs[inputs[p+3]] = inputs[inputs[p+1]] * inputs[inputs[p+2]]
				} else {
					log.Print("error:", inputs[p])
				}
			}

			if inputs[0] == 19690720 {

				log.Println("Noun:", noun)
				log.Println("Verb:", verb)
				log.Println("Result:", 100*noun+verb)

				os.Exit(0)
			}

		}

	}

}
