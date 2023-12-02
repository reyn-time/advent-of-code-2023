package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("cube.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	red := regexp.MustCompile(` ([0-9]+) red`)
	green := regexp.MustCompile(` ([0-9]+) green`)
	blue := regexp.MustCompile(` ([0-9]+) blue`)
	lineNum := 1
	sum := 0
	for scanner.Scan() {
		line := scanner.Text()

		if *secondFlag {
			// The second part of the puzzle
			sum += max(getValues(line, red)) * max(getValues(line, green)) * max(getValues(line, blue))
			continue
		}

		// The first part of the puzzle
		good := all(getValues(line, red), func(v int) bool { return v <= 12 }) &&
			all(getValues(line, green), func(v int) bool { return v <= 13 }) &&
			all(getValues(line, blue), func(v int) bool { return v <= 14 })

		if good {
			sum += lineNum
		}
		lineNum++
	}

	fmt.Println(sum)
}

func getValues(line string, regex *regexp.Regexp) []int {
	var values []int
	for _, matches := range regex.FindAllStringSubmatch(line, -1) {
		v, _ := strconv.Atoi(matches[1])
		values = append(values, v)
	}
	return values
}

func all[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if !pred(t) {
			return false
		}
	}
	return true
}

func max(ts []int) int {
	if len(ts) == 0 {
		return 0
	}
	return slices.Max(ts)
}
