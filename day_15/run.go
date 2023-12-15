package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type  lense struct {
	code string
	focalLength int
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

func hashRune(r rune, currentValue int) int {
	currentValue += int(r)
	currentValue *= 17
	return currentValue % 256
}

func hashStr(s string, currentValue int) int {
	for _, r := range s {
		currentValue = hashRune(r, currentValue)	
	}
	return currentValue
}

func toLense(s string) (l lense, shouldRemove bool) {
	r := regexp.MustCompile(`[a-z]+`)
	code := r.FindString(s)
	a := regexp.MustCompile(`\-`)

	if a.MatchString(s){
		l = lense{code, -1}
		return l, true
	}

	n := regexp.MustCompile(`[0-9]+`)
	fl, err:= strconv.Atoi(n.FindString(s))
	if err != nil{
		log.Fatal("Could not parse focal length")
	}

	l = lense{code, fl}
	return l, false 

	
}

type box struct {
	b []lense
}

func (b *box) remove(l lense) {
	left := []lense{}
	for i :=0; i < len(b.b); i++ {
		if b.b[i].code == l.code {
			left = b.b[:i]


			if i + 1 == len(b.b){
				b.b = left
				break
			}

			right := b.b[i+1:]
			b.b = append(left, right...)
			break
		}
	}
}

func (b *box) add(l lense) {
	for i :=0; i < len(b.b); i++ {
		if b.b[i].code == l.code {
			b.b[i] = l
			return
		}
	}
	b.b = append(b.b, l)
}

func (b box)score(i int) int{
	total := 0
	for j := 0; j< len(b.b); j++ {
		total += (i + 1) * (j + 1) * b.b[j].focalLength
	}
	return total
}


func main(){
	line := readLines("day_15/input.txt")[0]

	// part 1
	total := 0
	for _, s := range strings.Split(line, ","){
		currentValue := 0
		currentValue = hashStr(s, currentValue)
		total += currentValue
	}
	fmt.Println(total)

	// part 2
	boxes := [256]box{}
	for _, s := range strings.Split(line, ","){
		l, shouldRemove := toLense(s)
		idx := hashStr(l.code, 0)

		if shouldRemove{
			boxes[idx].remove(l)
		} else {
			boxes[idx].add(l)
		}
	}
	total = 0
	for i, b := range boxes{
		total += b.score(i)
	}
	fmt.Println(total)

}
