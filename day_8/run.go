package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func readLines(filename string) (lines []string) {
	data, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return
}

func parseNode(line string, nodes map[string][2]string) map[string][2]string{
	r := regexp.MustCompile(`[A-Z0-9]+`)
	strings := r.FindAllString(line, -1)

	nodes[strings[0]] = [2]string{strings[1], strings[2]}
	return nodes
}

func followRoute(route string, graph map[string][2]string) int {
	i := 0
	node := "AAA"
	for {
		d := route[i % len(route)]
		switch d {
		case 'L':
			node = graph[node][0]
		case 'R':
			node = graph[node][1]
		}
		i += 1

		if node == "ZZZ" {
			return i	
		}
	}
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

func followAllRoutes(nodes []string, route string, graph map[string][2]string) int {
	i := 0
	var idx int

	distances := []int{}
	for i:=0; i < len(nodes); i++ {
		distances = append(distances, -1)
	}

	for {
		d := route[i % len(route)]
		switch d {
		case 'L':
			idx = 0
		case 'R':
			idx = 1
		}
		i += 1

		// Follow each path until the first Z
		foundEnd := true
		for j, n := range nodes {
			nodes[j] = graph[n][idx]	
			if nodes[j][2] == 'Z' && distances[j] == -1{
				distances[j] = i 
			}
			if distances[j] == -1 {
				foundEnd = false
			}
		}

		// find the LCM of each path
		if foundEnd { 
			lcm := distances[0]
			for _, v := range distances[1:]{
				lcm = LCM(v, lcm)
			}
			return lcm
		}
	}
}

func main() {
	lines := readLines("day_8/input.txt")

	route := lines[0]
	nodes := make(map[string][2]string)

	// part 1
	for _, l := range lines[2:] {
		nodes = parseNode(l, nodes)
	}
	steps := followRoute(route, nodes)
	fmt.Println(steps)

	// part 2
	starts := []string{}
	for key, _ := range nodes {
		if key[2] == 'A' {
			starts = append(starts, key)
		}
	}
	fmt.Println(starts)
	steps = followAllRoutes(starts, route, nodes)
	fmt.Println(steps)
	



}
