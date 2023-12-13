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

	file, err := os.Open("reflection.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var maze []string
	total := 0
	for scanner.Scan() {
		t := scanner.Text()
		if t != "" {
			maze = append(maze, t)
		} else {
			total += solve(maze, *secondFlag)
			maze = nil
		}
	}
	total += solve(maze, *secondFlag)

	fmt.Println(total)
}

func solve(maze []string, allowMismatch bool) int {
	rows := len(maze)
	cols := len(maze[0])

	// Assume horizontal reflection
	for i := 0; i < rows-1; i++ {
		mismatch := 0
		for j := 0; i-j >= 0 && i+1+j < rows && mismatch <= 1; j++ {
			for k := 0; k < cols && mismatch <= 1; k++ {
				if maze[i-j][k] != maze[i+1+j][k] {
					mismatch++
				}
			}
		}
		if (allowMismatch && mismatch == 1) || (!allowMismatch && mismatch == 0) {
			return 100 * (i + 1)
		}
	}

	// Assume vertical reflection
	for i := 0; i < cols-1; i++ {
		mismatch := 0
		for j := 0; i-j >= 0 && i+1+j < cols && mismatch <= 1; j++ {
			for k := 0; k < rows && mismatch <= 1; k++ {
				if maze[k][i-j] != maze[k][i+1+j] {
					mismatch++
				}
			}
		}
		if (allowMismatch && mismatch == 1) || (!allowMismatch && mismatch == 0) {
			return i + 1
		}
	}

	return 0
}
