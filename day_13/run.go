package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type pair struct {
	left  int
	right int
	idx   int
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

func group(lines []string) (groups [][]string) {
	g := []string{}
	for _, l := range lines {
		if l == "" {
			groups = append(groups, g)
			g = []string{}
		} else {
			g = append(g, l)
		}
	}
	groups = append(groups, g)
	return
}

func hash(s string) int {
	h := 0
	for _, i := range s {
		switch i {
		case '.':
			h = h << 1
		case '#':
			h = h << 1
			h = h ^ 0x1
		}
	}
	return h
}

func rowHash(group []string) (hashes []int) {
	// create a unique hash for each row
	for _, l := range group {
		hashes = append(hashes, hash(l))
	}
	return
}

func colHash(group []string) (hashes []int) {
	for i := 0; i < len(group[0]); i++ {
		s := ""
		for j := 0; j < len(group); j++ {
			s += string(group[j][i])
		}
		hashes = append(hashes, hash(s))
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func findReflection(hashes []int) (left int, right int, ok bool) {
	for i := 0; i < len(hashes)-1; i++ {
		for j := 0; j < len(hashes); j++ {
			l := i - j
			r := i + j + 1

			// check for no match
			if hashes[l] != hashes[r] {
				break
			}

			// reached left
			if l == 0 {
				return i + 1, len(hashes) - i + 1, true
			}
			// reached right
			if r == len(hashes)-1 {
				return i + 1, len(hashes) - i + 1, true
			}
		}
	}
	return -1, -1, false
}

func findOffByOne(hashes []int) (miss pair, ok bool) {
	for i := 0; i < len(hashes)-1; i++ {
		miss := pair{-1, -1, -1}

		for j := 0; j < len(hashes); j++ {
			l := i - j
			r := i + j + 1

			// check for no match
			if hashes[l] != hashes[r] {
				if miss.idx == -1 {
					miss = pair{hashes[l], hashes[r], i + 1}
				} else {
					break
				}
			}

			// reached end
			if l == 0 || r == len(hashes)-1 {
				if miss.idx != -1 {
					if miss.oneAway() {
						return miss, true
					}
				}
				break
			}
		}
	}
	return miss, false
}

// Check whether left can be converted to right with
// a single bit flip
func (p *pair) oneAway() bool {
	diff := p.left ^ p.right
	// check if diff contains exactly one 1 (power of 2)
	return (diff & (diff - 1)) == 0
}

func main() {
	lines := readLines("day_13/input.txt")
	groups := group(lines)

	// part 1
	total := 0
	for i, g := range groups {
		rows := rowHash(g)
		l, _, ok := findReflection(rows)
		if ok {
			total += l * 100
			continue
		}

		cols := colHash(g)
		l, _, ok = findReflection(cols)
		if ok {
			total += l
			continue
		}
		fmt.Println(i, g, rows, cols)
		panic(-1)
	}
	fmt.Println(total)

	// part 2
	total = 0
	for i, g := range groups {
		rows := rowHash(g)
		p, ok := findOffByOne(rows)
		if ok {
			total += p.idx * 100
			continue
		}

		cols := colHash(g)
		p, ok = findOffByOne(cols)
		if ok {
			total += p.idx
			continue
		}
		fmt.Println(i, g, rows, cols)
		panic(-1)
	}
	fmt.Println(total)
}
