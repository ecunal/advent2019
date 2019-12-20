package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

const ore = "ORE"

func main() {
	log.SetFlags(log.Lshortfile)
	resources := parseInput()
	onefuel := oreCount(1, resources)
	orecount := 1000000000000
	maxFuel := 0
	for i := orecount / onefuel; ; i += 1000 {
		reqOre := oreCount(i, resources)
		if reqOre > orecount {
			for j := i - 1000; j < i; j++ {
				reqOre := oreCount(j, resources)
				if reqOre > orecount {
					break
				}
				if j > maxFuel {
					maxFuel = j
				}
			}
			break
		}
		if i > maxFuel {
			maxFuel = i
		}
	}
	log.Println(maxFuel)
}

func oreCount(fuel int, resources map[string]resource) int {
	leftoverResources := make(map[string]int)
	currentResources := make(map[string]int)

	currentResources["FUEL"] = fuel

	for {
		if _, ok := currentResources[ore]; ok && len(currentResources) == 1 {
			break
		}

		tmp := make(map[string]int)

		for k, v := range currentResources {

			if k == ore {
				tmp[ore] += v
				continue
			}

			reqAmount := v
			if lf, ok := leftoverResources[k]; ok && lf > 0 {
				if lf >= reqAmount {
					leftoverResources[k] = lf - reqAmount
					continue
				}
				leftoverResources[k] = 0
				reqAmount -= lf
			}

			times := reqAmount / resources[k].amount
			if reqAmount%resources[k].amount != 0 {
				times++
			}

			for _, r := range resources[k].req {
				tmp[r.name] += r.amount * times
			}

			if leftover := resources[k].amount*times - reqAmount; leftover > 0 {
				leftoverResources[k] += leftover
			}
		}

		currentResources = tmp
	}

	return currentResources[ore]
}

type resource struct {
	name   string
	amount int
	req    []resource
}

func parseInput() map[string]resource {
	f, err := os.Open("input")
	if err != nil {
		log.Fatal("cannot open input:", err)
	}
	defer f.Close()
	inputList := make(map[string]resource)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := scanner.Text()
		trimmed := strings.TrimSpace(input)
		if trimmed == "" {
			continue
		}
		sides := strings.Split(trimmed, "=>")
		req := strings.Split(sides[0], ",")
		reqResources := make([]resource, len(req))
		for i, line := range req {
			reqResources[i] = parseResource(line)
		}
		resultResource := parseResource(sides[1])
		resultResource.req = reqResources
		inputList[resultResource.name] = resultResource
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning:", err)
	}
	return inputList
}

func parseResource(line string) resource {
	result := strings.Split(strings.TrimSpace(line), " ")
	amount, err := strconv.Atoi(result[0])
	if err != nil {
		log.Fatal(err)
	}
	return resource{
		name:   result[1],
		amount: amount,
	}
}
