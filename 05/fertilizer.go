package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type seedRange struct {
	start, end int
}

type rangeMap struct {
	destRangeStart, sourceRangeStart, length int
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("fertilizer.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	var seeds []seedRange
	var vals []int
	for _, seedString := range strings.Fields(strings.Split(scanner.Text(), ":")[1]) {
		val, err := strconv.Atoi(seedString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse seed: %v\n", err)
			os.Exit(1)
		}

		if *secondFlag {
			// Part 2: Seed ranges are of variable length. To process later.
			vals = append(vals, val)
		} else {
			// Part 1: All seed ranges are of length 1
			seeds = append(seeds, seedRange{val, val})
		}
	}

	if *secondFlag {
		// Part 2: Variable length seed ranges processed here.
		for i := 0; i < len(vals); i += 2 {
			seeds = append(seeds, seedRange{vals[i], vals[i] + vals[i+1] - 1})
		}
	}
	scanner.Scan() // Skip the empty line

	for scanner.Scan() {
		var ranges []rangeMap
		for {
			if !scanner.Scan() {
				break
			}
			t := scanner.Text()
			if t == "" {
				break
			}
			var r rangeMap
			fmt.Sscanf(t, "%d %d %d", &r.destRangeStart, &r.sourceRangeStart, &r.length)
			ranges = append(ranges, r)
		}

		sort.SliceStable(ranges, func(i, j int) bool {
			return ranges[i].sourceRangeStart < ranges[j].sourceRangeStart
		})

		var newSeeds []seedRange
		for _, seedRange := range seeds {
			newSeeds = append(newSeeds, mapSourceToDest(ranges, seedRange)...)
		}
		seeds = newSeeds
	}

	min := seeds[0].start
	for _, seedRange := range seeds {
		if seedRange.start < min {
			min = seedRange.start
		}
	}
	fmt.Println(min)
}

// Maps the seed range to the destination ranges.
func mapSourceToDest(ranges []rangeMap, seed seedRange) []seedRange {
	var dest []seedRange
	currStart := seed.start
	for _, r := range ranges {
		if currStart > seed.end {
			break
		}
		if currStart < r.sourceRangeStart {
			dest = append(dest, seedRange{currStart, min(r.sourceRangeStart-1, seed.end)})
			if seed.end >= r.sourceRangeStart {
				dest = append(dest, seedRange{r.destRangeStart, r.destRangeStart + min(r.length-1, seed.end-r.sourceRangeStart)})
			}
			currStart = r.sourceRangeStart + r.length
		} else if currStart < r.sourceRangeStart+r.length {
			dest = append(dest, seedRange{r.destRangeStart + currStart - r.sourceRangeStart, r.destRangeStart + min(r.length-1, seed.end-r.sourceRangeStart)})
			currStart = r.sourceRangeStart + r.length
		}
	}
	if currStart <= seed.end {
		dest = append(dest, seedRange{currStart, seed.end})
	}
	return dest
}
