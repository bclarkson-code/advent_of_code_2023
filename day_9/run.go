package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

func toInts(line string) []int {
    r := regexp.MustCompile(`[\-\d]+`)
    strs := r.FindAllString(line, -1)
    ints := []int{}

    for _, i := range strs {
        v, err := strconv.Atoi(i)

        if err != nil {
            log.Fatal(err)  
        }
        ints = append(ints, v)
    }
    
    return ints
}

func difference(ints []int) int {
    diffs := [][]int{ints}
    rowIdx := 0

    for {
        diff := []int{}
        row := diffs[rowIdx]
        for i:=0; i < len(row) - 1; i++ {
            diff = append(diff, row[i + 1] - row[i])
        }
        diffs = append(diffs, diff)

        allZero := true
        for _, d := range diff {
            if d != 0 {
                allZero = false
                break
            }
        }
        if allZero {
            break
        }
        rowIdx += 1
    }

    a := 0
    fmt.Println(diffs[len(diffs) -1], a)
    for i := len(diffs)-2; i >= 0; i-- {
        row := diffs[i]
        a = row[len(row) - 1] + a
        fmt.Println(row, a)
    }
    fmt.Println()

    return a
}

func reverse(arr []int) []int {
    out := []int{}
    for i := len(arr) - 1; i >=0; i-- {
        out = append(out, arr[i])
    }
    return out
}

func main() {
	lines := readLines("day_9/input.txt")
    total := 0
    // part 1
    for _, line := range lines {
         total += difference(toInts(line))
    }
    fmt.Println(total)

    // part 2
    total = 0
    for _, line := range lines {
         total += difference(reverse(toInts(line)))
    }
    fmt.Println(total)
}
