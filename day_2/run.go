package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type countLimit struct {
	blue int
	red int
	green int
}

func (c countLimit) possible(count []int) bool {
	return (count[0] <= c.blue ) && (count[1] <= c.red ) && ( count[2] <= c.green)
}

func readLines(filename string) []string {

	data, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines

}

func extractId(line string) int {
	idRegex := regexp.MustCompile(`(?P<id>\d+):`)
	match := idRegex.FindStringSubmatch(line)
	
	if match == nil {
		log.Fatal("Could not parse ID from line: " + line)
	}
	idStr := match[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	return id

}


func splitLine(line string) []string {
	var splits []string
	var split string

	reachedColon := false
	for _, c := range line {
		if !reachedColon {
			reachedColon = c == ':'
			continue
		}
		if c == ';' {
			splits = append(splits, split)
			split = ""
			continue
		}
		split += string(c)
	}
	splits = append(splits, split)
	return splits
}

func parseSplit(split string) []int {
	patterns := []*regexp.Regexp{regexp.MustCompile(`(?P<blue>\d+) blue`),
		regexp.MustCompile(`(?P<red>\d+) red`),
		regexp.MustCompile(`(?P<green>\d+) green`)}

	
	var counts []int
	for _, re := range patterns {
		match := re.FindStringSubmatch(split)
		if match == nil {
			counts = append(counts, 0)
			continue
		}
		str := match[1]
		count, err := strconv.Atoi(str)
		if err != nil {
			log.Fatal(err)
		}
		counts = append(counts, count)
	} 
	return counts
}

func largestCount(counts [][]int) []int {
	largest := []int{0,0,0}
	for _, count := range counts {
		for i := 0; i < 3; i++ {
			if largest[i] < count[i] {
				largest[i] = count[i]
			}
		}
	}
	return largest
}

func power(counts []int) int {
	return counts[0] * counts[1] * counts[2]
}

func main() {
	// Part 1
	limit := countLimit{
		red: 12,
		green: 13,
		blue: 14,
	}

	total := 0
	for _, line := range readLines("input.txt"){
		var counts [][]int
		for _, split := range splitLine(line) {
			counts = append(counts, parseSplit(split))
		}
		largest := largestCount(counts)
		if limit.possible(largest){
			total += extractId(line)
		}
	}
	fmt.Println(total)

	// Part 2
	total = 0
	for _, line := range readLines("input.txt"){
		var counts [][]int
		for _, split := range splitLine(line) {
			counts = append(counts, parseSplit(split))
		}
		largest := largestCount(counts)
		total += power(largest)
	}
	fmt.Println(total)
}
