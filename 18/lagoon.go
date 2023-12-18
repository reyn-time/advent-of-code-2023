package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Line struct {
	start, end [2]int
	l, r       int
	direction  string
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("lagoon.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	currX, currY := 0, 0

	path := []Line{}
	xSet := map[int]bool{}
	perimeter := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, steps := getLineFromString(scanner.Text(), currX, currY, *secondFlag)
		currX = line.end[0]
		currY = line.end[1]
		perimeter += steps
		path = append(path, line)

		if line.direction == "L" || line.direction == "R" {
			xSet[currX] = true
		}
	}

	xs := []int{}
	for x := range xSet {
		xs = append(xs, x)
	}
	slices.Sort(xs)

	area := 0
	for i, x := range xs {
		area += rowArea(path, x)
		if i > 0 && xs[i-1]+1 < x {
			area += rowArea(path, xs[i-1]+1) * (x - xs[i-1] - 1)
		}
	}
	fmt.Println(perimeter + area)
}

func getLineFromString(s string, currX, currY int, second bool) (Line, int) {
	fields := strings.Fields(s)

	if !second {
		offset := map[string][2]int{
			"R": {0, 1},
			"L": {0, -1},
			"U": {-1, 0},
			"D": {1, 0},
		}[fields[0]]
		steps, _ := strconv.Atoi(fields[1])
		return Line{
			start:     [2]int{currX, currY},
			end:       [2]int{currX + offset[0]*steps, currY + offset[1]*steps},
			l:         min(currY, currY+offset[1]*steps),
			r:         max(currY, currY+offset[1]*steps),
			direction: fields[0],
		}, steps
	}

	steps64, _ := strconv.ParseInt(fields[2][2:7], 16, 64)
	steps := int(steps64)
	offset := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}[int(fields[2][7])-int('0')]
	return Line{
		start:     [2]int{currX, currY},
		end:       [2]int{currX + offset[0]*steps, currY + offset[1]*steps},
		l:         min(currY, currY+offset[1]*steps),
		r:         max(currY, currY+offset[1]*steps),
		direction: []string{"R", "D", "L", "U"}[int(fields[2][7])-int('0')],
	}, steps
}

// Compute the number of cells enclosed by the path at row i.
func rowArea(path []Line, i int) int {
	area := 0
	relevantLines := []Line{}
	for _, line := range path {
		if (line.start[0] <= i && line.end[0] >= i) || (line.start[0] >= i && line.end[0] <= i) {
			relevantLines = append(relevantLines, line)
		}
	}
	slices.SortFunc(relevantLines, func(a, b Line) int {
		if a.l != b.l {
			return a.l - b.l
		}
		return a.r - b.r
	})

	flips := 0
	prevDirection := ""
	for i, line := range relevantLines {
		if line.direction == "U" || line.direction == "D" {
			if prevDirection != line.direction {
				flips++
				if flips%2 == 0 {
					area += max(line.l-relevantLines[i-1].r-1, 0)
				}
				prevDirection = line.direction
			}
		}
	}

	return area
}
