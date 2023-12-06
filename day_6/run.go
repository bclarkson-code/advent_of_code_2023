package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"math"
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

func extractInts(s string) (ints []int) {
	r := regexp.MustCompile(`\d+`)	
	for _, i := range r.FindAllString(s, -1) {
		v, err := strconv.Atoi(i)
		if err != nil {
			log.Fatal(err)
		}
		ints = append(ints, v)
	}
	return
}

func toInt(s string) int {
	r := regexp.MustCompile(`\d`)
	combined := ""
	for _, i := range r.FindAllString(s, -1) {
		combined += i
	}
	v, err := strconv.Atoi(combined)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func winRange(t int, d int) (min float64, max float64) {
	// I did a bit of algebra to figure out the answer without needing
	// to do a loop
	// # mafs
	a := float64(t) / 2.0	
	b := math.Sqrt(math.Pow(a, 2.) - float64(d))
	
	// add or subtract as little bit to handle getting
	// an exact integer 
	min = math.Ceil(a - b + 1e-10)	
	max = math.Floor(a + b - 1e-10)	

	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	return

}

func waysToWin(l float64, r float64) int {
	// I know that I dont need these checks but it hurts me to not
	if l == 0 && r == 0 {
		return 0
	}
	if l > r {
		return 0
	}
	return int(r - l + 1)
}

func main() {
	lines := readLines("day_6/input.txt")
	times := extractInts(lines[0])
	distances := extractInts(lines[1])

	// part 1
	score := 1
	for i :=0; i < len(times); i++ {
		l, r:= winRange(times[i], distances[i])
		score *= waysToWin(l, r)
	}
	fmt.Println(score)

	// part 2
	time := toInt(lines[0])
	distance := toInt(lines[1])
	l, r := winRange(time, distance)
	fmt.Println(waysToWin(l, r))
}
