package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("gear.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var puzzle []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// Add padding to the left and right.
		// This makes it easier to check if a number is adjacent to a symbol.
		puzzle = append(puzzle, "."+scanner.Text()+".")
	}

	total := 0
	var gearNumbers = make(map[gear][]int)

	for i, line := range puzzle {
		is_digit := false
		start_j := 0

		for j, char := range line {
			if char >= '0' && char <= '9' && !is_digit {
				is_digit = true
				start_j = j
			} else if (char < '0' || char > '9') && is_digit {
				is_digit = false
				if !*secondFlag {
					// Number found. Check if it is adjacent to a symbol.
					if is_adjacent_to_symbol(puzzle, i, start_j, j-1) {
						val, _ := strconv.Atoi(line[start_j:j])
						total += val
					}
				} else {
					// Number found. Check if it is adjacent to gears.
					// Keep track of the numbers for each gear.
					gears := all_adjacent_gears(puzzle, i, start_j, j-1)
					for _, gear := range gears {
						val, _ := strconv.Atoi(line[start_j:j])
						gearNumbers[gear] = append(gearNumbers[gear], val)
					}
				}
			}
		}
	}

	if *secondFlag {
		for _, numbers := range gearNumbers {
			if len(numbers) == 2 {
				total += numbers[0] * numbers[1]
			}
		}
	}
	fmt.Println(total)
}

func is_adjacent_to_symbol(puzzle []string, i, start_j, end_j int) bool {
	// Check top and bottom rows adjacent to the number.
	if i > 0 {
		for j := start_j - 1; j <= end_j+1; j++ {
			if is_symbol(puzzle[i-1][j]) {
				return true
			}
		}
	}
	if i+1 < len(puzzle) {
		for j := start_j - 1; j <= end_j+1; j++ {
			if is_symbol(puzzle[i+1][j]) {
				return true
			}
		}
	}

	// Check left and right cell adjacent to the number.
	return is_symbol(puzzle[i][start_j-1]) || is_symbol(puzzle[i][end_j+1])
}

type gear struct {
	i, j int
}

func all_adjacent_gears(puzzle []string, i, start_j, end_j int) []gear {
	var gears []gear

	// Check top and bottom rows adjacent to the number.
	if i > 0 {
		for j := start_j - 1; j <= end_j+1; j++ {
			if is_gear(puzzle[i-1][j]) {
				gears = append(gears, gear{i - 1, j})
			}
		}
	}
	if i+1 < len(puzzle) {
		for j := start_j - 1; j <= end_j+1; j++ {
			if is_gear(puzzle[i+1][j]) {
				gears = append(gears, gear{i + 1, j})
			}
		}
	}

	// Check left and right cell adjacent to the number.
	if is_gear(puzzle[i][start_j-1]) {
		gears = append(gears, gear{i, start_j - 1})
	}
	if is_gear(puzzle[i][end_j+1]) {
		gears = append(gears, gear{i, end_j + 1})
	}
	return gears
}

func is_symbol(char byte) bool {
	return char != '.' && (char < '0' || char > '9')
}

func is_gear(char byte) bool {
	return char == '*'
}
