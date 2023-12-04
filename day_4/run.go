package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

type line struct {
	idx   int
	left  []int
	right []int
}

func (l *line) Matches() (matches []int) {
	for _, left := range l.left {
		for _, right := range l.right {
			if left == right {
				matches = append(matches, left)
				break
			}
		}
	}
	return
}

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

func split(line string) (idx string, left string, right string) {
	inIdx := true
	inLeft := true
	for _, c := range line {
		switch c {
		case ':':
			inIdx = false
		case '|':
			inLeft = false
		default:
			if inIdx {
				idx += string(c)
			} else if inLeft {
				left += string(c)
			} else {
				right += string(c)
			}
		}
	}
	return
}

func extractInts(s string) (ints []int) {
	r := regexp.MustCompilePOSIX(`[0-9]+`)
	for _, i := range r.FindAllString(s, -1) {
		val, err := strconv.Atoi(i)
		if err != nil {
			log.Fatal(err)
		}
		ints = append(ints, val)
	}
	return
}

func points(ints []int) int {
	l := len(ints)
	if l == 0 {
		return 0
	}
	return int(math.Pow(2, float64(l-1)))
}

func min(a int, b int) int {
	if a < b{
		return a
	}
	return b

}

func main() {
	rawLines := readLines("day_4/input.txt")
	var lines []line

	for _, l := range rawLines {
		idx, left, right := split(l)
		parsedLine := line{
			idx:   extractInts(idx)[0],
			left:  extractInts(left),
			right: extractInts(right),
		}
		lines = append(lines, parsedLine)

	}

	// Part 1
	score := 0
	for _, l := range lines {
		score += points(l.Matches())
	}
	fmt.Println(score)

	// Part 2
	counts := make([]int, len(lines))
	for i := 0; i < len(lines); i++ {
		counts[i] = 1	
	}

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		c := counts[i]

		m := len(l.Matches())
		for j := i+1; j < min(i+1+m, len(lines)); j++ {
			counts[j] += c 
		}
	}
	
	score = 0
	for _, val := range counts {
		score += val
	}
	fmt.Println(score)
}
