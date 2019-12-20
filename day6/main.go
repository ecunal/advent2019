package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const parent = "COM"

func main() {
	inputs := parseInput()

	graph := make(map[string][]string)

	for _, p := range inputs {
		if v, ok := graph[p.parent]; !ok || v == nil {
			graph[p.parent] = make([]string, 0)
		}
		graph[p.parent] = append(graph[p.parent], p.child)
	}

	you := parents("YOU", graph)
	santa := parents("SAN", graph)

	firstCommonParent := ""
	total := 0

	for i := 0; i < len(you); i++ {
		if i >= len(santa) {
			break
		}
		if you[i] != santa[i] {
			break
		}
		firstCommonParent = you[i]
		total = len(you) - i + len(santa) - i - 4
	}

	log.Println("Common parent:", firstCommonParent)
	log.Println("Total:", total)
}

func parents(child string, graph map[string][]string) []string {
	result := make([]string, 0)

	return dfs(parent, child, graph, result)
}

func dfs(current, goal string, graph map[string][]string, stack []string) []string {
	if current == goal {
		return stack
	}
	if v, ok := graph[current]; !ok || len(v) == 0 {
		return nil
	}

	for _, child := range graph[current] {
		stack = append(stack, child)

		result := dfs(child, goal, graph, stack)
		if result != nil {
			return result
		}

		stack = stack[:len(stack)-1]
	}

	return nil
}

func day1(graph map[string][]string) {
	queue := make([]string, 1)
	queue[0] = parent

	total, level := 0, 0

	for len(queue) != 0 {

		tmpqueue := make([]string, 0)

		for len(queue) != 0 {
			top := queue[0]

			for _, child := range graph[top] {
				tmpqueue = append(tmpqueue, child)
			}

			queue = queue[1:]
			total += level
		}

		queue = append(queue, tmpqueue...)

		level++
	}

	log.Println("Total:", total)
}

type pair struct {
	parent, child string
}

func parseInput() []pair {
	f, err := os.Open("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}
	defer f.Close()

	inputList := make([]pair, 0)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		input := scanner.Text()

		trimmed := strings.TrimSpace(input)
		if trimmed == "" {
			continue
		}

		split := strings.Split(trimmed, ")")
		inputList = append(inputList, pair{
			parent: split[0],
			child:  split[1],
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}
