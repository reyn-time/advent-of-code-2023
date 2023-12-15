package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Lens struct {
	focal      int
	prev, next *Lens
}

type Box struct {
	end *Lens
	m   map[string]*Lens
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("hash.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	input := scanner.Text()

	if !*secondFlag {
		// Part 1
		total := 0
		for _, s := range strings.Split(input, ",") {
			total += hash(s)
		}
		fmt.Println(total)
		return
	}

	// Part 2
	boxes := make([]Box, 256)
	for i := range boxes {
		boxes[i].m = make(map[string]*Lens)
	}
	for _, s := range strings.Split(input, ",") {
		if s[len(s)-1] == '-' {
			label := s[:len(s)-1]
			h := hash(label)
			// Remove this lens from the box
			if _, found := boxes[h].m[label]; found {
				if boxes[h].m[label].prev != nil {
					boxes[h].m[label].prev.next = boxes[h].m[label].next
				}
				if boxes[h].m[label].next != nil {
					boxes[h].m[label].next.prev = boxes[h].m[label].prev
				}
				if boxes[h].end == boxes[h].m[label] {
					boxes[h].end = boxes[h].m[label].prev
				}
				delete(boxes[h].m, label)
			}
		} else {
			label := s[:len(s)-2]
			focal := int(s[len(s)-1]) - int('0')
			h := hash(label)
			// Add this lens to the box
			if _, found := boxes[h].m[label]; !found {
				l := &Lens{focal, boxes[h].end, nil}
				boxes[h].m[label] = l
				if boxes[h].end != nil {
					boxes[h].end.next = l
				}
				boxes[h].end = l
			} else {
				boxes[h].m[label].focal = focal
			}
		}
	}

	total := 0
	for i, box := range boxes {
		slotNumber := len(box.m)
		for l := box.end; l != nil; l = l.prev {
			total += slotNumber * l.focal * (i + 1)
			slotNumber--
		}
	}
	fmt.Println(total)
}

func hash(s string) int {
	var hash int
	for _, c := range s {
		hash = (hash + int(c)) * 17 % 256
	}
	return hash
}
