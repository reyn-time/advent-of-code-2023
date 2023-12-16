package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("mirror.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var maze []string
	for scanner.Scan() {
		maze = append(maze, scanner.Text())
	}

	if !*secondFlag {
		// Part 1
		count := countEnergizedTiles(maze, 0, 0, 1)
		fmt.Println(count)
		return
	}
	// Part 2
	maxCount := 0
	h, w := len(maze), len(maze[0])
	for i := 0; i < h; i++ {
		count := countEnergizedTiles(maze, i, 0, 1)
		if count > maxCount {
			maxCount = count
		}
		count = countEnergizedTiles(maze, i, w-1, 3)
		if count > maxCount {
			maxCount = count
		}
	}
	for j := 0; j < w; j++ {
		count := countEnergizedTiles(maze, 0, j, 2)
		if count > maxCount {
			maxCount = count
		}
		count = countEnergizedTiles(maze, h-1, j, 0)
		if count > maxCount {
			maxCount = count
		}
	}
	fmt.Println(maxCount)
}

// Direction: 0: up, 1: right, 2: down, 3: left
func countEnergizedTiles(maze []string, startX, startY, startDir int) int {
	h, w := len(maze), len(maze[0])
	visited := make([][]int, h)
	for i := 0; i < h; i++ {
		visited[i] = make([]int, w)
	}
	visited[startX][startY] |= (1 << startDir)

	nextDirection := map[rune][][]int{
		'.':  {{0}, {1}, {2}, {3}},
		'/':  {{1}, {0}, {3}, {2}},
		'\\': {{3}, {2}, {1}, {0}},
		'-':  {{1, 3}, {1}, {1, 3}, {3}},
		'|':  {{0}, {0, 2}, {2}, {0, 2}},
	}

	offset := [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	stack := [][3]int{{startX, startY, startDir}}
	for len(stack) > 0 {
		// Pop from stack
		x, y, dir := stack[len(stack)-1][0], stack[len(stack)-1][1], stack[len(stack)-1][2]
		stack = stack[:len(stack)-1]

		// Main loop
		for _, nextDir := range nextDirection[rune(maze[x][y])][dir] {
			dx, dy := offset[nextDir][0], offset[nextDir][1]
			if x+dx < 0 || x+dx >= h || y+dy < 0 || y+dy >= w {
				// Out of range
				continue
			}
			if visited[x+dx][y+dy]&(1<<nextDir) != 0 {
				// Already visited (x+dx, y+dy, nextDir)
				continue
			}
			// Mark (x+dx, y+dy, nextDir) as visited
			visited[x+dx][y+dy] |= (1 << nextDir)
			stack = append(stack, [3]int{x + dx, y + dy, nextDir})
		}
	}

	// Count visited cells
	count := 0
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if visited[i][j] != 0 {
				count++
			}
		}
	}
	return count
}
