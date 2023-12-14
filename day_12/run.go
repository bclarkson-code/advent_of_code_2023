package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	// "errors"
)

type count struct {
	char  rune
	count int
}

type Stack struct {
	items []interface{}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (item interface{}) {
	if len(s.items) == 0 {
		return nil
	}

	item = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	return
}

func (s *Stack) Peek() (item interface{}) {
	if len(s.items) == 0 {
		return nil
	}

	return s.items[len(s.items)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
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

func extractInts(input string) []int {
	re := regexp.MustCompile(`-?\d+`)

	matches := re.FindAllString(input, -1)

	var ints []int
	for _, match := range matches {
		num, err := strconv.Atoi(match)
		if err == nil {
			ints = append(ints, num)
		}
	}

	return ints
}

func parse(line string) (string, int, []int) {
	s := strings.Split(line, " ")

	extracted := extractInts(s[1])

	spaces := 0
	for i := 0; i < len(s[0]); i++ {
		if s[0][i] == '?' {
			spaces += 1
		}
	}
	return s[0], spaces, extracted
}

func Count(line string) (counts []count) {
	prev := rune(line[0])
	cur := count{rune(line[0]), 0}
	for _, c := range line {
		if prev == c {
			cur.count += 1
		} else {
			counts = append(counts, cur)
			cur = count{c, 1}
		}
		prev = c

	}
	counts = append(counts, cur)
	return
}

func isValid(line string, budget []int) bool {
	used := []count{}
	for _, l := range Count(line) {
		if l.char == '#' {
			used = append(used, l)
		}
	}

	if len(used) != len(budget) {
		return false
	}
	for i := 0; i < len(used); i++ {
		l := used[i]
		b := budget[i]

		if b != l.count {
			return false
		}
	}
	return true
}

func Combos(spaces int) (combos [][]int) {
	for n := 0; n < int(math.Pow(2, float64(spaces))); n++ {
		combo := []int{}
		for j := 0; j < spaces; j++ {
			bit := (n >> j) & 1
			if bit == 1 {
				combo = append(combo, '#')
			} else {
				combo = append(combo, '.')
			}
		}
		combos = append(combos, combo)
	}
	return combos
}

func findCombinations(l string, spaces int, budget []int) int {
	combos := Combos(spaces)
	count := 0
	for _, combo := range combos {
		cIdx := 0
		t := ""
		for _, c := range l {
			switch c {
			case '.', '#':
				t += string(c)
			case '?':
				t += string(combo[cIdx])
				cIdx++
			}
		}
		// fmt.Println(t, isValid(t, budget))
		if isValid(t, budget) {
			count += 1
		}
	}
	return count
}

func findCombinationsFaster(line string, budget []int, total int, soln []count) int {
	sString := ""
	if len(soln) > 0 {
		for _, s := range soln[1:] {
			for i := 0; i < s.count; i++ {
				sString += string(s.char)
			}
		}
	}

	if len(soln) == 0 {
		soln = append(soln, count{'.', 0})
	}
	s := soln[len(soln)-1]

	// no more characters
	if len(line) == 0 {
		switch s.char {
		case '.':
			if len(budget) == 0 {
				return total + 1
			}
			return total
		case '#':
			if budget[0] == s.count && len(budget) == 1 {
				return total + 1
			}
			return total
		}
	}
	c := line[0]
	line = line[1:]

	if c == '.' {
		switch s.char {
		case '.':
			soln[len(soln)-1].count += 1
			return findCombinationsFaster(line, budget, total, soln)
		case '#':
			if budget[0] != s.count {
				return total
			}
			budget = budget[1:]
			s = count{rune(c), 1}
			soln = append(soln, s)
			return findCombinationsFaster(line, budget, total, soln)
		}
	}
	if c == '#' {
		switch s.char {
		case '.':
			if len(budget) == 0 || budget[0] == 0 {
				return total
			}
			s = count{rune(c), 1}
			soln = append(soln, s)
			return findCombinationsFaster(line, budget, total, soln)
		case '#':
			if len(budget) == 0 || budget[0] == 0 {
				return total
			}
			soln[len(soln)-1].count += 1
			return findCombinationsFaster(line, budget, total, soln)
		}
	}
	if c == '?' {
		budgetCopy := make([]int, len(budget))
		copy(budgetCopy, budget)

		solnCopy := make([]count, len(soln))
		copy(solnCopy, soln)

		// try .
		line = string('.') + line
		total = findCombinationsFaster(line, budget, total, soln)

		// try #
		line = string('#') + line[1:]
		return findCombinationsFaster(line, budgetCopy, total, solnCopy)

	}
	panic(-1)
}

func repeat(s string, n int) string {
	if n == 0 {
		return ""
	}
	return s + repeat(s, n-1)
}

func repeatInts(s []int, n int) []int {
	if n == 0 {
		return []int{}
	}
	return append(s, repeatInts(s, n-1)...)
}

func main() {
	lines := readLines("day_12/example.txt")

	// part 1
	total := 0
	for i, line := range lines {
		fmt.Println(i, len(lines))
		l, _, budget := parse(line)
		count := findCombinationsFaster(l, budget, 0, []count{})
		total += count
	}
	fmt.Println(total)

	// part 2
	fmt.Println()
	total = 0
	for _, line := range lines {
		l, _, budget := parse(line)
		count := findCombinationsFaster(repeat(l, 5), repeatInts(budget, 5), 0, []count{})
		fmt.Println(count)
		total += count
	}
	fmt.Println(total)
}
