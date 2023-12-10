package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"errors"
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

func neighbours(l loc, pipes []string) ([2]loc, error) {
	x := l.x
	y := l.y
	pipe := pipes[y][x]

	switch pipe {
	case '|':
		return [2]loc{loc{x, y-1}, loc{x, y + 1}}, nil
	case '-':
		return [2]loc{loc{x - 1,y}, loc{ x + 1, y}}, nil
	case 'L':
		return [2]loc{loc{x,y - 1}, loc{ x + 1, y}}, nil
	case 'J':
		return [2]loc{loc{x,y - 1}, loc{ x - 1, y}}, nil
	case '7':
		return [2]loc{loc{x - 1,y}, loc{ x, y + 1}}, nil
	case 'F':
		return [2]loc{loc{x + 1,y}, loc{ x, y + 1}}, nil
	case '.':
		return [2]loc{loc{-1,-1}, loc{-1, -1}}, errors.New("Outside of pipe system")
	case 'S':
		out := [2]loc{}
		outIdx := 0
		dirs := [4]loc{
			loc{x + 1, y},
			loc{x - 1, y},
			loc{x, y + 1},
			loc{x, y - 1},
		}
		for _, d := range dirs {
			// out of bounds
			if d.x < 0 || d.x > len(pipes[0])-1 || d.y < 0 || d.y > len(pipes)-1 {
				continue
			}
			possibleNeighbours, err := neighbours(d, pipes)
			if err != nil {
				continue
			}
			for i := 0; i < 2; i++{
				n := possibleNeighbours[i]
				if n.x == l.x && n.y == l.y {
					out[outIdx] = d
					outIdx += 1
					break
				}
			}
		}
		return out, nil
	}
	return [2]loc{loc{-1,-1}, loc{-1, -1}}, errors.New("Could not handle input")
}

func follow(pipes []string) (distances [][]int, maxDist int) {
	sx, sy := -1, -1
	for i := 0; i < len(pipes); i++ {
		row := []int{}
		for j := 0; j < len(pipes[0]); j++ {
			if pipes[i][j] == byte('S') {
				sx = j
				sy = i
			}
			row = append(row, -1)
		}
		distances = append(distances, row)
	}
	l := loc{sx, sy}
	distances[sy][sx] = 0
	maxDist = 0
	// FIFO queue
	queue := []loc{l}

	for {
		if len(queue) == 0 {
			break
		}
		l, queue = queue[0], queue[1:]
		nextTo, err := neighbours(l, pipes) 
		if err != nil {
			log.Fatal(err)
		}
		
		for i:=0; i<2; i++ {
			n := nextTo[i]
			// out of bounds
			if n.x < 0 || n.x > len(distances[0]) || n.y < 0 || n.y > len(distances) {
				continue
			}
			// not yet visited
			if distances[n.y][n.x] == -1 {
				distances[n.y][n.x] = distances[l.y][l.x] + 1
				if distances[l.y][l.x] + 1 > maxDist {
					maxDist = distances[l.y][l.x] + 1
				}
				queue = append(queue, n)
			}	
		}
	}

	return distances, maxDist

}

func extractPipe(distances [][]int, pipes []string) (pipe []string) {
	for y, d := range distances{
		l := ""
		for x, c := range d {
			switch c{
				case -1:
					l += " "
				default:
					l += string(pipes[y][x])
			}
		}
		pipe = append(pipe, l)
	}
	return
}

func expandGrid(pipe []string) (expanded []string){
	fmt.Println(len(pipe), len(pipe[0]))
	for y := 0; y < len(pipe); y++ {
		row := ""
		nextRow := ""
		for x := 0; x < len(pipe[0]); x++ {
			l := pipe[y][x]
			
			// fill gaps in current row
			row += string(l)
			switch l {
			case ' ', '7', 'J', '|':
				row += " "
			case 'F', 'L', '-', 'S':
				row += "-"
			}
			// fill gaps in next row
			switch l {
			case ' ', 'L', 'J', '-':
				nextRow += "  "
			case 'F', '7', '|', 'S':
				nextRow += "| "
			}
			
		}
		expanded = append(expanded, row)
		expanded = append(expanded, nextRow)
	}
	return 
}

func (l loc)neighbours(grid [][]rune) (validNeighbours[]loc) {
	dirs := []loc{
		loc{-1, 0},
		loc{0, -1},
		loc{1, 0},
		loc{0, 1},
	}
	for _, d := range dirs {
		n := loc{l.x + d.x, l.y + d.y}
		// check bounds
		if n.x < 0 || n.x >= len(grid[0]) || n.y < 0 || n.y >= len(grid) {
			continue
		}
		validNeighbours = append(validNeighbours, n)

	}
	return
}
func printGrid(grid [][]rune) {
	for _, row := range grid {
		for _, char := range row {
			switch char {
			case 'I':
				fmt.Printf("\x1b[31m%s\x1b[0m", string(char)) // Red color for 'I'
			case ' ':
				fmt.Print(string(char))
			default:
				fmt.Printf("\x1b[32m%s\x1b[0m", string(char)) // Green for pipes
			}
		}
		fmt.Println()
	}
}

func moveCursorUp(lines int) {
	for i := 0; i < lines; i++ {
		fmt.Print("\033[F") // Move cursor up one line
	}
}
func floodFill(startGrid []string, start loc) ([][]rune){
	l := start
	queue := []loc{l}
	grid := [][]rune{}

	for _, row := range startGrid {
		runes := []rune(row)
		grid = append(grid, runes)
		fmt.Println(string(runes))
	}

	for {
		if len(queue) == 0 {
			break
			fmt.Println(queue)
		}

		l, queue = queue[0], queue[1:]
		if grid[l.y][l.x] == ' ' {
			grid[l.y][l.x] = 'I'

			for _, n := range l.neighbours(grid) {
				queue = append(queue, n)
			}
		}
		// printGrid(grid)
		// moveCursorUp(len(grid))
	}
	return grid
	
}

func undoExpand(expanded [][]rune) (grid [][]rune){
	// Reverse expansion back to the original grid with the correct 
	// squares filled in
	for y := 0; y < len(expanded); y+=2 {
		row := []rune{}

		for x := 0; x < len(expanded[0]); x+=2 {
			if y % 2 == 0 && x % 2 == 0{
				row = append(row, expanded[y][x])
			}
		}
		grid = append(grid, row)
	}
	return
}

func main() {
	pipes := readLines("day_10/input.txt")

	// part 1
	distances, maxDist := follow(pipes) 
	fmt.Println(maxDist)
	
	// part 2
	pipe := extractPipe(distances, pipes)
	expanded := expandGrid(pipe)
	filled := floodFill(expanded, loc{140, 140})

	squished := undoExpand(filled)
	printGrid(squished)
	total := 0
	for _, row := range squished {
		for _, c := range row {
			if c == 'I' {
				total += 1
			}	
		}
	}
	fmt.Println(total)
}
