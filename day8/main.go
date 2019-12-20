package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	input, err := ioutil.ReadFile("input")
	if err != nil {
		log.Fatal(err)
	}

	cols, rows := 25, 6

	lastImage := make([]byte, cols*rows)

	for j := 0; j < len(input); j += cols * rows {

		if input[j] == '\n' {
			break
		}

		for i := j; i < j+cols*rows; i++ {
			if lastImage[i-j] == 0 || lastImage[i-j] == '2' {
				lastImage[i-j] = input[i]
			}
		}
	}

	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			px := string(lastImage[cols*j+i])
			if px == "0" {
				fmt.Print(" ")
			} else {
				fmt.Print(px)
			}
		}
		fmt.Println()
	}

}
