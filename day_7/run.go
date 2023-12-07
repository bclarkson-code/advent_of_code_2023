package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Tuple struct {
	score int
	id    int
	bid   int
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

func parseLine(line string) (cards []int, bid int) {
	split := strings.SplitN(line, " ", 2)
	cardString := split[0]
	bidString := split[1]

	for _, c := range cardString {
		switch c {
		case 'A':
			cards = append(cards, 14)
		case 'K':
			cards = append(cards, 13)
		case 'Q':
			cards = append(cards, 12)
		case 'J':
			cards = append(cards, 11)
		case 'T':
			cards = append(cards, 10)
		default:
			n, err := strconv.Atoi(string(c))
			if err != nil {
				log.Fatal("Could not convert to int")
			}
			cards = append(cards, n)
		}
	}

	bid, err := strconv.Atoi(bidString)
	if err != nil {
		log.Fatal("Could not parse score")
	}

	return
}

func parseLinePart2(line string) (cards []int, bid int) {
	split := strings.SplitN(line, " ", 2)
	cardString := split[0]
	bidString := split[1]

	for _, c := range cardString {
		switch c {
		case 'A':
			cards = append(cards, 14)
		case 'K':
			cards = append(cards, 13)
		case 'Q':
			cards = append(cards, 12)
		case 'J':
			cards = append(cards, 1)
		case 'T':
			cards = append(cards, 10)
		default:
			n, err := strconv.Atoi(string(c))
			if err != nil {
				log.Fatal("Could not convert to int")
			}
			cards = append(cards, n)
		}
	}

	bid, err := strconv.Atoi(bidString)
	if err != nil {
		log.Fatal("Could not parse score")
	}

	return
}

func uniqueId(cards []int) (score int) {
	for i, c := range cards {
		score += c * int(math.Pow(14, float64(5-i)))
	}

	return score
}

func orderScore(t []Tuple, l int, r int) bool {
	if t[l].score > t[r].score {
		return false
	} else if t[l].score < t[r].score {
		return true
	}
	// exact match
	if t[l].id > t[r].id {
		return false
	}
	return true
}

func sortTuple(t []Tuple) []Tuple {
	sort.Slice(t, func(i, j int) bool {
		return orderScore(t, i, j)
	})
	return t
}

func Score(cards []int) (score int) {
	count := [13]int{}
	for _, c := range cards {
		count[c-2] += 1
	}

	pairs := []int{}
	three := -1
	for i, c := range count {
		// 5 of a kind
		if c == 5 {
			return 6
		}
		// 4 of a kind
		if c == 4 {
			return 5
		}
		if c == 3 {
			three = i + 1
		}
		if c == 2 {
			pairs = append(pairs, i+1)
		}
	}
	// full house
	if three != -1 && len(pairs) == 1 {
		return 4
	}
	// three of a kind
	if three != -1 && len(pairs) == 0 {
		return 3
	}
	// two pair
	if len(pairs) == 2 {
		return 2
	}
	// one pair
	if len(pairs) == 1 {
		return 1
	}

	return 0
}

func ScorePart2(cards []int) (score int) {
	count := [14]int{}
	for _, c := range cards {
		count[c-1] += 1
	}

	const FIVE_OF_A_KIND = 6
	const FOUR_OF_A_KIND = 5
	const FULL_HOUSE = 4
	const THREE_OF_A_KIND = 3
	const TWO_PAIR = 2
	const ONE_PAIR = 1
	const HIGH_CARD = 0

	jokers := count[0]
	pairs := []int{}
	three := false
	for i, c := range count {
		// 5 of a kind
		if c == 5 {
			return FIVE_OF_A_KIND
		}
		// 4 of a kind
		if c == 4 {
			if jokers == 1 || jokers == 4{
				return FIVE_OF_A_KIND
			}
			return FOUR_OF_A_KIND
		}
		if c == 3 {
			three = true
		}
		if c == 2 {
			pairs = append(pairs, i+1)
		}
	}
	// full house
	if three && len(pairs) == 1 {
		if jokers == 2 || jokers == 3 {
			return FIVE_OF_A_KIND
		}
		return FULL_HOUSE
	}
	// three of a kind
	if three && len(pairs) == 0 {
		if jokers == 3 || jokers == 1 {
			return FOUR_OF_A_KIND
		}
		return THREE_OF_A_KIND
	}
	// two pair
	if len(pairs) == 2 {
		if jokers == 1 {
			return FULL_HOUSE 
		} else if jokers == 2 {
			return FOUR_OF_A_KIND 
		}
		return TWO_PAIR
	}
	// one pair
	if len(pairs) == 1 {
		if jokers == 2 || jokers == 1 {
			return THREE_OF_A_KIND
		}
		return ONE_PAIR
	}
	if jokers == 1 {
		return ONE_PAIR
	}
	return HIGH_CARD
}

func main() {
	lines := readLines("day_7/input.txt")
	cards := [][]int{}
	bids := []int{}

	// part 1
	for _, l := range lines {
		c, b := parseLine(l)
		cards = append(cards, c)
		bids = append(bids, b)
	}
	pairs := []Tuple{}
	for i := 0; i < len(cards); i++ {
		p := Tuple{
			score: Score(cards[i]),
			id:    uniqueId(cards[i]),
			bid:   bids[i],
		}
		pairs = append(pairs, p)
	}
	pairs = sortTuple(pairs)

	score := 0
	for i, p := range pairs {
		score += p.bid * (i + 1)
	}
	fmt.Println(score)

	// part 2
	cards = [][]int{}
	bids = []int{}
	for _, l := range lines {
		c, b := parseLinePart2(l)
		cards = append(cards, c)
		bids = append(bids, b)
	}
	pairs = []Tuple{}
	for i := 0; i < len(cards); i++ {
		p := Tuple{
			score: ScorePart2(cards[i]),
			id:    uniqueId(cards[i]),
			bid:   bids[i],
		}
		pairs = append(pairs, p)
	}
	pairs = sortTuple(pairs)

	score = 0
	for i, p := range pairs {
		score += p.bid * (i + 1)
	}
	fmt.Println(score)
}
