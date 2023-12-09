package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("oasis.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	total := 0
	for scanner.Scan() {
		var vals []int
		for _, field := range strings.Fields(scanner.Text()) {
			val, _ := strconv.Atoi(field)
			vals = append(vals, val)
		}
		total += predict(vals, *secondFlag)
	}

	fmt.Println(total)
}

func predict(vals []int, prev bool) int {
	n := len(vals)
	matrix := make([][]int, n)
	for i := n - 1; i >= 0; i-- {
		matrix[i] = make([]int, i+1)
		for j := 0; j <= i; j++ {
			if i == n-1 {
				matrix[i][j] = vals[j]
			} else {
				matrix[i][j] = matrix[i+1][j+1] - matrix[i+1][j]
			}
		}
	}

	total := 0
	if prev {
		for i := 0; i < n; i++ {
			total = matrix[i][0] - total
		}
	} else {
		for i := 0; i < n; i++ {
			total += matrix[i][i]
		}
	}

	return total
}
