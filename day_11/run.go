package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	// "errors" // I wont have any errors so dont need this, 100% bug free!
)

type loc struct {
	x int
	y int
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

func toRunes(lines []string) (runes [][]rune) {

	for _, l := range lines {
		runes = append(runes, []rune(l))

	}
	return
}

func expand(runes [][]rune) (expanded [][]rune) {
	// expand rows
	for _, row := range runes {
		expanded = append(expanded, row)
		empty := true
		for _, c := range row {
			if c == '#' {
				empty = false
				break
			}
		}
		if empty {
			expanded = append(expanded, row)
		}
	}

	// expand columns
	newCol := []rune{}
	offset := 0
	for k := 0; k < len(expanded); k++{
		newCol = append(newCol, '.')	
	}

	for i := 0; i < len(runes[0]); i++ {
		empty := true
		for j := 0; j < len(runes); j++ {
			if runes[j][i] == '#' {
				empty = false
				break
			}
		}

		if empty {
			for k := 0; k < len(expanded); k++{
				left, right := expanded[k][:i + offset], expanded[k][i + offset:]
				newRow := []rune{}
				newRow = append(newRow, left...)
				newRow = append(newRow, '.')
				newRow = append(newRow, right...)
				expanded[k] = newRow
			}
			offset += 1
		}
	}

	return
}

func emptyLines(runes [][]rune) (rows []int, cols []int) {
	for i, row := range runes {
		empty := true
		for _, c := range row {
			if c == '#' {
				empty = false
				break
			}	
		}
		if empty {
			rows = append(rows, i)
		}
	}	

	for i := 0; i < len(runes[0]); i++ {
		empty := true
		for _, row := range  runes {
			if row[i] == '#' {
				empty = false
				break
			}
		}
		if empty {
			cols = append(cols, i)
		}
	}

	return
}

func getLocs(runes [][]rune) (locs []loc) {
	for y, row := range runes {
		for x, c := range row {
			if c == '#' {
				locs = append(locs, loc{x, y})
			}
		}
	}
	return locs
}

func abs(val int) int {
	if val < 0 {
		return val * -1
	}
	return val
}

func (l loc) distanceTo(r loc) int {
	return abs(l.x - r.x) + abs(l.y - r.y)
}

func (l loc) expandedDistanceTo(r loc, eRows []int, eCols []int, m int) int {
	eRowsBetween := 0
	eColsBetween := 0

	minRow := 0
	maxRow := 0
	if r.y < l.y {
		minRow = r.y
		maxRow = l.y
	} else {
		minRow = l.y
		maxRow = r.y
	}

	minCol := 0
	maxCol := 0
	if r.x < l.x {
		minCol = r.x
		maxCol = l.x
	} else {
		minCol = l.x
		maxCol = r.x
	}

	for _, row := range eRows {
		if row > minRow && row < maxRow {
			eRowsBetween += 1	
		} else if row > maxRow {
			break
		}
	}

	for _, col := range eCols {
		if col > minCol && col < maxCol {
			eColsBetween += 1	
		} else if col > maxCol {
			break
		}
	}

	rowDist := abs(l.x - r.x) + m * eRowsBetween
	colDist := abs(l.y - r.y) + m * eColsBetween 
	
	return rowDist + colDist
}

func main() {
	lines := readLines("day_11/input.txt")

	// part 1
	runes := toRunes(lines)
	runes = expand(runes)
	fmt.Println()
	for _, l := range runes {
		fmt.Println(string(l))
	}

	total := 0
	locs := getLocs(runes)
	for i, l := range locs[:len(locs)-1] {
		for _, r := range locs[i:] {
			total += l.distanceTo(r)
		}
	}
	fmt.Println(total)

	// part 2
	total = 0
	runes = toRunes(lines)
	rows, cols := emptyLines(runes)
	locs = getLocs(runes)
	for i, l := range locs[:len(locs)-1] {
		for _, r := range locs[i:] {
			d := l.expandedDistanceTo(r, rows, cols, 1_000_000-1)
			total += d
		}
	}
	fmt.Println(total)
}
