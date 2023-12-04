package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("scratch.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	total := 0
	next := make([]int, 1000)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ": ")
		parts = strings.Split(parts[1], " | ")
		winning := strings.Fields(parts[0])
		all := strings.Fields(parts[1])

		matches := 0
		for _, num := range winning {
			if slices.Contains(all, num) {
				matches += 1
			}
		}

		if !*secondFlag {
			if matches > 0 {
				total += (1 << (matches - 1))
			}
		} else {
			cardCount := next[0] + 1
			total += cardCount

			// Update next cards count
			for i := 1; i <= matches; i++ {
				next[i] += cardCount
			}

			// Dequeue current card
			next = next[1:]
		}
	}

	fmt.Println(total)
}
