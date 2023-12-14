package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

type pair struct {
	x int
	y int
}

func toRunes(l []string) (s [][]rune) {
	for _, row := range l {
		parsed := []rune{}
		for _, c := range row {
			parsed = append(parsed, c)
		}
		s = append(s, parsed)
	}
	return s
}

func slideUp(s [][]rune) [][]rune {
	for i := 0; i < len(s); i++ {
		for j := i; j > 0; j-- {
			for k := 0; k < len(s[0]); k++ {
				if s[j][k] == 'O' && s[j-1][k] == '.' {
					s[j][k] = '.'
					s[j-1][k] = 'O'
				}
			}
		}
	}

	return s
}

func slideDown(s [][]rune) [][]rune {
	for i := len(s)-1; i >= 0; i-- {
		for j := i; j < len(s)-1; j++ {
			for k := 0; k < len(s[0]); k++ {
				if s[j][k] == 'O' && s[j+1][k] == '.' {
					s[j][k] = '.'
					s[j+1][k] = 'O'
				}
			}
		}
	}

	return s
}

func slideLeft(s [][]rune) [][]rune {
	for i := 0; i < len(s[0]); i++ {
		for j := i; j > 0; j-- {
			for k := 0; k < len(s); k++ {
				if s[k][j] == 'O' && s[k][j-1] == '.' {
					s[k][j] = '.'
					s[k][j-1] = 'O'
				}
			}
		}
	}

	return s
}

func slideRight(s [][]rune) [][]rune {
	for i := len(s[0])-1; i >= 0; i-- {
		for j := i; j < len(s[0])-1; j++ {
			for k := 0; k < len(s); k++ {
				if s[k][j] == 'O' && s[k][j+1] == '.' {
					s[k][j] = '.'
					s[k][j+1] = 'O'
				}
			}
		}
	}

	return s
}

func hash(l [][]rune) (h []pair) {
	for i, row := range l {
		for j, c := range row {
			if c == 'O' {
				h = append(h, pair{i, j})
			}
		}
	}
	return h
}

func findWeight(s [][]rune) int {
	total := 0
	for i, row := range s {
		for _, c := range row {
			if c == 'O' {
				total += len(s) - i
			}
		}
	}
	return total
}

func inHashes(h []pair, hashes [][]pair) (int, bool) {
	for idx, prevHash := range hashes {
		foundMatch := true
		for i := 0; i < len(prevHash); i++ {
			if h[i].x != prevHash[i].x || h[i].y != prevHash[i].y {
				foundMatch = false
				break
			}
		}
		if foundMatch {
			return idx, true
		}
	}
	return -1, false
}

func main() {
	lines := readLines("day_14/input.txt")
	r := toRunes(lines)
	s := slideUp(r)
	fmt.Println(findWeight(s))

	// part 2
	s = toRunes(lines)
	hashes := [][]pair{hash(s)}

	i := 0
	for {
		
		s = slideUp(s)
		s = slideLeft(s)
		s = slideDown(s)
		s = slideRight(s)


		offset, inHash := inHashes(hash(s), hashes)
		if inHash {
			cycleLen := len(hashes) - offset
			cycleIdx := (1000000000 - offset) % cycleLen

			fmt.Println(cycleIdx)
			fmt.Printf("Length: %d, Index: %d, Offset: %d\n", cycleLen, cycleIdx, offset)

			// do some extra cycles to check visually we
			// are at the right idx
			for j := 0; j < cycleIdx + cycleLen*3; j++ {
				s = slideUp(s)
				s = slideLeft(s)
				s = slideDown(s)
				s = slideRight(s)
				fmt.Println(findWeight(s))
				if j % cycleLen  == 0{
					fmt.Println()
				}
			}

			// o, _ := inHashes(hash(s), hashes)
			// fmt.Println(o)

			fmt.Println()
			fmt.Println(findWeight(s))
			break
		}

		hashes = append(hashes, hash(s))
		i += 1
	}
}
