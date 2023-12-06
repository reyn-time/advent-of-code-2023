package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("race.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var times, distances []int
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if *secondFlag {
		times = lineToInts2(scanner.Text())
	} else {
		times = lineToInts(scanner.Text())
	}
	scanner.Scan()
	if *secondFlag {
		distances = lineToInts2(scanner.Text())
	} else {
		distances = lineToInts(scanner.Text())
	}

	product := 1
	for i, time := range times {
		product *= waysToReachDistance(time, distances[i])
	}
	fmt.Println(product)
}

func lineToInts(line string) []int {
	fields := strings.Fields(strings.Split(line, ":")[1])
	vals := make([]int, len(fields))
	for i, field := range fields {
		vals[i], _ = strconv.Atoi(field)
	}
	return vals
}
func lineToInts2(line string) []int {
	numString := strings.Replace(strings.Split(line, ":")[1], " ", "", -1)
	val, _ := strconv.Atoi(numString)
	return []int{val}
}

// Counts number of x such that x*(time-x) > distance.
// That's equivalent to solving the quadratic equation x^2 - time*x + distance = 0.
func waysToReachDistance(time, distance int) int {
	root := math.Sqrt(float64(time*time - 4*distance))
	left := int(math.Ceil((float64(time) - root) / 2))
	right := int(math.Floor((float64(time) + root) / 2))
	return right - left + 1
}
