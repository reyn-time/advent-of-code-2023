package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Cell struct {
	i, j, direction, stepLeft int
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("crucible.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var heatMap [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var row []int
		for _, c := range scanner.Text() {
			row = append(row, int(c-'0'))
		}
		heatMap = append(heatMap, row)
	}

	if !*secondFlag {
		fmt.Println(minHeatAtBottomRight(heatMap, 3, 1))
	} else {
		fmt.Println(minHeatAtBottomRight(heatMap, 10, 4))
	}
}

// Returns the minimum heat dissipated when you end at the bottom right corner.
// maxMove is the maximum number of blocks that the crucible can move in one direction.
// minMove is the minimum number of blocks the crucible must move in one direction.
func minHeatAtBottomRight(heatMap [][]int, maxMove, minMove int) int {
	// Directions: 0: up, 1: right, 2: down, 3: left
	infinity := 1000000000
	h, w := len(heatMap), len(heatMap[0])
	turns := [][]int{{1, 3}, {0, 2}, {1, 3}, {0, 2}}
	dist := [][][4][10]int{}

	for i := 0; i < h; i++ {
		var row [][4][10]int
		for j := 0; j < w; j++ {
			var cell [4][10]int
			for k := 0; k < 4; k++ {
				for l := 0; l < 10; l++ {
					cell[k][l] = infinity
				}
			}
			row = append(row, cell)
		}
		dist = append(dist, row)
	}

	// Run SPFA for each cell
	queue := []Cell{{0, 1, 1, maxMove - 1}, {1, 0, 2, maxMove - 1}}
	inQueue := map[Cell]bool{}
	inQueue[Cell{0, 1, 1, maxMove - 1}] = true
	inQueue[Cell{1, 0, 2, maxMove - 1}] = true
	dist[0][1][1][maxMove-1] = heatMap[0][1]
	dist[1][0][2][maxMove-1] = heatMap[1][0]
	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]
		inQueue[cell] = false

		// Can turn only when you have moved at least minMove blocks
		if maxMove-cell.stepLeft >= minMove {
			for _, nextDir := range turns[cell.direction] {
				next := Cell{cell.i, cell.j, nextDir, maxMove - 1}
				next.i += []int{-1, 0, 1, 0}[next.direction]
				next.j += []int{0, 1, 0, -1}[next.direction]
				if next.i < 0 || next.i >= h || next.j < 0 || next.j >= w {
					// Out of bounds
					continue
				}
				if nextHeat := dist[cell.i][cell.j][cell.direction][cell.stepLeft] + heatMap[next.i][next.j]; nextHeat < dist[next.i][next.j][next.direction][next.stepLeft] {
					dist[next.i][next.j][next.direction][next.stepLeft] = nextHeat
					if !inQueue[next] {
						queue = append(queue, next)
						inQueue[next] = true
					}
				}
			}
		}

		// Can move forward if you still have steps left
		if cell.stepLeft > 0 {
			next := Cell{cell.i, cell.j, cell.direction, cell.stepLeft - 1}
			next.i += []int{-1, 0, 1, 0}[next.direction]
			next.j += []int{0, 1, 0, -1}[next.direction]
			if next.i < 0 || next.i >= h || next.j < 0 || next.j >= w {
				// Out of bounds
				continue
			}
			if nextHeat := dist[cell.i][cell.j][cell.direction][cell.stepLeft] + heatMap[next.i][next.j]; nextHeat < dist[next.i][next.j][next.direction][next.stepLeft] {
				dist[next.i][next.j][next.direction][next.stepLeft] = nextHeat
				if !inQueue[next] {
					queue = append(queue, next)
					inQueue[next] = true
				}
			}
		}
	}

	// Find the minimum heat at (h-1, w-1)
	minHeat := infinity
	for i := 0; i < 4; i++ {
		for j := 0; j <= maxMove-minMove; j++ {
			if dist[h-1][w-1][i][j] < minHeat {
				minHeat = dist[h-1][w-1][i][j]
			}
		}
	}

	return minHeat
}
