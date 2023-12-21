package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("step.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var puzzle [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		puzzle = append(puzzle, []rune(scanner.Text()))
	}

	startY, startX := 0, 0
	for y, row := range puzzle {
		for x, char := range row {
			if char == 'S' {
				startY, startX = y, x
			}
		}
	}

	// Assumption: Let f(n) be the number of points that can be reached in n steps.
	// I assume that for all i, f(i), f(i+d), f(i+2d)... forms a quadratic sequence. (d = batch size)
	// We only need three batches to determine the quadratic sequence.
	batchSize := len(puzzle)
	batchCount := 3
	n := batchSize * batchCount
	points := []Point{{startX, startY}}
	numOfPointsVisited := make([][]int, batchCount)
	for i := 0; i < batchCount; i++ {
		numOfPointsVisited[i] = make([]int, batchSize)
	}

	for i := 0; i < n; i++ {
		visited := map[Point]bool{}
		for _, point := range points {
			for _, offset := range []Point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
				newPoint := Point{point.x + offset.x, point.y + offset.y}
				refX := newPoint.x % len(puzzle[0])
				if refX < 0 {
					refX += len(puzzle[0])
				}
				refY := newPoint.y % len(puzzle)
				if refY < 0 {
					refY += len(puzzle)
				}
				if puzzle[refY][refX] == '#' {
					continue
				}
				visited[newPoint] = true
			}
		}
		points = []Point{}
		for point := range visited {
			points = append(points, point)
		}

		numOfPointsVisited[i/batchSize][i%batchSize] = len(points)
	}

	steps := 64
	if *secondFlag {
		steps = 26501365
	}
	refIndex := (steps - 1) % batchSize
	row := (steps - 1) / batchSize
	// Suppose the quadratic sequence is f(n) = an^2 + bn + c.
	// f(0) = c
	// f(1) = a + b + c
	// f(2) = 4a + 2b + c
	// a = (f(2) - 2f(1) + f(0)) / 2
	// b = f(1) - f(0) - a
	// c = f(0)
	f0 := numOfPointsVisited[0][refIndex]
	f1 := numOfPointsVisited[1][refIndex]
	f2 := numOfPointsVisited[2][refIndex]
	a := (f2 - 2*f1 + f0) / 2
	b := f1 - f0 - a
	c := f0
	fmt.Println(a*row*row + b*row + c)
}
