package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type node struct {
	l, r string
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("camel.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	nameToNode := make(map[string]node)
	scanner := bufio.NewScanner(file)

	// Scan instructions
	scanner.Scan()
	moves := scanner.Text()
	scanner.Scan() // Skip blank line

	for scanner.Scan() {
		t := scanner.Text()
		name := t[:3]
		l := t[7:10]
		r := t[12:15]
		nameToNode[name] = node{l, r}
	}

	steps := nextZ("AAA", moves, nameToNode)
	if *secondFlag {
		// Find all names that has an A in it
		names := make([]string, 0)
		for name := range nameToNode {
			if strings.Contains(name, "A") {
				names = append(names, name)
			}
		}

		steps = 1
		for _, name := range names {
			// Implicitly assumes that a name always reaches Z in the same number of steps
			s := nextZ(name, moves, nameToNode)
			steps = lcm(steps, s)
		}
	}

	fmt.Println(steps)
}

// Returns the number of steps to the next Z
func nextZ(name string, moves string, nameToNode map[string]node) int {
	curr := name
	count := 0

	for count == 0 || !strings.Contains(curr, "Z") {
		move := moves[count%len(moves)]
		if move == 'L' {
			curr = nameToNode[curr].l
		} else {
			curr = nameToNode[curr].r
		}
		count += 1
	}

	return count
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
