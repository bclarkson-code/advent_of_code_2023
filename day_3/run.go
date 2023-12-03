package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type number struct {
	start int
	end   int
	row   int
	str   string
}

func (n *number) val() int {
	val, err := strconv.Atoi(n.str)
	if err != nil {
		log.Fatal(err)
	}
	return val
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

func getNumbers(lines []string) (numbers []number) {
	re := regexp.MustCompile(`\d+`)
	for i, line := range lines {
		for _, n := range re.FindAllStringIndex(line, -1) {
			n := number{n[0], n[1], i, line[n[0]:n[1]]}
			numbers = append(numbers, n)
		}
	}
	return

}

func getGears(lines []string) (gears [][2]int){
	re := regexp.MustCompile(`\*`)
	for i, line := range lines {
		for _, n := range re.FindAllStringIndex(line, -1) {
			gear := [2]int{i, n[0]}
			gears = append(gears, gear)
		}
	}
	return
} 

func isSymbol(val byte) bool {
	if '0' <= val && val <= '9' {
		return false
	}
	if val == '.' {
		return false
	}
	return true
}

func (n *number) nextToSymbol(lines []string) bool {
	var left int
	var right int

	if n.start > 0 {
		left = n.start - 1
	} else {
		left = 0
	}

	if n.end < len(lines[0])-1 {
		right = n.end
	} else {
		right = len(lines[0]) - 1
	}

	// above
	if n.row > 0 {
		for i := left; i < right+1; i++ {
			if isSymbol(lines[n.row-1][i]) {
				return true
			}
		}
	}

	// left
	if isSymbol(lines[n.row][left]) {
		return true
	}

	// right
	if isSymbol(lines[n.row][right]) {
		return true
	}

	// below
	if n.row < len(lines)-2 {
		for i := left; i < right+1; i++ {
			if isSymbol(lines[n.row+1][i]) {
				return true
			}
		}
	}
	return false
}

func (n *number) nextToGear(lines []string, gear [2]int) bool {
	if n.row < gear[0] - 1 || n.row > gear[0] + 1{
		return false
	}
	if n.start > gear[1] + 1 || n.end < gear[1]{
		return false
	}
	return true
}


func main() {
	lines := readLines("day_3/input.txt")
	nums := getNumbers(lines)

	// part 1
	total := 0
	for _, n := range nums {
		if n.nextToSymbol(lines) {
			total += n.val()
		}
	}
	fmt.Println(total)

	// part 2
	total = 0
	gears := getGears(lines)
	for _, gear := range gears {
		var nextTo []int
		for _, num := range nums {
			if num.nextToGear(lines, gear) {
				nextTo = append(nextTo, num.val())	
			}
		}
		if len(nextTo) == 2 {
			total += nextTo[0] * nextTo[1]
		}
	}
	fmt.Println(total)

}
