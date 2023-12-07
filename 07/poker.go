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

type Hand struct {
	vals []int
	bid  int
	rank int
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("poker.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var hands []Hand

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		vals := handVals(fields[0], *secondFlag)
		bid, _ := strconv.Atoi(fields[1])
		hands = append(hands, Hand{vals, bid, handRank(vals)})
	}

	sort.SliceStable(hands, func(i, j int) bool {
		if hands[i].rank == hands[j].rank {
			for k := 0; k < 5; k++ {
				if hands[i].vals[k] == hands[j].vals[k] {
					continue
				}
				return hands[i].vals[k] < hands[j].vals[k]
			}
		}
		return hands[i].rank < hands[j].rank
	})

	sum := 0
	for i, hand := range hands {
		sum += (i + 1) * hand.bid
	}
	fmt.Println(sum)
}

func handVals(handStr string, secondFlag bool) []int {
	vals := make([]int, 5)
	for i, card := range handStr {
		switch card {
		case 'T':
			vals[i] = 10
		case 'J':
			if secondFlag {
				vals[i] = 0
			} else {
				vals[i] = 11
			}
		case 'Q':
			vals[i] = 12
		case 'K':
			vals[i] = 13
		case 'A':
			vals[i] = 14
		default:
			vals[i] = int(card - '0')
		}
	}
	return vals
}

func handRank(hand []int) int {
	counts := make([]int, 15)
	for _, card := range hand {
		counts[card]++
	}

	// Handle jokers. Add them to the card with the highest count.
	if counts[0] > 0 {
		maxIndex := 1
		for i := 1; i < 15; i++ {
			if counts[i] > counts[maxIndex] {
				maxIndex = i
			}
		}
		counts[maxIndex] += counts[0]
		counts[0] = 0
	}

	revCount := make([]int, 6)
	for _, count := range counts {
		revCount[count]++
	}

	if revCount[5] == 1 {
		return 10 // Five of a kind
	} else if revCount[4] == 1 {
		return 9 // Four of a kind
	} else if revCount[3] == 1 && revCount[2] == 1 {
		return 8 // Full house
	} else if revCount[3] == 1 {
		return 7 // Three of a kind
	} else if revCount[2] == 2 {
		return 6 // Two pairs
	} else if revCount[2] == 1 {
		return 5 // One pair
	} else {
		return 4 // High card
	}
}
