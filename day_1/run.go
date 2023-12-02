package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
)

func readLines(filePath string) []string {
	// Return a scanner that reads from a file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var output []string

	for scanner.Scan() {
		output = append(output, scanner.Text())
	}
	return output
}

func toInt(value rune) int {
	return int(value) - int('0')
}

func firstLastDigits(line string) (int, int) {
	first := 'x'
	last := 'x'

	for _, char := range line {
		if char >= '0' && char <= '9' {
			if first == 'x' {
				first = char
			}
			last = char
		}
	}
	return toInt(first), toInt(last)
}

func extractNumbers(line string) ([]int, error) {
	regexes := [10]*regexp.Regexp{
		regexp.MustCompile(`zero`),
		regexp.MustCompile(`one`),
		regexp.MustCompile(`two`),
		regexp.MustCompile(`three`),
		regexp.MustCompile(`four`),
		regexp.MustCompile(`five`),
		regexp.MustCompile(`six`),
		regexp.MustCompile(`seven`),
		regexp.MustCompile(`eight`),
		regexp.MustCompile(`nine`),
	}
        numberRegex := regexp.MustCompile(`\d`)

	parsed := []int{}

	for range line {
		parsed = append(parsed, -1)
	}

        // find spelled out numbers
	for n, re := range regexes {
		for _, val := range re.FindAllStringIndex(line, -1) {
			start := val[0]
			parsed[start] = n
		}
	}

        // find digits
        for _, val := range numberRegex.FindAllStringIndex(line, -1) {
                start := val[0]
                parsed[start] = toInt(rune(line[start]))
        }

        // clean up the output
	out := []int{}
	for _, n := range parsed {
		if n != -1 {
			out = append(out, n)
		}
	}
	if len(out) == 0 {
		return nil, errors.New("Could not find any numbers in " + line)
	}

	return out, nil
}

func main() {
	// part 1
	lines := readLines("input.txt")
	total := 0
	for _, line := range lines {
		first, last := firstLastDigits(line)
		number := (10 * first) + last
		total += number
	}
	fmt.Println(total)

	// part 2
	lines = readLines("input.txt")
	total = 0
	for _, line := range lines {
		numbers, err := extractNumbers(line)
		if err != nil {
			panic(err)
		}
		first := numbers[0]
		last := numbers[len(numbers)-1]
		number := (10 * first) + last
		total += number
	}
	fmt.Println(total)
}
