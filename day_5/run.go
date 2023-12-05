package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type mapping struct {
	sources []int
	diffs   []int
	sizes   []int
	from    string
	to      string
}

type pair struct {
	min int
	max int
}

type op struct {
	min  int
	max  int
	diff int
}

// biggest number I can think of
const BIG_NUMBER int = 100000000000000

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

func chunk(lines []string) (chunks [][]string) {
	// group the lines by chunk
	var chunk []string
	for _, line := range lines {
		switch line {
		case "":
			chunks = append(chunks, chunk)
			chunk = []string{}
		default:
			chunk = append(chunk, line)
		}
	}
	chunks = append(chunks, chunk)
	return
}

func fromChunk(chunk []string) mapping {
	f := regexp.MustCompile(`(?P<from>[a-z]+)\-`)
	t := regexp.MustCompile(`to\-(?P<to>[a-z]+)`)
	n := regexp.MustCompile(`\d+`)

	fromRaw := f.FindStringSubmatch(chunk[0])
	if fromRaw == nil {
		log.Fatal("Could not parse 'from' from string")
	}
	from := fromRaw[1]

	toRaw := t.FindStringSubmatch(chunk[0])
	if fromRaw == nil {
		log.Fatal("Could not parse 'from' from string")
	}
	to := toRaw[1]

	var sources []int
	var diffs []int
	var sizes []int

	for _, line := range chunk[1:] {
		nums := n.FindAllString(line, -1)
		if nums == nil {
			log.Fatal("Could not parse numbers from string")
		}

		source, err := strconv.Atoi(nums[1])
		if err != nil {
			log.Fatal("Failed to parse source")
		}
		sources = append(sources, source)

		dest, err := strconv.Atoi(nums[0])
		if err != nil {
			log.Fatal("Failed to parse dest")
		}
		diffs = append(diffs, dest-source)

		size, err := strconv.Atoi(nums[2])
		if err != nil {
			log.Fatal("Failed to parse size")
		}
		sizes = append(sizes, size)
	}
	return mapping{sources, diffs, sizes, from, to}
}

func toSeeds(line string) (seeds []int) {
	r := regexp.MustCompile(`\d+`)
	nums := r.FindAllString(line, -1)
	if nums == nil {
		log.Fatal("Could not parse seeds")
	}
	for _, n := range nums {
		num, err := strconv.Atoi(n)
		if err != nil {
			log.Fatal("Could not parse seeds")
		}
		seeds = append(seeds, num)
	}
	return
}

func sortPairs(pairs []pair) []pair {
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].min < pairs[j].min })
	return pairs
}

func sortOps(ops []op) []op {
	sort.Slice(ops, func(i, j int) bool { return ops[i].min < ops[j].min })
	return ops
}

func fillGaps(pairs []pair) []pair {
	s := 0
	var filled []pair
	pairs = sortPairs(pairs)
	for _, p := range pairs {
		if p.min > s+1 {
			filled = append(filled, pair{s, p.min})
		}
		filled = append(filled, p)
		s = p.max
	}
	filled = append(filled, pair{s, BIG_NUMBER})
	return filled
}

func toPairs(seeds []int) (pairs []pair) {
	for i := 0; i < len(seeds); i += 2 {
		p := pair{seeds[i], seeds[i] + seeds[i+1] - 1}
		pairs = append(pairs, p)
	}
	// pairs = fillGaps(pairs)
	return
}

func (m mapping) apply(in int) int {
	for i := 0; i < len(m.sources); i++ {
		s := m.sources[i]
		d := m.diffs[i]
		n := m.sizes[i]
		if s <= in && in < s+n {
			return in + d
		}
	}
	return in
}

func (m mapping) applyArr(arr []int) []int {
	for i, v := range arr {
		arr[i] = m.apply(v)
	}
	return arr
}

func (m mapping) toOps() (ops []op) {
	var parsed []op
	for i := 0; i < len(m.sources); i++ {
		s := m.sources[i]
		n := m.sizes[i]
		d := m.diffs[i]

		p := op{s, s + n - 1, d}
		parsed = append(parsed, p)
	}
	parsed = sortOps(parsed)

	small := 0
	for _, p := range parsed {
		if p.min > small+1 {
			ops = append(ops, op{small, p.min - 1, 0})
		}
		ops = append(ops, p)
		small = p.max + 1
	}
	ops = append(ops, op{small, BIG_NUMBER, 0})

	return
}

func (m mapping) applyPair(p pair) (in []pair, out []pair) {
	for i := 0; i < len(m.sources); i++ {
		s := m.sources[i]
		n := m.sizes[i]
		if p.min > s+n-1 || p.max < s {
			continue
		}

		var l int
		var r int
		if p.min >= s {
			l = p.min
		} else {
			in = append(in, pair{p.min, s - 1})
			l = s
		}

		if p.max <= s+n-1 {
			r = p.max
		} else {
			in = append(in, pair{s + n - 1, p.max})
			r = s + n - 1
		}
		in = append(in, pair{l, r})
	}
	if len(in) == 0 {
		in = append(in, p)
	}
	for _, v := range in {
		out = append(out, pair{m.apply(v.min), m.apply(v.max)})
	}
	return
}

func (m mapping) applyPairs(pairs []pair) (in []pair, out []pair) {
	for _, p := range pairs {
		i, o := m.applyPair(p)
		fmt.Println(i, o, m)
		for j := 0; j < len(i); j++ {
			in = append(in, i[j])
			out = append(out, o[j])
		}
	}
	fmt.Println()
	return
}

func min(arr []int) int {
	lowest := BIG_NUMBER
	for _, v := range arr {
		if v < lowest {
			lowest = v
		}
	}
	return lowest
}

func (o op) splitAt(i int) (op, op) {
	if i < o.min || i > o.max+1 {
		log.Fatal("Cannot split")
	}
	left := op{o.min, i, o.diff}
	right := op{i + 1, o.max, o.diff}
	return left, right
}

func applyOps(ops []op, p pair) (result []pair) {
	var applies []op
	for _, o := range ops {
		if o.max < p.min {
			continue
		}
		if o.min > p.max{
			continue
		}
		if p.min <= o.min && p.max >= o.max{
			applies = append(applies, o)
			continue
		}
		toStore := op{o.min, o.max, o.diff}
		if o.min < p.min {
			toStore.min	 = p.min
		}
		if o.max > p.max {
			toStore.max= p.max
		}
		applies = append(applies, toStore)
	}
	for _, o := range applies {
		outPair := pair{o.min + o.diff, o.max + o.diff}
		result = append(result, outPair)
	}
	return
}

// Combine two successive op lists
func fuse(left []op, right []op) (fused []op) {
	il := 0
	ir := 0
	l := left[il]
	r := right[ir]
	var lSplit []op
	var rSplit []op
	fmt.Println()
	fmt.Println(left, right)
	fmt.Println()
	for {
		var n op

		if l.min > r.min {
			n, r = r.splitAt(l.min)
			rSplit = append(rSplit, n)
		} else if l.min < r.min {
			n, l = l.splitAt(r.min)
			lSplit = append(lSplit, n)
		} else {
			// do nothing
		}

		if l.max < r.max {
			n, r = r.splitAt(l.max)
			lSplit = append(lSplit, l)
			rSplit = append(rSplit, n)
			if il == len(left)-1 && ir == len(right)-1 {
				break
			}
			il++
			l = left[il]
		} else if l.max > r.max {
			n, l = l.splitAt(r.max)
			lSplit = append(lSplit, n)
			rSplit = append(rSplit, r)
			if il == len(left)-1 && ir == len(right)-1 {
				break
			}
			ir++
			r = right[ir]
		} else {
			lSplit = append(lSplit, l)
			rSplit = append(rSplit, r)
			if il == len(left)-1 && ir == len(right)-1 {
				break
			}
			il++
			ir++
			l = left[il]
			r = right[ir]
		}
	}

	for i := 0; i < len(lSplit); i++ {
		fusedOp := op{lSplit[i].min, lSplit[i].max, lSplit[i].diff + rSplit[i].diff}
		fused = append(fused, fusedOp)
	}
	return
}

func main() {
	lines := readLines("day_5/input.txt")
	chunks := chunk(lines)

	var maps []mapping
	for _, c := range chunks[1:] {
		maps = append(maps, fromChunk(c))
	}

	// part 1
	seeds := toSeeds(chunks[0][0])
	next := "seed"

	for {
		for _, m := range maps {
			if m.from == next {
				m.applyArr(seeds)
				next = m.to
				break
			}
		}
		if next == "location" {
			break
		}
	}
	fmt.Println(min(seeds))

	// part 2
	seeds = toSeeds(chunks[0][0])
	pairs := toPairs(seeds)

	next = "seed"

	combined := maps[0].toOps()
	for i := 1; i < len(maps); i++ {
		combined = fuse(combined, maps[i].toOps())
	}
	var out []pair
	for _, p := range pairs {
		for _, o := range applyOps(combined, p){
			out = append(out, o)
		}
	}

	var starts []int
	for _, p := range out {
		starts = append(starts, p.min)
	}
	fmt.Println(min(starts))

}
