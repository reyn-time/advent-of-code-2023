package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
)

type cell struct {
	row, col int
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("galaxy.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Scan values into cells
	scanner := bufio.NewScanner(file)
	i := 0
	var cells []cell
	rowMap := make(map[int]int) // row -> new row number after shift
	colMap := make(map[int]int) // col -> new col number after shift
	for scanner.Scan() {
		line := scanner.Text()
		for j, c := range line {
			if c == '#' {
				cells = append(cells, cell{i, j})
				rowMap[i] = i
				colMap[j] = j
			}
		}
		i++
	}

	expandMap(rowMap, *secondFlag)
	expandMap(colMap, *secondFlag)

	// Compute the sum of all pair shortest path
	sum := 0
	for i, c1 := range cells {
		for _, c2 := range cells[i+1:] {
			sum += abs(rowMap[c1.row]-rowMap[c2.row]) + abs(colMap[c1.col]-colMap[c2.col])
		}
	}
	fmt.Println(sum)
}

func expandMap(m map[int]int, second bool) {
	var vs []int
	for k := range m {
		vs = append(vs, k)
	}
	slices.Sort(vs)
	currSpace := 0
	for i, v := range vs {
		// If the closest value is not adjacent, the space multiplies
		if i != 0 {
			gap := v - vs[i-1] - 1
			if second {
				currSpace += gap * 999999
			} else {
				currSpace += gap
			}
		}
		m[v] = v + currSpace
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
