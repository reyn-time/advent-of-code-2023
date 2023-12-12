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

	file, err := os.Open("springs.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Scan values into cells
	scanner := bufio.NewScanner(file)

	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		// Puzzle part.
		puzzleString := parts[0]
		if *secondFlag {
			// Part 2: repeat the puzzle 5 times with ? sandwiched between.
			puzzleString = parts[0] + "?" + parts[0] + "?" + parts[0] + "?" + parts[0] + "?" + parts[0]
		}
		puzzle := []rune(puzzleString)

		// Hint part.
		hintStrings := strings.Split(parts[1], ",")
		hints := make([]int, len(hintStrings))
		for i, s := range hintStrings {
			v, _ := strconv.Atoi(s)
			hints[i] = v
		}
		if *secondFlag {
			// Part 2: repeat the hints 5 times.
			hintsCopy := hints
			for i := 0; i < 4; i++ {
				hints = append(hints, hintsCopy...)
			}
		}

		solution := solve(puzzle, hints)
		// fmt.Println(string(puzzle), solution)
		total += solution
	}

	fmt.Println(total)
}

func solve(puzzle []rune, hints []int) int {
	puzzle = append(puzzle, '.') // Add padding to the end for easier processing
	count := 0
	solveHelper(puzzle, hints, 0, make([]int, len(hints)), 0, false, &count, make(map[string]int))
	return count
}

// Solve the puzzle by exhaustion with memoization
func solveHelper(puzzle []rune, hints []int, i int, segmentLens []int, segment int, onSegment bool, count *int, memory map[string]int) {
	if i == len(puzzle) {
		// Check if the solution matches the hints
		if segment == len(hints) {
			*count++
		}
		return
	}

	// Check if the solution is already computed
	if v, ok := memory[solveHelperIndex(i, segmentLens, segment, onSegment)]; ok {
		*count += v
		return
	}

	prevCount := *count
	switch puzzle[i] {
	case '.':
		if onSegment && segmentLens[segment] == hints[segment] {
			// Leave segment
			solveHelper(puzzle, hints, i+1, segmentLens, segment+1, false, count, memory)
		} else if !onSegment {
			solveHelper(puzzle, hints, i+1, segmentLens, segment, false, count, memory)
		}
	case '#':
		if segment < len(hints) && segmentLens[segment]+1 <= hints[segment] {
			segmentLens[segment]++
			// Continue/start segment.
			solveHelper(puzzle, hints, i+1, segmentLens, segment, true, count, memory)
			segmentLens[segment]--
		}
	case '?':
		// Try '.'
		if onSegment && segmentLens[segment] == hints[segment] {
			puzzle[i] = '.'
			solveHelper(puzzle, hints, i+1, segmentLens, segment+1, false, count, memory)
		} else if !onSegment {
			puzzle[i] = '.'
			solveHelper(puzzle, hints, i+1, segmentLens, segment, false, count, memory)
		}

		// Try '#'
		if segment < len(hints) && segmentLens[segment]+1 <= hints[segment] {
			puzzle[i] = '#'
			segmentLens[segment]++
			solveHelper(puzzle, hints, i+1, segmentLens, segment, true, count, memory)
			segmentLens[segment]--
		}
		puzzle[i] = '?'
	}

	// Save the solution to memory
	memory[solveHelperIndex(i, segmentLens, segment, onSegment)] = *count - prevCount
}

func solveHelperIndex(i int, segmentLens []int, segment int, onSegment bool) string {
	if segment == len(segmentLens) {
		return fmt.Sprintf("%d-%d-%d-%t", i, 0, segment, onSegment)
	}
	return fmt.Sprintf("%d-%d-%d-%t", i, segmentLens[segment], segment, onSegment)
}
